#include "skycoin_crypto.h"

extern void bn_print(const bignum256 *a);

void create_node(const char* seed_str, HDNode* node)
{
    const char* curve_name = SECP256K1_NAME;
    
    hdnode_from_seed((const uint8_t *)seed_str, strlen(seed_str), curve_name, node);
    hdnode_fill_public_key(node);

    printf("Pub key: ");
    for (int i = 0; i < 33; ++i)
    {
    printf("%02x", node->public_key[i]);
    }
    printf("\n");


    printf("Sec key: ");
    for (int i = 0; i < 32; ++i)
    {
    printf("%02x", node->private_key[i]);
    }
    printf("\n");
}


void generate_shared_key(const char *seed_str) {

    HDNode alice;
    // shared key variables
    int res, key_size;
    uint8_t session_key1[65] = {0};

    create_node(seed_str, &alice);

    res = hdnode_get_shared_key(&alice, alice.public_key, session_key1, &key_size);


    printf("Shared key (status: %d, key_size: %d): \n", res, key_size);
    for (int i = 0; i < key_size; ++i)
    {
        printf("%02x", session_key1[i]);
    }

    bignum256 bigshared;
    bn_read_be(session_key1, &bigshared);
    printf("\nbignum sharedkey : \n");
    bn_print(&bigshared);
    printf("\ncurve prime : \n");
    bn_print(&alice.curve->params->order);

    printf("\nbignum sharedkey after mod : \n");
    bn_fast_mod(&bigshared, &alice.curve->params->order);
    bn_print(&bigshared);
    printf("\n");

    uint8_t digest[SHA256_DIGEST_LENGTH]= {0};
    compute_sha256sum((const char*) session_key1, digest, sizeof(session_key1));
    
    printf("\nSha256 of output: \n");
    for(uint i=0;i<SHA256_DIGEST_LENGTH;i++) {
        printf("%02x", digest[i]);
    }
    printf("\n");
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
