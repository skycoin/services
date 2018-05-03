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

#ifndef __LAYOUT2_H__
#define __LAYOUT2_H__

#include "layout.h"
#include "types.pb.h"
#include "bitmaps.h"
#include "bignum.h"
#include "trezor.h"

extern void *layoutLast;

#if DEBUG_LINK
#define layoutSwipe oledClear
#else
#define layoutSwipe oledSwipeLeft
#endif

void layoutDialogSwipe(const BITMAP *icon, const char *btnNo, const char *btnYes, const char *desc, const char *line1, const char *line2, const char *line3, const char *line4, const char *line5, const char *line6);
void layoutProgressSwipe(const char *desc, int permil);

void layoutRawMessage(char* msg);
void layoutScreensaver(void);
void layoutHome(void);
void layoutSignMessage(const uint8_t *msg, uint32_t len);
void layoutVerifyAddress(const char *address);
void layoutVerifyMessage(const uint8_t *msg, uint32_t len);
void layoutCipherKeyValue(bool encrypt, const char *key);
void layoutEncryptMessage(const uint8_t *msg, uint32_t len, bool signing);
void layoutDecryptMessage(const uint8_t *msg, uint32_t len, const char *address);
void layoutResetWord(const char *word, int pass, int word_pos, bool last);
void layoutPublicKey(const uint8_t *pubkey);
void layoutSignIdentity(const IdentityType *identity, const char *challenge);
void layoutDecryptIdentity(const IdentityType *identity);
void layoutU2FDialog(const char *verb, const char *appname, const BITMAP *appicon);

void layoutCosiCommitSign(const uint32_t *address_n, size_t address_n_count, const uint8_t *data, uint32_t len, bool final_sign);

#endif
