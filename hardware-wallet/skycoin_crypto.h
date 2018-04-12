#ifndef SKYCOIN_CRYPTO_H
#define SKYCOIN_CRYPTO_H

#include <stdint.h>
#include <stddef.h>

void ecdh(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* ecdh_key /*should be size SHA256_DIGEST_LENGTH*/);
void ecdh_shared_secret(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* shared_secret /*should be size SHA256_DIGEST_LENGTH*/);
void secp256k1Hash(const char* seed, uint8_t* secp256k1Hash_digest);
void generate_deterministic_key_pair_iterator(const char* seed, uint8_t* seckey, uint8_t* pubkey);
void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght);
void generate_pubkey_from_seckey(const uint8_t* seckey, uint8_t* pubkey);
void generate_deterministic_key_pair(const uint8_t* seed, const size_t seed_length, uint8_t* seckey, uint8_t* pubkey);
void generate_base58_address_from_pubkey(const uint8_t* pubkey, char* address, size_t *size_address);
void generate_bitcoin_address_from_pubkey(const uint8_t* pubkey, char* address, size_t *size_address);
void generate_bitcoin_private_address_from_seckey(const uint8_t* pubkey, char* address, size_t *size_address);
int recover_pubkey_from_signed_message(const char* message, const uint8_t* signature, uint8_t* pubkey);
int ecdsa_skycoin_sign(const uint8_t *priv_key, const uint8_t *digest, uint8_t *sig, uint8_t *pby);
void tohex(char * str, const uint8_t* buffer, int bufferLength);
#endif
