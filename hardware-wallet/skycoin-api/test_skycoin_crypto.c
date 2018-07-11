#include "skycoin_crypto.h"
#include "skycoin_check_signature.h"

#include <stdio.h>
#include <string.h>

#include "check.h"
#include "sha2.h" //SHA256_DIGEST_LENGTH
#include "base58.h"
#include "ecdsa.h"
#include "secp256k1.h"
#include "check_digest.h"

#define FROMHEX_MAXLEN 512

const uint8_t *fromhex(const char *str)
{
    static uint8_t buf[FROMHEX_MAXLEN];
    size_t len = strlen(str) / 2;
    if (len > FROMHEX_MAXLEN) len = FROMHEX_MAXLEN;
    for (size_t i = 0; i < len; i++) {
        uint8_t c = 0;
        if (str[i * 2] >= '0' && str[i*2] <= '9') c += (str[i * 2] - '0') << 4;
        if ((str[i * 2] & ~0x20) >= 'A' && (str[i*2] & ~0x20) <= 'F') c += (10 + (str[i * 2] & ~0x20) - 'A') << 4;
        if (str[i * 2 + 1] >= '0' && str[i * 2 + 1] <= '9') c += (str[i * 2 + 1] - '0');
        if ((str[i * 2 + 1] & ~0x20) >= 'A' && (str[i * 2 + 1] & ~0x20) <= 'F') c += (10 + (str[i * 2 + 1] & ~0x20) - 'A');
        buf[i] = c;
    }
    return buf;
}

START_TEST(test_base58_decode)
{
    uint8_t addrhex[25] = {0};
    uint8_t signhex[65] = {0};
    char address[36] = "2EVNa4CK9SKosT4j1GEn8SuuUUEAXaHAMbM";
    char signature[90] = "GA82nXSwVEPV5soMjCiQkJb4oLEAo6FMK8CAE2n2YBTm7xjhAknUxtZrhs3RPVMfQsEoLwkJCEgvGj8a2vzthBQ1M";

    size_t sz = sizeof(signhex);
    b58tobin(signhex, &sz, signature);
    ck_assert_int_eq(sz, 65);
    ck_assert_mem_eq(signhex , fromhex("abc30130e2d9561fa8eb9871b75b13100689937dfc41c98d611b985ca25258c960be25c0b45874e1255f053863f6e175300d7e788d8b93d6dcfa9377120e4d3500"), sz);

    sz = sizeof(addrhex);
    b58tobin(addrhex, &sz, address);
    ck_assert_int_eq(sz, 25);
    ck_assert_mem_eq(addrhex , fromhex("b1aa8dd3e68d1d9b130c67ea1339ac9250b7d845002437a5a0"), sz);

}
END_TEST

START_TEST(test_generate_public_key_from_seckey)
{
    uint8_t seckey[32] = {0};
    uint8_t pubkey[33] = {0};

    memcpy(seckey, fromhex("a7e130694166cdb95b1e1bbce3f21e4dbd63f46df42b48c5a1f8295033d57d04"), sizeof(seckey));
    generate_pubkey_from_seckey(seckey, pubkey);
    ck_assert_mem_eq(pubkey, fromhex("0244350faa76799fec03de2f324acd077fd1b686c3a89babc0ef47096ccc5a13fa"), SHA256_DIGEST_LENGTH);

    memcpy(seckey, fromhex("c89b70a1f7b960c08068de9f2d3b32287833b26372935aa5042f7cc1dc985335"), sizeof(seckey));
    generate_pubkey_from_seckey(seckey, pubkey);
    ck_assert_mem_eq(pubkey, fromhex("03b17c7b7c564385be66f9c1b9da6a0b5aea56f0cb70548e6528a2f4f7b27245d8"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_generate_key_pair_from_seed)
{
    char seed[256] = "seed";
    uint8_t seckey[32] = {0};
    uint8_t pubkey[33] = {0};
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum((const uint8_t*)seed, digest, strlen(seed));
    generate_deterministic_key_pair(digest, SHA256_DIGEST_LENGTH, seckey, pubkey);
    ck_assert_mem_eq(seckey, fromhex("a7e130694166cdb95b1e1bbce3f21e4dbd63f46df42b48c5a1f8295033d57d04"), SHA256_DIGEST_LENGTH);
    ck_assert_mem_eq(pubkey, fromhex("0244350faa76799fec03de2f324acd077fd1b686c3a89babc0ef47096ccc5a13fa"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_secp256k1Hash)
{
    char seed[256] = "seed";
    uint8_t secp256k1Hash_digest[SHA256_DIGEST_LENGTH] = {0};
    secp256k1Hash((const uint8_t*)seed, strlen(seed), secp256k1Hash_digest);
    ck_assert_mem_eq(secp256k1Hash_digest, fromhex("c79454cf362b3f55e5effce09f664311650a44b9c189b3c8eed1ae9bd696cd9e"), SHA256_DIGEST_LENGTH);

    strcpy(seed, "random_seed");
    memset(secp256k1Hash_digest, 0, SHA256_DIGEST_LENGTH);
    secp256k1Hash((const uint8_t*)seed, strlen(seed), secp256k1Hash_digest);
    ck_assert_mem_eq(secp256k1Hash_digest, fromhex("5e81d46f56767496bc05ed177c5237cd4fe5013e617c726af43e1cba884f17d1"), SHA256_DIGEST_LENGTH);

    strcpy(seed, "random_seed");
    memset(secp256k1Hash_digest, 0, SHA256_DIGEST_LENGTH);
    secp256k1Hash((const uint8_t*)seed, strlen(seed), secp256k1Hash_digest);
    ck_assert_mem_eq(secp256k1Hash_digest, fromhex("5e81d46f56767496bc05ed177c5237cd4fe5013e617c726af43e1cba884f17d1"), SHA256_DIGEST_LENGTH);

}
END_TEST

START_TEST(test_generate_deterministic_key_pair_iterator)
{
    char seed[256] = "seed";
    uint8_t seckey[32] = {0};
    uint8_t pubkey[33] = {0};
    uint8_t nextSeed[SHA256_DIGEST_LENGTH] = {0};
    generate_deterministic_key_pair_iterator((const uint8_t*)seed, strlen(seed), nextSeed, seckey, pubkey);
    ck_assert_mem_eq(pubkey, fromhex("02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"), 33);
    ck_assert_mem_eq(seckey, fromhex("001aa9e416aff5f3a3c7f9ae0811757cf54f393d50df861f5c33747954341aa7"), 32);
    ck_assert_mem_eq(nextSeed, fromhex("c79454cf362b3f55e5effce09f664311650a44b9c189b3c8eed1ae9bd696cd9e"), 32);

    strcpy(seed, "random_seed");
    memset(pubkey, 0, sizeof(pubkey));
    memset(seckey, 0, sizeof(seckey));
    generate_deterministic_key_pair_iterator((const uint8_t*)seed, strlen(seed), nextSeed, seckey, pubkey);
    ck_assert_mem_eq(pubkey, fromhex("030e40dda21c27126d829b6ae57816e1440dcb2cc73e37e860af26eff1ec55ed73"), 33);
    ck_assert_mem_eq(seckey, fromhex("ff671860c58aad3f765d8add25046412dabf641186472e1553435e6e3c4a6fb0"), 32);
    ck_assert_mem_eq(nextSeed, fromhex("5e81d46f56767496bc05ed177c5237cd4fe5013e617c726af43e1cba884f17d1"), 32);

    strcpy(seed, "hello seed");
    memset(pubkey, 0, sizeof(pubkey));
    memset(seckey, 0, sizeof(seckey));
    generate_deterministic_key_pair_iterator((const uint8_t*)seed, strlen(seed), nextSeed, seckey, pubkey);
    ck_assert_mem_eq(pubkey, fromhex("035843e72258696b391cf1d898fc65f31e66876ea0c9e101f8ddc3ebb4b87dc5b0"), 33);
    ck_assert_mem_eq(seckey, fromhex("84fdc649964bf299a787cb78cd975910e197dbddd7db776ece544f41c44b3056"), 32);
    ck_assert_mem_eq(nextSeed, fromhex("70d382540812d4abc969dcc2adc66e805db96f7e1dcbe1ae6bbf2878211cbcf6"), 32);

    strcpy(seed, "skycoin5");
    memset(pubkey, 0, sizeof(pubkey));
    memset(seckey, 0, sizeof(seckey));
    generate_deterministic_key_pair_iterator((const uint8_t*)seed, strlen(seed), nextSeed, seckey, pubkey);
    ck_assert_mem_eq(pubkey, fromhex("03b17c7b7c564385be66f9c1b9da6a0b5aea56f0cb70548e6528a2f4f7b27245d8"), 33);
    ck_assert_mem_eq(seckey, fromhex("c89b70a1f7b960c08068de9f2d3b32287833b26372935aa5042f7cc1dc985335"), 32);
}
END_TEST

START_TEST(test_base58_address_from_pubkey)
{
    uint8_t pubkey[33] = {0};
    char address[256] = {0};
    size_t size_address = sizeof(address);
    memcpy(pubkey, fromhex("02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"), 33);
    generate_base58_address_from_pubkey(pubkey, address, &size_address);
    ck_assert_str_eq(address, "2EVNa4CK9SKosT4j1GEn8SuuUUEAXaHAMbM");

    memcpy(pubkey, fromhex("030e40dda21c27126d829b6ae57816e1440dcb2cc73e37e860af26eff1ec55ed73"), 33);
    generate_base58_address_from_pubkey(pubkey, address, &size_address);
    ck_assert_str_eq(address, "2EKq1QXRmfe7jsWzNdYsmyoz8q3VkwkLsDJ");

    memcpy(pubkey, fromhex("035843e72258696b391cf1d898fc65f31e66876ea0c9e101f8ddc3ebb4b87dc5b0"), 33);
    generate_base58_address_from_pubkey(pubkey, address, &size_address);
    ck_assert_str_eq(address, "5UgkXRHrf5XRk41BFq1DVyeFZHTQXirhUu");
}
END_TEST


START_TEST(test_bitcoin_address_from_pubkey)
{
    uint8_t pubkey[33] = {0};
    char address[256] = {0};
    size_t size_address = sizeof(address);
    memcpy(pubkey, fromhex("02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"), 33);
    generate_bitcoin_address_from_pubkey(pubkey, address, &size_address);
    ck_assert_str_eq(address, "1CN7JTzTTpmh1dsHeUSosXmNL2GLTwt78g");

    memcpy(pubkey, fromhex("030e40dda21c27126d829b6ae57816e1440dcb2cc73e37e860af26eff1ec55ed73"), 33);
    generate_bitcoin_address_from_pubkey(pubkey, address, &size_address);
    ck_assert_str_eq(address, "1DkKGd1YV9nhBKHWT9Aa2JzbEus98y6oU9");

    memcpy(pubkey, fromhex("035843e72258696b391cf1d898fc65f31e66876ea0c9e101f8ddc3ebb4b87dc5b0"), 33);
    generate_bitcoin_address_from_pubkey(pubkey, address, &size_address);
    ck_assert_str_eq(address, "1Ba2hpHH2o6H1NSrFpJTz5AbxdB2BdK5L2");
}
END_TEST


START_TEST(test_bitcoin_private_address_from_seckey)
{
    uint8_t seckey[32] = {0};
    char address[256] = {0};
    size_t size_address = sizeof(address);
    memcpy(seckey, fromhex("001aa9e416aff5f3a3c7f9ae0811757cf54f393d50df861f5c33747954341aa7"), 32);
    generate_bitcoin_private_address_from_seckey(seckey, address, &size_address);
    ck_assert_str_eq(address, "KwDuvkABDqb4WQiwc92DpXtBBiEywuKv46ZUvz5Gi5Xyn9gbcTJt");

    memcpy(seckey, fromhex("ff671860c58aad3f765d8add25046412dabf641186472e1553435e6e3c4a6fb0"), 32);
    generate_bitcoin_private_address_from_seckey(seckey, address, &size_address);
    ck_assert_str_eq(address, "L5nBR59QkW6kyXFvyqNbncWo2jPMoBXSH9fGUkh3n2RQn5Mj3vfY");

    memcpy(seckey, fromhex("84fdc649964bf299a787cb78cd975910e197dbddd7db776ece544f41c44b3056"), 32);
    generate_bitcoin_private_address_from_seckey(seckey, address, &size_address);
    ck_assert_str_eq(address, "L1gEDGuLTpMjybHnsJ24bUHhueocDrrKVdM3rj1rqXFHfyM2WtwD");
}
END_TEST

START_TEST(test_compute_sha256sum)
{
    char seed[256] = "seed";
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum((const uint8_t*)seed, digest, strlen(seed));

    ck_assert_mem_eq(digest, fromhex("19b25856e1c150ca834cffc8b59b23adbd0ec0389e58eb22b3b64768098d002b"), SHA256_DIGEST_LENGTH);

    strcpy(seed, "random_seed");
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    compute_sha256sum((const uint8_t*)seed, digest, strlen(seed));

    ck_assert_mem_eq(digest, fromhex("7b491face15c5be43df3affe42e6e4aab48522a3b564043de464e8de50184a5d"), SHA256_DIGEST_LENGTH);


    strcpy(seed, "024f7fd15da6c7fc7d0410d184073ef702104f82452da9b3e3792db01a8b7907c3");
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    compute_sha256sum((const uint8_t*)seed, digest, strlen(seed));

    ck_assert_mem_eq(digest, fromhex("a5daa8c9d03a9ec500088bdf0123a9d865725b03895b1291f25500737298e0a9"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_compute_ecdh)
{
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    uint8_t remote_pubkey[33];
    uint8_t my_seckey[32];

    memcpy(my_seckey, fromhex("8f609a12bdfc8572590c66763bb05ce609cc0fdcd0c563067e91c06bfd5f1027"), sizeof(my_seckey));
    memcpy(remote_pubkey, fromhex("03008fa0a5668a567cb28ab45e4b6747f5592690c1d519c860f748f6762fa13103"), sizeof(remote_pubkey));
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    ecdh_shared_secret(my_seckey, remote_pubkey, digest);
    ck_assert_mem_eq(digest, fromhex("907d3c524abb561a80644cdb0cf48e6c71ce33ed6a2d5eed40a771bcf86bd081"), SHA256_DIGEST_LENGTH);

    memcpy(my_seckey, fromhex("ec4c3702ae8dc5d3aaabc230d362f1ccc1ad2222353d006a057969bf2cc749c1"), sizeof(my_seckey));
    memcpy(remote_pubkey, fromhex("03b5d8432d20e55590b3e1e74a86f4689a5c1f5e25cc58840741fe1ac044d5e65c"), sizeof(remote_pubkey));
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    ecdh_shared_secret(my_seckey, remote_pubkey, digest);
    ck_assert_mem_eq(digest, fromhex("c59b456353d0fbceadc06d7794c42ebf413ab952b29ecf6052d30c7c1a50acda"), SHA256_DIGEST_LENGTH);

    memcpy(my_seckey, fromhex("19adca686f1ca7befc30af65765597a4d033ac7479850e79cef3ce5cb5b95da4"), sizeof(my_seckey));
    memcpy(remote_pubkey, fromhex("0328bd053c69d9c3dd1e864098e503de9839e990c63c48d8a4d6011c423658c4a9"), sizeof(remote_pubkey));
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    ecdh_shared_secret(my_seckey, remote_pubkey, digest);
    ck_assert_mem_eq(digest, fromhex("1fd2c655bcf19202ee004a3e0ae8f5c64ad1c0ce3b69f32ba18da188bb4d1eea"), SHA256_DIGEST_LENGTH);

    memcpy(my_seckey, fromhex("085d62c27a37889e02a183ee29962d5f4377831b4a70834ccea24a209e201404"), sizeof(my_seckey));
    memcpy(remote_pubkey, fromhex("030684d74471053ac6395ef74a86f88daa25f501329734c837c8c79c600423b220"), sizeof(remote_pubkey));
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    ecdh_shared_secret(my_seckey, remote_pubkey, digest);
    ck_assert_mem_eq(digest, fromhex("4225281b8498f05e0eaac02be79ce72471c2ddd8c127908b1f717bf64177b287"), SHA256_DIGEST_LENGTH);

    memcpy(my_seckey, fromhex("3c4289a9d884f74bd05c352fa1c08ce0d65955b59b24a572f46e02807dd42e62"), sizeof(my_seckey));
    memcpy(remote_pubkey, fromhex("0223496e9caa207e0f8cc283e970b85f2831732d5e0be2bcf9fa366f7e064a25dd"), sizeof(remote_pubkey));
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    ecdh_shared_secret(my_seckey, remote_pubkey, digest);
    ck_assert_mem_eq(digest, fromhex("70e5d568b31ed601fcb7f3144888d0633938817ae85417de1fbd0d52e29b5d7c"), SHA256_DIGEST_LENGTH);

    memcpy(my_seckey, fromhex("a7e130694166cdb95b1e1bbce3f21e4dbd63f46df42b48c5a1f8295033d57d04"), sizeof(my_seckey));
    memcpy(remote_pubkey, fromhex("02683e90daa5b0dd195b69e01386390284d3b3723121ce213771d9a0815d12b86c"), sizeof(remote_pubkey));
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    ecdh_shared_secret(my_seckey, remote_pubkey, digest);
    ck_assert_mem_eq(digest, fromhex("9ab65c0e99605712aac66be1eccccb6dacb867ebaf2b1ebf96d3d92524f247fd"), SHA256_DIGEST_LENGTH);
}
END_TEST


START_TEST(test_recover_pubkey_from_signed_message)
{
    int res;
    // uint8_t message[32];
    char message[256] = "Hello World!";
    uint8_t signature[65];
    uint8_t pubkey[33];
    // memcpy(message, fromhex("5dfbea13c81c48f7261994c148a7a39b9b51107d22b57bfd4613dce02dee46ee"), 32);
    memcpy(signature, fromhex("abc30130e2d9561fa8eb9871b75b13100689937dfc41c98d611b985ca25258c960be25c0b45874e1255f053863f6e175300d7e788d8b93d6dcfa9377120e4d3500"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"), 33);

    sprintf(message, "Hello World, it's me!");
    memcpy(signature, fromhex("54d7572cf5066225f349d89ad6d19e19e64d14711083f6607258b37407e5f0d26c6328d7c3ecb31eb4132f6b983f8ec33cdf3664c1df617526bbac140cdac75b01"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"), 33);

    memcpy(signature, fromhex("00b0dbb50c8b8f6c5be2bdee786a658a0ea22872ce90b21fbc0eb4f1d1018a043f93216a6af467acfb44aef9ab07e0a65621128504f3a61dfa0014b1cdd6d9c701"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"), 33);

    // sign with a different key pair
    // the seed for key pair generation was 'different'
    memcpy(signature, fromhex("5feef64dd9b9465e0f66ac21c5078cee4504f15ad407093b58908b69bc717d1c456901b4dbf9dde3eb170bd7aaf4e7a62f260e6194cc9884037affbfda250f3501"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    sprintf(message, "The seed was 'different'");
    memcpy(signature, fromhex("b8a91946af3cfe42139c853f09d1bc087db3bea0ab8bb20ab13790f4ba08aa4c327a4f614c61b2c532c2bab3852817ecd17b1c607f52f52c9c561ddbb2e4418e01"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    memcpy(signature, fromhex("f2e863beed0c026d0c631712dbe5ecb7ed95166275586271b77181ee3e68502b24c7a5c32b26ca5424fadfd8488285ad6e3ff3b86ed6c5449102d3198712f57b00"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    sprintf(message, "This msg has 24 letters.");
    memcpy(signature, fromhex("eff089c10e4c8d3c7244a8bc75d5657153ec7b42ed6d01bcc75cd08271a4aa7c19d1bd3b60330c909600238c1f18d99f06d2573c27cb4f2dfb0f65666a5a523200"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    sprintf(message, "This msg has 31 characters: ok!");
    memcpy(signature, fromhex("3dc77d17eeed0d3fd3d34ca05e8a9e84fbf73529b96bde7548080ac35d81470a60d5b8b37f2bb2500cf6a9745cd1c6edb81ebb5419e4f4fda9271794c8daf54200"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    // testing message maximal length
    sprintf(message, "This msg has 32 characters: max.");
    memcpy(signature, fromhex("e092ce21dda29349bd1e4e8b7a26d701542ac972b4e319a60bd887b6e51853622300e4e847f01a9aff4f51caa969759f717a6e5439b6bc4a5305b10bab9b5cb201"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    sprintf(message, "This msg has 32 characters: max..");
    memcpy(signature, fromhex("e092ce21dda29349bd1e4e8b7a26d701542ac972b4e319a60bd887b6e51853622300e4e847f01a9aff4f51caa969759f717a6e5439b6bc4a5305b10bab9b5cb201"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    sprintf(message, "This msg has 32 characters: max... What ever I write here is ignored.");
    memcpy(signature, fromhex("e092ce21dda29349bd1e4e8b7a26d701542ac972b4e319a60bd887b6e51853622300e4e847f01a9aff4f51caa969759f717a6e5439b6bc4a5305b10bab9b5cb201"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    memcpy(message, fromhex("176b81623cf98f45879f3a48fa34af77dde44b2ffa0ddd2bf9edb386f76ec0ef"), 32);
    memcpy(signature, fromhex("864c6abf85214be99fed3dc37591a74282f566fb52fb56ab21dabc0d120f29b848ffeb52a7843a49c411753c0edc12c0dedf6313266722bee982a0d3b384b62600"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("03b17c7b7c564385be66f9c1b9da6a0b5aea56f0cb70548e6528a2f4f7b27245d8"), 33);

    memcpy(message, fromhex("176b81623cf98f45879f3a48fa34af77dde44b2ffa0ddd2bf9edb386f76ec0ef"), 32);
    memcpy(signature, fromhex("631182b9722489eedd1a9eab36bf776c3e679aa2b1bd3fb346db0f776b982be25bdd33d4e893aca619eff3013e087307d22ca30644c96ea0fbdef06396d1bf9600"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("039f12c93645e35e5274dc38f191be0b6d1321ec35d2d2a3ddf7d13ed12f6da85b"), 33);

    memcpy(message, fromhex("176b81623cf98f45879f3a48fa34af77dde44b2ffa0ddd2bf9edb386f76ec0ef"), 32);
    memcpy(signature, fromhex("d2a8ec2b29ce3cf3e6048296188adff4b5dfcb337c1d1157f28654e445bb940b4e47d6b0c7ba43d072bf8618775f123a435e8d1a150cb39bbb1aa80da8c57ea100"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("03338ffc0ff42df07d27b0b4131cd96ffdfa4685b5566aafc7aa71ed10fd1cbd6f"), 33);

}
END_TEST

START_TEST(test_signature)
{	
    int res;
	uint8_t digest[32];
    uint8_t my_seckey[32];
    uint8_t signature[65];
    uint8_t pubkey[33];
    char* message = (char*)digest;
    memcpy(my_seckey, fromhex("597e27368656cab3c82bfcf2fb074cefd8b6101781a27709ba1b326b738d2c5a"), sizeof(my_seckey));
    memcpy(digest, fromhex("001aa9e416aff5f3a3c7f9ae0811757cf54f393d50df861f5c33747954341aa7"), 32);

    res = ecdsa_skycoin_sign(1, my_seckey, digest, signature);
    ck_assert_int_eq(res, 0);
	ck_assert_mem_eq(signature,  fromhex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"), 32);
	ck_assert_mem_eq(&signature[32],  fromhex("04641a7472bb90647fa60b4d30aef8c7279e4b68226f7b2713dab712ef122f8b01"), 32);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    res = ecdsa_skycoin_sign(0xfe25, my_seckey, digest, signature);
    ck_assert_int_eq(res, 0);
	ck_assert_mem_eq(signature,  fromhex("ee38f27be5f3c4b8db875c0ffbc0232e93f622d16ede888508a4920ab51c3c99"), 32);
	ck_assert_mem_eq(&signature[32],  fromhex("06ea7426c5e251e4bea76f06f554fa7798a49b7968b400fa981c51531a5748d801"), 32);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    res = ecdsa_skycoin_sign(0xfe250100, my_seckey, digest, signature);
    ck_assert_int_eq(res, 0);
	ck_assert_mem_eq(signature,  fromhex("d4d869ad39cb3a64fa1980b47d1f19bd568430d3f929e01c00f1e5b7c6840ba8"), 32);
	ck_assert_mem_eq(&signature[32],  fromhex("5e08d5781986ee72d1e8ebd4dd050386a64eee0256005626d2acbe3aefee9e2500"), 32);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132c7"), 33);

    // try of another key pair
    memcpy(my_seckey, fromhex("67a331669081d22624f16512ea61e1d44cb3f26af3333973d17e0e8d03733b78"), sizeof(my_seckey));
    
    res = ecdsa_skycoin_sign(0x1e2501ac, my_seckey, digest, signature);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(signature, fromhex("eeee743d79b40aaa52d9eeb48791b0ae81a2f425bf99cdbc84180e8ed429300d457e8d669dbff1716b123552baf6f6f0ef67f16c1d9ccd44e6785d424002212601"), 65);
    res = recover_pubkey_from_signed_message(message, signature, pubkey);
    ck_assert_int_eq(res, 0);
    ck_assert_mem_eq(pubkey,  fromhex("0270b763664593c5f84dfb20d23ef79530fc317e5ee2ece0d9c50f432f62426ff9"), 33);
}
END_TEST


START_TEST(test_checkdigest)
{
    ck_assert(is_digest("02df09821cff4874198a1dbdc462d224bd99728eeed024185879225762376132"));
    ck_assert(!is_digest("02df09821cff4874198a1dbdc462d224bd99728eeed0241858792257623761")); //too short
    ck_assert(!is_digest("02df09821cff4874198a1dbdc462d224bd99728eeed0241858792257623761256")); //too long
    ck_assert(!is_digest("02df09821cff4874198a1dbdc462d224bd99728eeed0241858792257623761r")); //non hex digits
}
END_TEST

// define test suite and cases
Suite *test_suite(void)
{
    Suite *s = suite_create("skycoin_crypto");
    TCase *tc;

    tc = tcase_create("checksums");
    tcase_add_test(tc, test_generate_public_key_from_seckey);
    tcase_add_test(tc, test_generate_key_pair_from_seed);
    tcase_add_test(tc, test_secp256k1Hash);
    tcase_add_test(tc, test_generate_deterministic_key_pair_iterator);
    tcase_add_test(tc, test_base58_address_from_pubkey);
    tcase_add_test(tc, test_bitcoin_address_from_pubkey);
    tcase_add_test(tc, test_bitcoin_private_address_from_seckey);
    tcase_add_test(tc, test_compute_sha256sum);
    tcase_add_test(tc, test_compute_ecdh);
    tcase_add_test(tc, test_recover_pubkey_from_signed_message);
    tcase_add_test(tc, test_base58_decode);
    tcase_add_test(tc, test_signature);
    tcase_add_test(tc, test_checkdigest);
    suite_add_tcase(s, tc);

    return s;
}


// run suite
int main(void)
{
    int number_failed;
    Suite *s = test_suite();
    SRunner *sr = srunner_create(s);
    srunner_run_all(sr, CK_VERBOSE);
    number_failed = srunner_ntests_failed(sr);
    srunner_free(sr);
    if (number_failed == 0) {
        printf("PASSED ALL TESTS\n");
    }
    return number_failed;
}
