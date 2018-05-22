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

#include <stdint.h>
#include <string.h>

#include "signatures.h"
#include "skycoin_check_signature.h"
#include "sha2.h"
#include "bootloader.h"

#define PUBKEYS 5

#if SIGNATURE_PROTECT
static const uint8_t * const pubkey[PUBKEYS] = {
	(const uint8_t *)"\x02\x58\x39\x07\x8e\x6c\x11\xc0\x9a\xd4\xb0\x00\x92\xf5\xfe\xef\xd5\x66\xf9\x2a\xf6\x0f\x2c\x71\xfa\xcf\xb0\x1d\x2b\x84\xc0\x43\xd4",
	(const uint8_t *)"\x02\x91\x70\x19\x2c\x2f\xde\xfb\x4d\x37\x7a\xf9\xe1\x96\xe0\x6d\x11\x76\xf8\x6f\x73\x35\x23\xd3\x95\x0f\x90\xff\x84\xc0\xcd\x02\xd3",
	(const uint8_t *)"\x03\x33\x8f\xfc\x0f\xf4\x2d\xf0\x7d\x27\xb0\xb4\x13\x1c\xd9\x6f\xfd\xfa\x46\x85\xb5\x56\x6a\xaf\xc7\xaa\x71\xed\x10\xfd\x1c\xbd\x6f",
	(const uint8_t *)"\x03\x9f\x12\xc9\x36\x45\xe3\x5e\x52\x74\xdc\x38\xf1\x91\xbe\x0b\x6d\x13\x21\xec\x35\xd2\xd2\xa3\xdd\xf7\xd1\x3e\xd1\x2f\x6d\xa8\x5b",
	(const uint8_t *)"\x03\xb1\x7c\x7b\x7c\x56\x43\x85\xbe\x66\xf9\xc1\xb9\xda\x6a\x0b\x5a\xea\x56\xf0\xcb\x70\x54\x8e\x65\x28\xa2\xf4\xf7\xb2\x72\x45\xd8",
};
#endif

#define SIGNATURES 3


#if SIGNATURE_DEBUG
static void displaySignatureDebug(const uint8_t *hash, const uint8_t *signature, 
const uint8_t *pubk, const uint8_t *stored_pubkey)
{
	layout32bits(hash, "Hash");
	layout32bits(signature, "Signature[0-31]");
	layout32bits(signature + 32, "Signature[32-64]");
	layout32bits(pubk, "Computed Pub");
	layout32bits(stored_pubkey, "Pubkey");

}
#endif

int signatures_ok(uint8_t *store_hash)
{
	if (!firmware_present()) return SIG_FAIL; // no firmware present

	const uint32_t codelen = *((const uint32_t *)FLASH_META_CODELEN);
	
	uint8_t hash[32];
	sha256_Raw((const uint8_t *)FLASH_APP_START, codelen, hash);
	if (store_hash) {
		memcpy(store_hash, hash, 32);
	}

#if SIGNATURE_PROTECT

	const uint8_t sigindex1 = *((const uint8_t *)FLASH_META_SIGINDEX1);
	const uint8_t sigindex2 = *((const uint8_t *)FLASH_META_SIGINDEX2);
	const uint8_t sigindex3 = *((const uint8_t *)FLASH_META_SIGINDEX3);

	if (sigindex1 < 1 || sigindex1 > PUBKEYS) return SIG_FAIL; // invalid index
	if (sigindex2 < 1 || sigindex2 > PUBKEYS) return SIG_FAIL; // invalid index
	if (sigindex3 < 1 || sigindex3 > PUBKEYS) return SIG_FAIL; // invalid index

	if (sigindex1 == sigindex2) return SIG_FAIL; // duplicate use
	if (sigindex1 == sigindex3) return SIG_FAIL; // duplicate use
	if (sigindex2 == sigindex3) return SIG_FAIL; // duplicate use

	uint8_t pubkey1[33];
	uint8_t pubkey2[33];
	uint8_t pubkey3[33];

	uint8_t sign1[65];
	uint8_t sign2[65];
	uint8_t sign3[65];

	memcpy(sign1, (const uint8_t *)FLASH_META_SIG1, 64);
	recover_pubkey_from_signed_message((char*)hash, sign1, pubkey1);
	if (0 != memcmp(pubkey1, pubkey[sigindex1 - 1], 33)) // failure
	{
#if SIGNATURE_DEBUG
		displaySignatureDebug(hash, sign1, pubkey1, pubkey[sigindex1 - 1]);
#endif
		return SIG_FAIL;
	} 
	memcpy(sign2, (const uint8_t *)FLASH_META_SIG2, 64);
	recover_pubkey_from_signed_message((char*)hash, sign2, pubkey2);
	if (0 != memcmp(pubkey2, pubkey[sigindex2 - 1], 33)) // failure
	{
#if SIGNATURE_DEBUG
		displaySignatureDebug(hash, sign2, pubkey2, pubkey[sigindex2 - 1]);
#endif
		return SIG_FAIL;
	} 
	memcpy(sign3, (const uint8_t *)FLASH_META_SIG3, 64);
	recover_pubkey_from_signed_message((char*)hash, sign3, pubkey3);
	if (0 != memcmp(pubkey3, pubkey[sigindex3 - 1], 33)) // failure
	{
#if SIGNATURE_DEBUG
		displaySignatureDebug(hash, sign3, pubkey3, pubkey[sigindex3 - 1]);
#endif
		return SIG_FAIL;
	} 
#endif

	return SIG_OK;
}
