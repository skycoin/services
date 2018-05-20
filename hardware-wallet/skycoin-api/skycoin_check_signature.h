
#ifndef SKYCOIN_CHECK_SIGNATURE_H
#define SKYCOIN_CHECK_SIGNATURE_H

#include <stdint.h>

int recover_pubkey_from_signed_message(const char* message, const uint8_t* signature, uint8_t* pubkey);

#endif