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
#include <ctype.h>

#include "layout2.h"
#include "storage.h"
#include "oled.h"
#include "bitmaps.h"
#include "string.h"
#include "util.h"
#include "timer.h"
#include "bignum.h"
#include "secp256k1.h"
#include "gettext.h"

#define BITCOIN_DIVISIBILITY (8)
#define BIP32_MAX_LAST_ELEMENT 1000000

// split longer string into 4 rows, rowlen chars each
static const char **split_message(const uint8_t *msg, uint32_t len, uint32_t rowlen)
{
	static char str[4][32 + 1];
	if (rowlen > 32) {
		rowlen = 32;
	}
	memset(str, 0, sizeof(str));
	strlcpy(str[0], (char *)msg, rowlen + 1);
	if (len > rowlen) {
		strlcpy(str[1], (char *)msg + rowlen, rowlen + 1);
	}
	if (len > rowlen * 2) {
		strlcpy(str[2], (char *)msg + rowlen * 2, rowlen + 1);
	}
	if (len > rowlen * 3) {
		strlcpy(str[3], (char *)msg + rowlen * 3, rowlen + 1);
	}
	if (len > rowlen * 4) {
		str[3][rowlen - 1] = '.';
		str[3][rowlen - 2] = '.';
		str[3][rowlen - 3] = '.';
	}
	static const char *ret[4] = { str[0], str[1], str[2], str[3] };
	return ret;
}

void *layoutLast = layoutHome;

void layoutDialogSwipe(const BITMAP *icon, const char *btnNo, const char *btnYes, const char *desc, const char *line1, const char *line2, const char *line3, const char *line4, const char *line5, const char *line6)
{
	layoutLast = layoutDialogSwipe;
	layoutSwipe();
	layoutDialog(icon, btnNo, btnYes, desc, line1, line2, line3, line4, line5, line6);
}

void layoutProgressSwipe(const char *desc, int permil)
{
	if (layoutLast == layoutProgressSwipe) {
		oledClear();
	} else {
		layoutLast = layoutProgressSwipe;
		layoutSwipe();
	}
	layoutProgress(desc, permil);
}

void layoutScreensaver(void)
{
	layoutLast = layoutScreensaver;
	oledClear();
	oledRefresh();
}

void layoutRawMessage(char* msg)
{
	oledClear();
	oledDrawStringCenter(OLED_HEIGHT/2, msg, FONT_STANDARD);
	oledRefresh();
}

void layoutHome(void)
{
	if (layoutLast == layoutHome || layoutLast == layoutScreensaver) {
		oledClear();
	} else {
		layoutSwipe();
	}
	layoutLast = layoutHome;
	const char *label = storage_isInitialized() ? storage_getLabel() : _("Go to trezor.io/start");
	const uint8_t *homescreen = bmp_skycoin_logo64.data;
	if (homescreen) {
		BITMAP b;
		b.width = 128;
		b.height = 64;
		b.data = homescreen;
		oledDrawBitmap(0, 0, &b);
	} else {
		if (label && strlen(label) > 0) {
			oledDrawBitmap(44, 4, &bmp_logo48);
			oledDrawStringCenter(OLED_HEIGHT - 8, label, FONT_STANDARD);
		} else {
			oledDrawBitmap(40, 0, &bmp_logo64);
		}
	}
	if (storage_needsBackup()) {
		oledBox(0, 0, 127, 8, false);
		oledDrawStringCenter(0, "NEEDS BACKUP!", FONT_STANDARD);
	}
	oledRefresh();

	// Reset lock screen timeout
	system_millis_lock_start = timer_ms();
}

void layoutSignMessage(const uint8_t *msg, uint32_t len)
{
	const char **str = split_message(msg, len, 16);
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"),
		_("Sign message?"),
		str[0], str[1], str[2], str[3], NULL, NULL);
}

void layoutVerifyAddress(const char *address)
{
	const char **str = split_message((const uint8_t *)address, strlen(address), 17);
	layoutDialogSwipe(&bmp_icon_info, _("Cancel"), _("Confirm"),
		_("Confirm address?"),
		_("Message signed by:"),
		str[0], str[1], str[2], NULL, NULL);
}

void layoutVerifyMessage(const uint8_t *msg, uint32_t len)
{
	const char **str = split_message(msg, len, 16);
	layoutDialogSwipe(&bmp_icon_info, _("Cancel"), _("Confirm"),
		_("Verified message"),
		str[0], str[1], str[2], str[3], NULL, NULL);
}

void layoutCipherKeyValue(bool encrypt, const char *key)
{
	const char **str = split_message((const uint8_t *)key, strlen(key), 16);
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"),
		encrypt ? _("Encrypt value of this key?") : _("Decrypt value of this key?"),
		str[0], str[1], str[2], str[3], NULL, NULL);
}

void layoutEncryptMessage(const uint8_t *msg, uint32_t len, bool signing)
{
	const char **str = split_message(msg, len, 16);
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"),
		signing ? _("Encrypt+Sign message?") : _("Encrypt message?"),
		str[0], str[1], str[2], str[3], NULL, NULL);
}

void layoutDecryptMessage(const uint8_t *msg, uint32_t len, const char *address)
{
	const char **str = split_message(msg, len, 16);
	layoutDialogSwipe(&bmp_icon_info, NULL, _("OK"),
		address ? _("Decrypted signed message") : _("Decrypted message"),
		str[0], str[1], str[2], str[3], NULL, NULL);
}

void layoutResetWord(const char *word, int pass, int word_pos, bool last)
{
	layoutLast = layoutResetWord;
	layoutSwipe();

	const char *btnYes;
	if (last) {
		if (pass == 1) {
			btnYes = _("Finish");
		} else {
			btnYes = _("Again");
		}
	} else {
		btnYes = _("Next");
	}

	const char *action;
	if (pass == 1) {
		action = _("Please check the seed");
	} else {
		action = _("Write down the seed");
	}

	char index_str[] = "##th word is:";
	if (word_pos < 10) {
		index_str[0] = ' ';
	} else {
		index_str[0] = '0' + word_pos / 10;
	}
	index_str[1] = '0' + word_pos % 10;
	if (word_pos == 1 || word_pos == 21) {
		index_str[2] = 's'; index_str[3] = 't';
	} else
	if (word_pos == 2 || word_pos == 22) {
		index_str[2] = 'n'; index_str[3] = 'd';
	} else
	if (word_pos == 3 || word_pos == 23) {
		index_str[2] = 'r'; index_str[3] = 'd';
	}

	int left = 0;
	oledClear();
	oledDrawBitmap(0, 0, &bmp_icon_info);
	left = bmp_icon_info.width + 4;

	oledDrawString(left, 0 * 9, action, FONT_STANDARD);
	oledDrawString(left, 2 * 9, word_pos < 10 ? index_str + 1 : index_str, FONT_STANDARD);
	oledDrawString(left, 3 * 9, word, FONT_STANDARD | FONT_DOUBLE);
	oledHLine(OLED_HEIGHT - 13);
	layoutButtonYes(btnYes);
	oledRefresh();
}

void layoutPublicKey(const uint8_t *pubkey)
{
	char hex[32 * 2 + 1], desc[16];
	strlcpy(desc, "Public Key: 00", sizeof(desc));
	if (pubkey[0] == 1) {
		/* ed25519 public key */
		// pass - leave 00
	} else {
		data2hex(pubkey, 1, desc + 12);
	}
	data2hex(pubkey + 1, 32, hex);
	const char **str = split_message((const uint8_t *)hex, 32 * 2, 16);
	layoutDialogSwipe(&bmp_icon_question, NULL, _("Continue"), NULL,
		desc, str[0], str[1], str[2], str[3], NULL);
}

void layoutSignIdentity(const IdentityType *identity, const char *challenge)
{
	char row_proto[8 + 11 + 1];
	char row_hostport[64 + 6 + 1];
	char row_user[64 + 8 + 1];

	bool is_gpg = (strcmp(identity->proto, "gpg") == 0);

	if (identity->has_proto && identity->proto[0]) {
		if (strcmp(identity->proto, "https") == 0) {
			strlcpy(row_proto, _("Web sign in to:"), sizeof(row_proto));
		} else if (is_gpg) {
			strlcpy(row_proto, _("GPG sign for:"), sizeof(row_proto));
		} else {
			strlcpy(row_proto, identity->proto, sizeof(row_proto));
			char *p = row_proto;
			while (*p) { *p = toupper((int)*p); p++; }
			strlcat(row_proto, _(" login to:"), sizeof(row_proto));
		}
	} else {
		strlcpy(row_proto, _("Login to:"), sizeof(row_proto));
	}

	if (identity->has_host && identity->host[0]) {
		strlcpy(row_hostport, identity->host, sizeof(row_hostport));
		if (identity->has_port && identity->port[0]) {
			strlcat(row_hostport, ":", sizeof(row_hostport));
			strlcat(row_hostport, identity->port, sizeof(row_hostport));
		}
	} else {
		row_hostport[0] = 0;
	}

	if (identity->has_user && identity->user[0]) {
		strlcpy(row_user, _("user: "), sizeof(row_user));
		strlcat(row_user, identity->user, sizeof(row_user));
	} else {
		row_user[0] = 0;
	}

	if (is_gpg) {
		// Split "First Last <first@last.com>" into 2 lines:
		// "First Last"
		// "first@last.com"
		char *email_start = strchr(row_hostport, '<');
		if (email_start) {
			strlcpy(row_user, email_start + 1, sizeof(row_user));
			*email_start = 0;
			char *email_end = strchr(row_user, '>');
			if (email_end) {
				*email_end = 0;
			}
		}
	}

	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"),
		_("Do you want to sign in?"),
		row_proto[0] ? row_proto : NULL,
		row_hostport[0] ? row_hostport : NULL,
		row_user[0] ? row_user : NULL,
		challenge,
		NULL,
		NULL);
}

void layoutDecryptIdentity(const IdentityType *identity)
{
	char row_proto[8 + 11 + 1];
	char row_hostport[64 + 6 + 1];
	char row_user[64 + 8 + 1];

	if (identity->has_proto && identity->proto[0]) {
		strlcpy(row_proto, identity->proto, sizeof(row_proto));
		char *p = row_proto;
		while (*p) { *p = toupper((int)*p); p++; }
		strlcat(row_proto, _(" decrypt for:"), sizeof(row_proto));
	} else {
		strlcpy(row_proto, _("Decrypt for:"), sizeof(row_proto));
	}

	if (identity->has_host && identity->host[0]) {
		strlcpy(row_hostport, identity->host, sizeof(row_hostport));
		if (identity->has_port && identity->port[0]) {
			strlcat(row_hostport, ":", sizeof(row_hostport));
			strlcat(row_hostport, identity->port, sizeof(row_hostport));
		}
	} else {
		row_hostport[0] = 0;
	}

	if (identity->has_user && identity->user[0]) {
		strlcpy(row_user, _("user: "), sizeof(row_user));
		strlcat(row_user, identity->user, sizeof(row_user));
	} else {
		row_user[0] = 0;
	}

	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"),
		_("Do you want to decrypt?"),
		row_proto[0] ? row_proto : NULL,
		row_hostport[0] ? row_hostport : NULL,
		row_user[0] ? row_user : NULL,
		NULL,
		NULL,
		NULL);
}

void layoutU2FDialog(const char *verb, const char *appname, const BITMAP *appicon) {
	if (!appicon) {
		appicon = &bmp_icon_question;
	}
	layoutDialog(appicon, NULL, verb, NULL, verb, _("U2F security key?"), NULL, appname, NULL, NULL);
}

static inline bool is_slip18(const uint32_t *address_n, size_t address_n_count)
{
	return address_n_count == 2 && address_n[0] == (0x80000000 + 10018) && (address_n[1] & 0x80000000) && (address_n[1] & 0x7FFFFFFF) <= 9;
}

void layoutCosiCommitSign(const uint32_t *address_n, size_t address_n_count, const uint8_t *data, uint32_t len, bool final_sign)
{
	char *desc = final_sign ? _("CoSi sign message?") : _("CoSi commit message?");
	char desc_buf[32];
	if (is_slip18(address_n, address_n_count)) {
		if (final_sign) {
			strlcpy(desc_buf, _("CoSi sign index #?"), sizeof(desc_buf));
			desc_buf[16] = '0' + (address_n[1] & 0x7FFFFFFF);
		} else {
			strlcpy(desc_buf, _("CoSi commit index #?"), sizeof(desc_buf));
			desc_buf[18] = '0' + (address_n[1] & 0x7FFFFFFF);
		}
		desc = desc_buf;
	}
	char str[4][17];
	if (len == 32) {
		data2hex(data     , 8, str[0]);
		data2hex(data +  8, 8, str[1]);
		data2hex(data + 16, 8, str[2]);
		data2hex(data + 24, 8, str[3]);
	} else {
		strlcpy(str[0], "Data", sizeof(str[0]));
		strlcpy(str[1], "of", sizeof(str[1]));
		strlcpy(str[2], "unsupported", sizeof(str[2]));
		strlcpy(str[3], "length", sizeof(str[3]));
	}
	layoutDialogSwipe(&bmp_icon_question, _("Cancel"), _("Confirm"), desc,
		str[0], str[1], str[2], str[3], NULL, NULL);
}
