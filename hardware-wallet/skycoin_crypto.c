#include "skycoin_crypto.h"

extern void bn_print(const bignum256 *a);

void create_node(const char* seed_str, HDNode* node)
{
    const char* curve_name = SECP256K1_NAME; 
    hdnode_from_seed((const uint8_t *)seed_str, strlen(seed_str), curve_name, node);
    hdnode_fill_public_key(node);
}

void genereate_deterministic_key_pair_seckey(const char* seed, uint8_t* seckey)
{
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum(seed, digest, strlen(seed));

    compute_sha256sum((const char * )digest, seckey, SHA256_DIGEST_LENGTH);
}


void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght)
{
    SHA256_CTX ctx;
    sha256_Init(&ctx);
    sha256_Update(&ctx, (const uint8_t*) seed, seed_lenght);
    sha256_Final(&ctx, digest);
}
