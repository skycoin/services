#ifndef SKYCOIN_CRYPTO_H
#define SKYCOIN_CRYPTO_H

#include <stdint.h>

#include "sha2.h"

#include "bip32.h"

void create_node(const char* seed_str, HDNode* node);
void ecdh(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* ecdh_key /*should be size SHA256_DIGEST_LENGTH*/);
void ecdh_shared_secret(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* shared_secret /*should be size SHA256_DIGEST_LENGTH*/);
void secp256k1Hash(const char* seed, uint8_t* secp256k1Hash_digest);
void generate_deterministic_key_pair_iterator(const char* seed, uint8_t* seckey, uint8_t* pubkey);
void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght);
void generate_pubkey_from_seckey(const uint8_t* seckey, uint8_t* pubkey);
void genereate_deterministic_key_pair(const uint8_t* seed, const size_t seed_length, uint8_t* seckey, uint8_t* pubkey);

#endif
