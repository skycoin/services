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

void genereate_deterministic_key_pair(const uint8_t* seed, uint8_t* seckey, uint8_t* pubkey)
{
    compute_sha256sum((const char * )seed, seckey, SHA256_DIGEST_LENGTH);
    generate_pubkey_from_seckey(seckey, pubkey);
}

void ecdh_shared_secret(const uint8_t* secret_key, const uint8_t* remote_public_key, uint8_t* shared_secret /*should be size SHA256_DIGEST_LENGTH*/)
{
    uint8_t session_key1[33] = {0};
    uint8_t mult[65] = {0};
    char seed_str[256] = "dummy seed";
	HDNode dummy_node;
    create_node(seed_str, &dummy_node);
	ecdh_multiply(dummy_node.curve->params, secret_key, remote_public_key, mult); //65
	memcpy(&session_key1[1], &mult[1], 32);
	if (mult[64] % 2 == 0)
	{
		session_key1[0] = 0x02;
	}
	else
	{
		session_key1[0] = 0x03;
	}
    compute_sha256sum((char*)(session_key1), shared_secret, 33);
}


void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght)
{
    SHA256_CTX ctx;
    sha256_Init(&ctx);
    sha256_Update(&ctx, (const uint8_t*) seed, seed_lenght);
    sha256_Final(&ctx, digest);
}
