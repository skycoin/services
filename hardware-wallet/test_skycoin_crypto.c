#include "skycoin_crypto.h"

#include <stdio.h>
#include <string.h>

#include "check.h"

#include "ecdsa.h"
#include "curves.h"

#define FROMHEX_MAXLEN 512

void tohex(char * str, const uint8_t* buffer, int bufferLength)
{
	int i;
	for (i = 0; i < bufferLength; ++i)
	{
		sprintf(&str[2*i], "%02x", buffer[i]);
	}
}

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

START_TEST(test_generate_deterministic_key_pair_seckey)
{
    char seed[256] = "seed";
    uint8_t seckey_digest[SHA256_DIGEST_LENGTH] = {0};
    genereate_deterministic_key_pair_seckey(seed, seckey_digest);
	ck_assert_mem_eq(seckey_digest, fromhex("a7e130694166cdb95b1e1bbce3f21e4dbd63f46df42b48c5a1f8295033d57d04"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_compute_sha256sum)
{
    char seed[256] = "seed";
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum(seed, digest, strlen(seed));

	ck_assert_mem_eq(digest, fromhex("19b25856e1c150ca834cffc8b59b23adbd0ec0389e58eb22b3b64768098d002b"), SHA256_DIGEST_LENGTH);

    strcpy(seed, "random_seed");
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    compute_sha256sum(seed, digest, strlen(seed));

	ck_assert_mem_eq(digest, fromhex("7b491face15c5be43df3affe42e6e4aab48522a3b564043de464e8de50184a5d"), SHA256_DIGEST_LENGTH);


    strcpy(seed, "024f7fd15da6c7fc7d0410d184073ef702104f82452da9b3e3792db01a8b7907c3");
    memset(digest, 0, SHA256_DIGEST_LENGTH);
    compute_sha256sum(seed, digest, strlen(seed));

	ck_assert_mem_eq(digest, fromhex("a5daa8c9d03a9ec500088bdf0123a9d865725b03895b1291f25500737298e0a9"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_compute_ecdh)
{
    char seed_str[256] = "seed";
	HDNode alice;
    uint8_t session_key1[32] = {0};

    create_node(seed_str, &alice);
	ck_assert_mem_eq(alice.public_key, fromhex("03008fa0a5668a567cb28ab45e4b6747f5592690c1d519c860f748f6762fa13103"), 33);
	ck_assert_mem_eq(alice.private_key, fromhex("8f609a12bdfc8572590c66763bb05ce609cc0fdcd0c563067e91c06bfd5f1027"), 32);


	char pubkey[66] = {0};
	tohex(pubkey, alice.public_key, 33);
	printf("pub key: %s\n", pubkey);


    uint8_t mult[65] = {0};
	char key_m[128] = {0};
	ecdh_multiply(alice.curve->params, alice.private_key, alice.public_key, mult); //65

	tohex(key_m, mult, 65);
	printf("ECDH key_mult: %s\n", key_m);
	memcpy(&session_key1[1], &mult[1], 31);
	if (mult[64] % 2 == 0)
	{
		session_key1[0] = 0x02;
	}
	else
	{
		session_key1[0] = 0x03;
	}
	
	char key[64] = {0};
	tohex(key, session_key1, 32);
	printf("ECDH key: %s\n", key);

	ck_assert_mem_eq(session_key1, fromhex("024f7fd15da6c7fc7d0410d184073ef702104f82452da9b3e3792db01a8b7907c3"), 32);

    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum(key, digest, strlen(key));

	ck_assert_mem_eq(digest, fromhex("907d3c524abb561a80644cdb0cf48e6c71ce33ed6a2d5eed40a771bcf86bd081"), SHA256_DIGEST_LENGTH);
}
END_TEST

// define test suite and cases
Suite *test_suite(void)
{
	Suite *s = suite_create("skycoin_crypto");
	TCase *tc;

	tc = tcase_create("checksums");
	tcase_add_test(tc, test_generate_deterministic_key_pair_seckey);
	tcase_add_test(tc, test_compute_sha256sum);
	tcase_add_test(tc, test_compute_ecdh);
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
