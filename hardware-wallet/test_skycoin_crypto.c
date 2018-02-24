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

START_TEST(test_generate_public_key_from_seckey)
{
	uint8_t seckey[32] = {0};
	uint8_t pubkey[33] = {0};

	memcpy(seckey, fromhex("a7e130694166cdb95b1e1bbce3f21e4dbd63f46df42b48c5a1f8295033d57d04"), sizeof(seckey));
    generate_pubkey_from_seckey(seckey, pubkey);

	ck_assert_mem_eq(pubkey, fromhex("0244350faa76799fec03de2f324acd077fd1b686c3a89babc0ef47096ccc5a13fa"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_generate_key_pair_from_seed)
{
    char seed[256] = "seed";
    uint8_t seckey[32] = {0};
    uint8_t pubkey[33] = {0};
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum((char *)seed, digest, strlen(seed));
    genereate_deterministic_key_pair(digest, seckey, pubkey);
    ck_assert_mem_eq(seckey, fromhex("a7e130694166cdb95b1e1bbce3f21e4dbd63f46df42b48c5a1f8295033d57d04"), SHA256_DIGEST_LENGTH);
    ck_assert_mem_eq(pubkey, fromhex("0244350faa76799fec03de2f324acd077fd1b686c3a89babc0ef47096ccc5a13fa"), SHA256_DIGEST_LENGTH);
}
END_TEST

START_TEST(test_secp256k1Hash)
{
    char seed[256] = "seed";
    uint8_t secp256k1Hash_digest[SHA256_DIGEST_LENGTH] = {0};
	secp256k1Hash(seed, secp256k1Hash_digest);
    ck_assert_mem_eq(secp256k1Hash_digest, fromhex("c79454cf362b3f55e5effce09f664311650a44b9c189b3c8eed1ae9bd696cd9e"), SHA256_DIGEST_LENGTH);
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

// define test suite and cases
Suite *test_suite(void)
{
	Suite *s = suite_create("skycoin_crypto");
	TCase *tc;

	tc = tcase_create("checksums");
	tcase_add_test(tc, test_generate_public_key_from_seckey);
	tcase_add_test(tc, test_generate_key_pair_from_seed);
    tcase_add_test(tc, test_secp256k1Hash);
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
