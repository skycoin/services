#ifndef SKYCOIN_CRYPTO_H
#define SKYCOIN_CRYPTO_H

#include <stdint.h>

#include "sha2.h"

#include "bip32.h"

void create_node(const char* seed_str, HDNode* node);
void ecdh_shared_secret(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* shared_secret /*should be size SHA256_DIGEST_LENGTH*/);
void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght);
void genereate_deterministic_key_pair_seckey(const char* seed, uint8_t* seckey);

#endif
