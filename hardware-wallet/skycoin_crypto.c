#include "skycoin_crypto.h"

#include <string.h>

#include "curves.h"

extern void bn_print(const bignum256 *a);

void create_node(const char* seed_str, HDNode* node)
{
    const char* curve_name = SECP256K1_NAME; 
    hdnode_from_seed((const uint8_t *)seed_str, strlen(seed_str), curve_name, node);
    hdnode_fill_public_key(node);
}

void generate_pubkey_from_seckey(const uint8_t* seckey, uint8_t* pubkey)
{
    char seed_str[256] = "dummy seed";
	HDNode dummy_node;
    create_node(seed_str, &dummy_node);
	ecdsa_get_public_key33(dummy_node.curve->params, seckey, pubkey);
}

void genereate_deterministic_key_pair(const uint8_t* seed, const size_t seed_length, uint8_t* seckey, uint8_t* pubkey)
{
    compute_sha256sum((const char * )seed, seckey, seed_length);
    generate_pubkey_from_seckey(seckey, pubkey);
}

void ecdh(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* ecdh_key/*should be size SHA256_DIGEST_LENGTH*/)
{
    uint8_t mult[65] = {0};
    char seed_str[256] = "dummy seed";
	HDNode dummy_node;
    create_node(seed_str, &dummy_node);
	ecdh_multiply(dummy_node.curve->params, secret_key, remote_public_key, mult); //65
	memcpy(&ecdh_key[1], &mult[1], 32);
	if (mult[64] % 2 == 0)
	{
		ecdh_key[0] = 0x02;
	}
	else
	{
		ecdh_key[0] = 0x03;
	}
}

void ecdh_shared_secret(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* shared_secret /*should be size SHA256_DIGEST_LENGTH*/)
{
    uint8_t ecdh_key[33] = {0};
    ecdh(secret_key, remote_public_key, ecdh_key);
    compute_sha256sum((char*)(ecdh_key), shared_secret, 33);
}

void secp256k1Hash(const char* seed, uint8_t* secp256k1Hash_digest)
{
    uint8_t seckey[32] = {0};
    uint8_t dummy_seckey[32] = {0};
    uint8_t pubkey[33] = {0};
    uint8_t hash[SHA256_DIGEST_LENGTH] = {0};
    uint8_t hash2[SHA256_DIGEST_LENGTH] = {0};
    uint8_t ecdh_key[33] = {0};
    uint8_t secp256k1Hash[SHA256_DIGEST_LENGTH + 33] = {0};
    compute_sha256sum(seed, hash, strlen(seed));
    compute_sha256sum((const char*)hash, seckey, sizeof(hash));
    compute_sha256sum((const char*)hash, hash2, sizeof(hash));
    genereate_deterministic_key_pair(hash2, SHA256_DIGEST_LENGTH, dummy_seckey, pubkey);
    ecdh(seckey, pubkey, ecdh_key);
    memcpy(secp256k1Hash, hash, sizeof(hash));
    memcpy(&secp256k1Hash[SHA256_DIGEST_LENGTH], ecdh_key, sizeof(ecdh_key));
    compute_sha256sum((const char *)secp256k1Hash, secp256k1Hash_digest, sizeof(secp256k1Hash));
}

void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght)
{
    SHA256_CTX ctx;
    sha256_Init(&ctx);
    sha256_Update(&ctx, (const uint8_t*) seed, seed_lenght);
    sha256_Final(&ctx, digest);
}
