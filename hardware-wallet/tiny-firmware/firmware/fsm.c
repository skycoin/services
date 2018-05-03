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
#include "pinmatrix.h"
#include "layout2.h"
#include "address.h"
#include "base58.h"
#include "ecdsa.h"
#include "reset.h"
#include "recovery.h"
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

#define CHECK_PIN \
	if (!protectPin(true)) { \
		layoutHome(); \
		return; \
	}

#define CHECK_PIN_UNCACHED \
	if (!protectPin(false)) { \
		layoutHome(); \
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

static HDNode *fsm_getDerivedNode(const char *curve, const uint32_t *address_n, size_t address_n_count, uint32_t *fingerprint)
{
	static CONFIDENTIAL HDNode node;
	if (fingerprint) {
		*fingerprint = 0;
	}
	if (!storage_getRootNode(&node, curve, true)) {
		fsm_sendFailure(FailureType_Failure_NotInitialized, _("Device not initialized or passphrase request cancelled or unsupported curve"));
		layoutHome();
		return 0;
	}
	if (!address_n || address_n_count == 0) {
		return &node;
	}
	if (hdnode_private_ckd_cached(&node, address_n, address_n_count, fingerprint) == 0) {
		fsm_sendFailure(FailureType_Failure_ProcessError, _("Failed to derive private key"));
		layoutHome();
		return 0;
	}
	return &node;
}

void fsm_msgInitialize(Initialize *msg)
{
	recovery_abort();
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

	if (msg->has_pin_protection && msg->pin_protection) {
		CHECK_PIN
	}

	if (msg->has_passphrase_protection && msg->passphrase_protection) {
		if (!protectPassphrase()) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
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

void fsm_msgChangePin(ChangePin *msg)
{
	bool removal = msg->has_remove && msg->remove;
	if (removal) {
		if (storage_hasPin()) {
			layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("remove current PIN?"), NULL, NULL, NULL, NULL);
		} else {
			fsm_sendSuccess(_("PIN removed"));
			return;
		}
	} else {
		if (storage_hasPin()) {
			layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("change current PIN?"), NULL, NULL, NULL, NULL);
		} else {
			layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("set new PIN?"), NULL, NULL, NULL, NULL);
		}
	}
	if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}

	CHECK_PIN_UNCACHED

	if (removal) {
		storage_setPin("");
		storage_update();
		fsm_sendSuccess(_("PIN removed"));
	} else {
		if (protectChangePin()) {
			fsm_sendSuccess(_("PIN changed"));
		} else {
			fsm_sendFailure(FailureType_Failure_PinMismatch, NULL);
		}
	}
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

void fsm_msgGetEntropy(GetEntropy *msg)
{
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("send entropy?"), NULL, NULL, NULL, NULL);
	if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}
	
	RESP_INIT(Entropy);
	uint32_t len = msg->size;
	if (len > 1024) {
		len = 1024;
	}
	resp->entropy.size = len;
	random_buffer(resp->entropy.bytes, len);
	msg_write(MessageType_MessageType_Entropy, resp);
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

void fsm_msgResetDevice(ResetDevice *msg)
{
	CHECK_NOT_INITIALIZED

	CHECK_PARAM(!msg->has_strength || msg->strength == 128 || msg->strength == 192 || msg->strength == 256, _("Invalid seed strength"));

	reset_init(
		msg->has_display_random && msg->display_random,
		msg->has_strength ? msg->strength : 128,
		msg->has_passphrase_protection && msg->passphrase_protection,
		msg->has_pin_protection && msg->pin_protection,
		msg->has_language ? msg->language : 0,
		msg->has_label ? msg->label : 0,
		msg->has_u2f_counter ? msg->u2f_counter : 0,
		msg->has_skip_backup ? msg->skip_backup : false
	);
}

void fsm_msgBackupDevice(BackupDevice *msg)
{
	CHECK_INITIALIZED

	CHECK_PIN_UNCACHED

	(void)msg;
	reset_backup(true);
}

void fsm_msgCancel(Cancel *msg)
{
	(void)msg;
	recovery_abort();
	fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
}

void fsm_msgCipherKeyValue(CipherKeyValue *msg)
{
	CHECK_INITIALIZED

	CHECK_PARAM(msg->has_key, _("No key provided"));
	CHECK_PARAM(msg->has_value, _("No value provided"));
	CHECK_PARAM(msg->value.size % 16 == 0, _("Value length must be a multiple of 16"));

	CHECK_PIN

	const HDNode *node = fsm_getDerivedNode(SECP256K1_NAME, msg->address_n, msg->address_n_count, NULL);
	if (!node) return;

	bool encrypt = msg->has_encrypt && msg->encrypt;
	bool ask_on_encrypt = msg->has_ask_on_encrypt && msg->ask_on_encrypt;
	bool ask_on_decrypt = msg->has_ask_on_decrypt && msg->ask_on_decrypt;
	if ((encrypt && ask_on_encrypt) || (!encrypt && ask_on_decrypt)) {
		layoutCipherKeyValue(encrypt, msg->key);
		if (!protectButton(ButtonRequestType_ButtonRequest_Other, false)) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
			layoutHome();
			return;
		}
	}

	uint8_t data[256 + 4];
	strlcpy((char *)data, msg->key, sizeof(data));
	strlcat((char *)data, ask_on_encrypt ? "E1" : "E0", sizeof(data));
	strlcat((char *)data, ask_on_decrypt ? "D1" : "D0", sizeof(data));

	hmac_sha512(node->private_key, 32, data, strlen((char *)data), data);

	RESP_INIT(CipheredKeyValue);
	if (encrypt) {
		aes_encrypt_ctx ctx;
		aes_encrypt_key256(data, &ctx);
		aes_cbc_encrypt(msg->value.bytes, resp->value.bytes, msg->value.size, ((msg->iv.size == 16) ? (msg->iv.bytes) : (data + 32)), &ctx);
	} else {
		aes_decrypt_ctx ctx;
		aes_decrypt_key256(data, &ctx);
		aes_cbc_decrypt(msg->value.bytes, resp->value.bytes, msg->value.size, ((msg->iv.size == 16) ? (msg->iv.bytes) : (data + 32)), &ctx);
	}
	resp->has_value = true;
	resp->value.size = msg->value.size;
	msg_write(MessageType_MessageType_CipheredKeyValue, resp);
	layoutHome();
}

void fsm_msgClearSession(ClearSession *msg)
{
	(void)msg;
	session_clear(true); // clear PIN as well
	layoutScreensaver();
	fsm_sendSuccess(_("Session cleared"));
}

void fsm_msgApplySettings(ApplySettings *msg)
{
	CHECK_PARAM(msg->has_label || msg->has_language || msg->has_use_passphrase || msg->has_homescreen, _("No setting provided"));

	CHECK_PIN

	if (msg->has_label) {
		layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("change name to"), msg->label, "?", NULL, NULL);
		if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
			layoutHome();
			return;
		}
	}
	if (msg->has_language) {
		layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("change language to"), msg->language, "?", NULL, NULL);
		if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
			layoutHome();
			return;
		}
	}
	if (msg->has_use_passphrase) {
		layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), msg->use_passphrase ? _("enable passphrase") : _("disable passphrase"), _("encryption?"), NULL, NULL, NULL);
		if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
			layoutHome();
			return;
		}
	}
	if (msg->has_homescreen) {
		layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you really want to"), _("change the home"), _("screen?"), NULL, NULL, NULL);
		if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
			fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
			layoutHome();
			return;
		}
	}

	if (msg->has_label) {
		storage_setLabel(msg->label);
	}
	if (msg->has_language) {
		storage_setLanguage(msg->language);
	}
	if (msg->has_use_passphrase) {
		storage_setPassphraseProtection(msg->use_passphrase);
	}
	if (msg->has_homescreen) {
		storage_setHomescreen(msg->homescreen.bytes, msg->homescreen.size);
	}
	storage_update();
	fsm_sendSuccess(_("Settings applied"));
	layoutHome();
}

void fsm_msgApplyFlags(ApplyFlags *msg)
{
	if (msg->has_flags) {
		storage_applyFlags(msg->flags);
	}
	fsm_sendSuccess(_("Flags applied"));
}

void fsm_msgEntropyAck(EntropyAck *msg)
{
	if (msg->has_entropy) {
		reset_entropy(msg->entropy.bytes, msg->entropy.size);
	} else {
		reset_entropy(0, 0);
	}
}

void fsm_msgRecoveryDevice(RecoveryDevice *msg)
{
	const bool dry_run = msg->has_dry_run ? msg->dry_run : false;
	if (dry_run) {
		CHECK_PIN
	} else {
		CHECK_NOT_INITIALIZED
	}

	CHECK_PARAM(!msg->has_word_count || msg->word_count == 12 || msg->word_count == 18 || msg->word_count == 24, _("Invalid word count"));

	recovery_init(
		msg->has_word_count ? msg->word_count : 12,
		msg->has_passphrase_protection && msg->passphrase_protection,
		msg->has_pin_protection && msg->pin_protection,
		msg->has_language ? msg->language : 0,
		msg->has_label ? msg->label : 0,
		msg->has_enforce_wordlist && msg->enforce_wordlist,
		msg->has_type ? msg->type : 0,
		msg->has_u2f_counter ? msg->u2f_counter : 0,
		dry_run
	);
}

void fsm_msgWordAck(WordAck *msg)
{
	recovery_word(msg->word);
}

void fsm_msgSetU2FCounter(SetU2FCounter *msg)
{
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), NULL, _("Do you want to set"), _("the U2F counter?"), NULL, NULL, NULL, NULL);
	if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}
	storage_setU2FCounter(msg->u2f_counter);
	storage_update();
	fsm_sendSuccess(_("U2F counter set"));
	layoutHome();
}

void fsm_msgCosiCommit(CosiCommit *msg)
{
	RESP_INIT(CosiCommitment);

	CHECK_INITIALIZED

	CHECK_PARAM(msg->has_data, _("No data provided"));

	layoutCosiCommitSign(msg->address_n, msg->address_n_count, msg->data.bytes, msg->data.size, false);
	if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}

	CHECK_PIN

	const HDNode *node = fsm_getDerivedNode(ED25519_NAME, msg->address_n, msg->address_n_count, NULL);
	if (!node) return;

	uint8_t nonce[32];
	sha256_Raw(msg->data.bytes, msg->data.size, nonce);
	rfc6979_state rng;
	init_rfc6979(node->private_key, nonce, &rng);
	generate_rfc6979(nonce, &rng);

	resp->has_commitment = true;
	resp->has_pubkey = true;
	resp->commitment.size = 32;
	resp->pubkey.size = 32;

	ed25519_publickey(nonce, resp->commitment.bytes);
	ed25519_publickey(node->private_key, resp->pubkey.bytes);

	msg_write(MessageType_MessageType_CosiCommitment, resp);
	layoutHome();
}

void fsm_msgCosiSign(CosiSign *msg)
{
	RESP_INIT(CosiSignature);

	CHECK_INITIALIZED

	CHECK_PARAM(msg->has_data, _("No data provided"));
	CHECK_PARAM(msg->has_global_commitment && msg->global_commitment.size == 32, _("Invalid global commitment"));
	CHECK_PARAM(msg->has_global_pubkey && msg->global_pubkey.size == 32, _("Invalid global pubkey"));

	layoutCosiCommitSign(msg->address_n, msg->address_n_count, msg->data.bytes, msg->data.size, true);
	if (!protectButton(ButtonRequestType_ButtonRequest_ProtectCall, false)) {
		fsm_sendFailure(FailureType_Failure_ActionCancelled, NULL);
		layoutHome();
		return;
	}

	CHECK_PIN

	const HDNode *node = fsm_getDerivedNode(ED25519_NAME, msg->address_n, msg->address_n_count, NULL);
	if (!node) return;

	uint8_t nonce[32];
	sha256_Raw(msg->data.bytes, msg->data.size, nonce);
	rfc6979_state rng;
	init_rfc6979(node->private_key, nonce, &rng);
	generate_rfc6979(nonce, &rng);

	resp->has_signature = true;
	resp->signature.size = 32;

	ed25519_cosi_sign(msg->data.bytes, msg->data.size, node->private_key, nonce, msg->global_commitment.bytes, msg->global_pubkey.bytes, resp->signature.bytes);

	msg_write(MessageType_MessageType_CosiSignature, resp);
	layoutHome();
}

#if DEBUG_LINK

void fsm_msgDebugLinkGetState(DebugLinkGetState *msg)
{
	(void)msg;

	// Do not use RESP_INIT because it clears msg_resp, but another message
	// might be being handled
	DebugLinkState resp;
	memset(&resp, 0, sizeof(resp));

	resp.has_layout = true;
	resp.layout.size = OLED_BUFSIZE;
	memcpy(resp.layout.bytes, oledGetBuffer(), OLED_BUFSIZE);

	if (storage_hasPin()) {
		resp.has_pin = true;
		strlcpy(resp.pin, storage_getPin(), sizeof(resp.pin));
	}

	resp.has_matrix = true;
	strlcpy(resp.matrix, pinmatrix_get(), sizeof(resp.matrix));

	resp.has_reset_entropy = true;
	resp.reset_entropy.size = reset_get_int_entropy(resp.reset_entropy.bytes);

	resp.has_reset_word = true;
	strlcpy(resp.reset_word, reset_get_word(), sizeof(resp.reset_word));

	resp.has_recovery_fake_word = true;
	strlcpy(resp.recovery_fake_word, recovery_get_fake_word(), sizeof(resp.recovery_fake_word));

	resp.has_recovery_word_pos = true;
	resp.recovery_word_pos = recovery_get_word_pos();

	if (storage_hasMnemonic()) {
		resp.has_mnemonic = true;
		strlcpy(resp.mnemonic, storage_getMnemonic(), sizeof(resp.mnemonic));
	}

	if (storage_hasNode()) {
		resp.has_node = true;
		storage_dumpNode(&(resp.node));
	}

	resp.has_passphrase_protection = true;
	resp.passphrase_protection = storage_hasPassphraseProtection();

	msg_debug_write(MessageType_MessageType_DebugLinkState, &resp);
}

void fsm_msgDebugLinkStop(DebugLinkStop *msg)
{
	(void)msg;
}

void fsm_msgDebugLinkMemoryRead(DebugLinkMemoryRead *msg)
{
	RESP_INIT(DebugLinkMemory);

	uint32_t length = 1024;
	if (msg->has_length && msg->length < length)
		length = msg->length;
	resp->has_memory = true;
	memcpy(resp->memory.bytes, (void*) msg->address, length);
	resp->memory.size = length;
	msg_debug_write(MessageType_MessageType_DebugLinkMemory, resp);
}

void fsm_msgDebugLinkMemoryWrite(DebugLinkMemoryWrite *msg)
{
	uint32_t length = msg->memory.size;
	if (msg->flash) {
		flash_clear_status_flags();
		flash_unlock();
		for (uint32_t i = 0; i < length; i += 4) {
			uint32_t word;
			memcpy(&word, msg->memory.bytes + i, 4);
			flash_program_word(msg->address + i, word);
		}
		flash_lock();
	} else {
		memcpy((void *) msg->address, msg->memory.bytes, length);
	}
}

void fsm_msgDebugLinkFlashErase(DebugLinkFlashErase *msg)
{
	flash_clear_status_flags();
	flash_unlock();
	flash_erase_sector(msg->sector, FLASH_CR_PROGRAM_X32);
	flash_lock();
}
#endif
