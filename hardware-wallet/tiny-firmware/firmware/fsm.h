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

#ifndef __FSM_H__
#define __FSM_H__

#include "messages.pb.h"

// message functions

void fsm_sendSuccess(const char *text);

void fsm_sendFailure(FailureType code, const char *text);

void fsm_msgInitialize(Initialize *msg);
void fsm_msgGetFeatures(GetFeatures *msg);
int fsm_getKeyPairAtIndex(uint32_t nbAddress, uint8_t* pubkey, uint8_t* seckey, ResponseSkycoinAddress* respSkycoinAddress, uint32_t start_index);
void fsm_msgSkycoinCheckMessageSignature(SkycoinCheckMessageSignature* msg);
void fsm_msgSkycoinSignMessage(SkycoinSignMessage* msg);
void fsm_msgSkycoinAddress(SkycoinAddress* msg);
void fsm_msgSetMnemonic(SetMnemonic* msg);
void fsm_msgPing(Ping *msg);
void fsm_msgChangePin(ChangePin *msg);
void fsm_msgWipeDevice(WipeDevice *msg);
void fsm_msgGetEntropy(GetEntropy *msg);
void fsm_msgEntropyAck(EntropyAck *msg);
void fsm_msgLoadDevice(LoadDevice *msg);
void fsm_msgResetDevice(ResetDevice *msg);
void fsm_msgBackupDevice(BackupDevice *msg);
void fsm_msgPinMatrixAck(PinMatrixAck *msg);
void fsm_msgCancel(Cancel *msg);
void fsm_msgRecoveryDevice(RecoveryDevice *msg);
void fsm_msgWordAck(WordAck *msg);

#endif
