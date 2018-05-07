/*
 * This file is part of the TREZOR project, https://trezor.io/
 *
 * Copyright (C) 2014 Pavol Rusnak <stick@satoshilabs.com>
 *
 * This library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this library.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <libopencm3/stm32/flash.h>

#include "trezor.h"
#include "fsm.h"
#include "messages.h"
#include "bip32.h"
#include "storage.h"
#include "rng.h"
#include "storage.h"
#include "oled.h"
#include "protect.h"
#include "layout2.h"
#include "address.h"
#include "base58.h"
#include "ecdsa.h"
#include "memory.h"
#include "usb.h"
#include "util.h"
#include "aes/aes.h"
#include "hmac.h"
#include "base58.h"
#include "bip39.h"
#include "ripemd160.h"
#include "curves.h"
#include "secp256k1.h"
#include "rfc6979.h"
#include "gettext.h"
#include "skycoin_crypto.h"

// message methods

static uint8_t msg_resp[MSG_OUT_SIZE] __attribute__ ((aligned));

#define RESP_INIT(TYPE) \
			TYPE *resp = (TYPE *) (void *) msg_resp; \
			_Static_assert(sizeof(msg_resp) >= sizeof(TYPE), #TYPE " is too large"); \
			memset(resp, 0, sizeof(TYPE));

#define CHECK_INITIALIZED \
	if (!storage_isInitialized()) { \
		fsm_sendFailure(FailureType_Failure_NotInitialized, NULL); \
		return; \
	}

#define CHECK_NOT_INITIALIZED \
	if (storage_isInitialized()) { \
		fsm_sendFailure(FailureType_Failure_UnexpectedMessage, _("Device is already initialized. Use Wipe first.")); \
		return; \
	}

#define CHECK_PARAM(cond, errormsg) \
	if (!(cond)) { \
		fsm_sendFailure(FailureType_Failure_DataError, (errormsg)); \
		layoutHome(); \
		return; \
	}

void fsm_sendSuccess(const char *text)
{
	RESP_INIT(Success);
	if (text) {
		resp->has_message = true;
		strlcpy(resp->message, text, sizeof(resp->message));
	}
	msg_write(MessageType_MessageType_Success, resp);
}

void fsm_sendFailure(FailureType code, const char *text)
{
	if (protectAbortedByInitialize) {
		fsm_msgInitialize((Initialize *)0);
		protectAbortedByInitialize = false;
		return;
	}
	RESP_INIT(Failure);
	resp->has_code = true;
	resp->code = code;
	if (!text) {
		switch (code) {
			case FailureType_Failure_UnexpectedMessage:
				text = _("Unexpected message");
				break;
			case FailureType_Failure_ButtonExpected:
				text = _("Button expected");
				break;
			case FailureType_Failure_DataError:
				text = _("Data error");
				break;
			case FailureType_Failure_ActionCancelled:
				text = _("Action cancelled by user");
				break;
			case FailureType_Failure_PinExpected:
				text = _("PIN expected");
				break;
			case FailureType_Failure_PinCancelled:
				text = _("PIN cancelled");
				break;
			case FailureType_Failure_PinInvalid:
				text = _("PIN invalid");
				break;
			case FailureType_Failure_InvalidSignature:
				text = _("Invalid signature");
				break;
			case FailureType_Failure_ProcessError:
				text = _("Process error");
				break;
			case FailureType_Failure_NotEnoughFunds:
				text = _("Not enough funds");
				break;
			case FailureType_Failure_NotInitialized:
				text = _("Device not initialized");
				break;
			case FailureType_Failure_PinMismatch:
				text = _("PIN mismatch");
				break;
			case FailureType_Failure_FirmwareError:
				text = _("Firmware error");
				break;
		}
	}
	if (text) {
		resp->has_message = true;
		strlcpy(resp->message, text, sizeof(resp->message));
	}
	msg_write(MessageType_MessageType_Failure, resp);
}

void fsm_msgInitialize(Initialize *msg)
{
	if (msg && msg->has_state && msg->state.size == 64) {
		uint8_t i_state[64];
		if (!session_getState(msg->state.bytes, i_state, NULL)) {
			session_clear(false); // do not clear PIN
		} else {
			if (0 != memcmp(msg->state.bytes, i_state, 64)) {
				session_clear(false); // do not clear PIN
			}
		}
	} else {
		session_clear(false); // do not clear PIN
	}
	layoutHome();
	fsm_msgGetFeatures(0);
}

void fsm_msgGetFeatures(GetFeatures *msg)
{
	(void)msg;
	RESP_INIT(Features);
	resp->has_vendor = true;         strlcpy(resp->vendor, "bitcointrezor.com", sizeof(resp->vendor));
	resp->has_major_version = true;  resp->major_version = VERSION_MAJOR;
	resp->has_minor_version = true;  resp->minor_version = VERSION_MINOR;
	resp->has_patch_version = true;  resp->patch_version = VERSION_PATCH;
	resp->has_device_id = true;      strlcpy(resp->device_id, storage_uuid_str, sizeof(resp->device_id));
	resp->has_pin_protection = true; resp->pin_protection = storage_hasPin();
	resp->has_passphrase_protection = true; resp->passphrase_protection = storage_hasPassphraseProtection();
#ifdef SCM_REVISION
	int len = sizeof(SCM_REVISION) - 1;
	resp->has_revision = true; memcpy(resp->revision.bytes, SCM_REVISION, len); resp->revision.size = len;
#endif
	resp->has_bootloader_hash = true; resp->bootloader_hash.size = memory_bootloader_hash(resp->bootloader_hash.bytes);
	if (storage_getLanguage()) {
		resp->has_language = true;
		strlcpy(resp->language, storage_getLanguage(), sizeof(resp->language));
	}
	if (storage_getLabel()) {
		resp->has_label = true;
		strlcpy(resp->label, storage_getLabel(), sizeof(resp->label));
	}
	
	resp->has_initialized = true; resp->initialized = storage_isInitialized();
	resp->has_imported = true; resp->imported = storage_isImported();
	resp->has_pin_cached = true; resp->pin_cached = session_isPinCached();
	resp->has_passphrase_cached = true; resp->passphrase_cached = session_isPassphraseCached();
	resp->has_needs_backup = true; resp->needs_backup = storage_needsBackup();
	resp->has_flags = true; resp->flags = storage_getFlags();
	resp->has_model = true; strlcpy(resp->model, "1", sizeof(resp->model));

	msg_write(MessageType_MessageType_Features, resp);
}

void fsm_msgSkycoinCheckMessageSignature(SkycoinCheckMessageSignature* msg)
{
	uint8_t sign[65];
    size_t size_sign = sizeof(sign);
    char pubkeybase58[36];
    uint8_t pubkey[33] = {0};
    uint8_t digest[32] = {0};

    RESP_INIT(Success);
    compute_sha256sum(msg->message, digest, strlen(msg->message));
    size_sign = sizeof(sign);
    b58tobin(sign, &size_sign, msg->signature);
    recover_pubkey_from_signed_message((char*)digest, sign, pubkey);
    size_sign = sizeof(pubkeybase58);
    generate_base58_address_from_pubkey(pubkey, pubkeybase58, &size_sign);
    if (memcmp(pubkeybase58, msg->address, size_sign) == 0)
    {
        layoutRawMessage("Verification success");
    }
    else {
        layoutRawMessage("Wrong signature");
    }
    memcpy(resp->message, pubkeybase58, size_sign);
    resp->has_message = true;
    msg_write(MessageType_MessageType_Success, resp);
}

int fsm_getKeyPairAtIndex(uint32_t index, uint8_t* pubkey, uint8_t* seckey)
{
    const uint8_t* mnemo = storage_getSeed(false);
    char seed[64] = {0};
    uint8_t nextSeed[SHA256_DIGEST_LENGTH] = {0};
    if (mnemo == NULL)
    {
        return -1;
    }
	memcpy(seed, mnemo, sizeof(seed));
	for (uint8_t i = 0; i < index; ++i)
	{
		generate_deterministic_key_pair_iterator(seed, nextSeed, seckey, pubkey);
		memcpy(seed, nextSeed, 32);
		seed[32] = 0;
	}
    return 0;
}

void fsm_msgSkycoinSignMessage(SkycoinSignMessage* msg)
{
    uint8_t pubkey[33] = {0};
    uint8_t seckey[32] = {0};
	uint8_t digest[32] = {0};
    size_t size_sign;
    uint8_t signature[65];
	char sign58[90] = {0};
	int res = 0;
	RESP_INIT(Success);
    fsm_getKeyPairAtIndex(msg->address_n, pubkey, seckey);
    compute_sha256sum(msg->message, digest, strlen(msg->message));
    res = ecdsa_skycoin_sign(rand(), seckey, digest, signature);
	if (res == 0)
	{
		layoutRawMessage("Signature success");
	}
	else
	{
		layoutRawMessage("Signature failed");
	}
	size_sign = sizeof(sign58);
    b58enc(sign58, &size_sign, signature, sizeof(signature));
	memcpy(resp->message, sign58, size_sign);
	resp->has_message = true;
	msg_write(MessageType_MessageType_Success, resp);
}

void fsm_msgSkycoinAddress(SkycoinAddress* msg)
{
    uint8_t seckey[32] = {0};
    uint8_t pubkey[33] = {0};

	RESP_INIT(Success);
	// reset_entropy((const uint8_t*)msg->seed, strlen(msg->seed));
	if (msg->has_address_type  && fsm_getKeyPairAtIndex(msg->address_n, pubkey, seckey) == 0)
	{
    	char address[256] = {0};
    	size_t size_address = sizeof(address);
		switch (msg->address_type)
		{
			case SkycoinAddressType_AddressTypeSkycoin:
				layoutRawMessage("Skycoin address");
    			generate_base58_address_from_pubkey(pubkey, address, &size_address);
				memcpy(resp->message, address, size_address);
				break;
			case SkycoinAddressType_AddressTypeBitcoin:
				layoutRawMessage("Bitcoin address");
				generate_bitcoin_address_from_pubkey(pubkey, address, &size_address);
				memcpy(resp->message, address, size_address);
				break;
			default:
				layoutRawMessage("Unknown address type");
				break;
		}
	}
	else {
		tohex(resp->message, pubkey, 33);
		layoutRawMessage(resp->message);
	}
	resp->has_message = true;
	msg_write(MessageType_MessageType_Success, resp);
}

void fsm_msgPing(Ping *msg)
{
	RESP_INIT(Success);

	if (msg->has_button_protection && msg->button_protection) {
		layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("answer to ping?"), NULL, NULL, NULL, NULL);
		if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
			layoutHome();
			return;
		}
	}

	if (msg->has_message) {
		resp->has_message = true;
		memcpy(&(resp->message), &(msg->message), sizeof(resp->message));
	}
	msg_write(MessageType_MessageType_Success, resp);
	layoutHome();
}

void fsm_msgWipeDevice(WipeDevice *msg)
{
	(void)msg;
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("wipe the device?"), NULL, _("All data will be lost."), NULL, NULL);
	if (!protectButton(ButtonRequestType_ButtonRequest_WipeDevice, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}
	storage_wipe();
	// the following does not work on Mac anyway :-/ Linux/Windows are fine, so it is not needed
	// usbReconnect(); // force re-enumeration because of the serial number change
	fsm_sendSuccess(_("Device wiped"));
	layoutHome();
}

void fsm_msgLoadDevice(LoadDevice *msg)
{
	CHECK_NOT_INITIALIZED

	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("I take the risk"), NULL, _("Loading private seed"), _("is not recommended."), _("Continue only if you"), _("know what you are"), _("doing!"), NULL);
	if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}

	if (msg->has_mnemonic && !(msg->has_skip_checksum && msg->skip_checksum) ) {
		if (!mnemonic_check(msg->mnemonic)) {
			fsm_sendFailure(FailureType_Failure_DataError, _("Mnemonic with wrong checksum provided"));
			layoutHome();
			return;
		}
	}

	storage_loadDevice(msg);
	fsm_sendSuccess(_("Device loaded"));
	layoutHome();
}

void fsm_msgCancel(Cancel *msg)
{
	(void)msg;
	fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
}
