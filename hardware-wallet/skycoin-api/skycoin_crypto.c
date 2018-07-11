#include "skycoin_crypto.h"

#include <string.h>
#include <stdio.h> //sprintf

#include "sha2.h"
#include "bip32.h"
#include "curves.h"
#include "ripemd160.h"
#include "base58.h"
#include "ecdsa.h"

extern void bn_print(const bignum256 *a);
static void create_node(const char* seed_str, HDNode* node);

static void create_node(const char* seed_str, HDNode* node)
{
    const char* curve_name = SECP256K1_NAME; 
    hdnode_from_seed((const uint8_t *)seed_str, strlen(seed_str), curve_name, node);
    hdnode_fill_public_key(node);
}

void tohex(char * str, const uint8_t* buffer, int bufferLength)
{
    int i;
    for (i = 0; i < bufferLength; ++i)
    {
        sprintf(&str[2*i], "%02x", buffer[i]);
    }
}

void writebuf_fromhexstr(const char *str, uint8_t* buf)
{
    size_t len = strlen(str) / 2;
    if (len > 512) len = 512;
    for (size_t i = 0; i < len; i++) {
        uint8_t c = 0;
        if (str[i * 2] >= '0' && str[i*2] <= '9') c += (str[i * 2] - '0') << 4;
        if ((str[i * 2] & ~0x20) >= 'A' && (str[i*2] & ~0x20) <= 'F') c += (10 + (str[i * 2] & ~0x20) - 'A') << 4;
        if (str[i * 2 + 1] >= '0' && str[i * 2 + 1] <= '9') c += (str[i * 2 + 1] - '0');
        if ((str[i * 2 + 1] & ~0x20) >= 'A' && (str[i * 2 + 1] & ~0x20) <= 'F') c += (10 + (str[i * 2 + 1] & ~0x20) - 'A');
        buf[i] = c;
    }
}

void generate_pubkey_from_seckey(const uint8_t* seckey, uint8_t* pubkey)
{
    char seed_str[256] = "dummy seed";
    HDNode dummy_node;
    create_node(seed_str, &dummy_node);
    ecdsa_get_public_key33(dummy_node.curve->params, seckey, pubkey);
}

void generate_deterministic_key_pair(const uint8_t* seed, const size_t seed_length, uint8_t* seckey, uint8_t* pubkey)
{
    compute_sha256sum(seed, seckey, seed_length);
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
    compute_sha256sum(ecdh_key, shared_secret, 33);
}

void secp256k1Hash(const uint8_t* seed, const size_t seed_length, uint8_t* secp256k1Hash_digest)
{
    uint8_t seckey[32] = {0};
    uint8_t dummy_seckey[32] = {0};
    uint8_t pubkey[33] = {0};
    uint8_t hash[SHA256_DIGEST_LENGTH] = {0};
    uint8_t hash2[SHA256_DIGEST_LENGTH] = {0};
    uint8_t ecdh_key[33] = {0};
    uint8_t secp256k1Hash[SHA256_DIGEST_LENGTH + 33] = {0};
    compute_sha256sum(seed, hash, seed_length);
    compute_sha256sum(hash, seckey, sizeof(hash));
    compute_sha256sum(hash, hash2, sizeof(hash));
    generate_deterministic_key_pair(hash2, SHA256_DIGEST_LENGTH, dummy_seckey, pubkey);
    ecdh(seckey, pubkey, ecdh_key);
    memcpy(secp256k1Hash, hash, sizeof(hash));
    memcpy(&secp256k1Hash[SHA256_DIGEST_LENGTH], ecdh_key, sizeof(ecdh_key));
    compute_sha256sum(secp256k1Hash, secp256k1Hash_digest, sizeof(secp256k1Hash));
}

// nextSeed should be 32 bytes (size of a secp256k1Hash digest)
void generate_deterministic_key_pair_iterator(const uint8_t* seed, const size_t seed_length, uint8_t* nextSeed, uint8_t* seckey, uint8_t* pubkey)
{
    uint8_t seed1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t seed2[SHA256_DIGEST_LENGTH] = {0};
    uint8_t keypair_seed[256] = {0};
    secp256k1Hash(seed, seed_length, seed1);
    memcpy(keypair_seed, seed, seed_length);
    memcpy(&keypair_seed[seed_length], seed1, SHA256_DIGEST_LENGTH);
    memcpy(nextSeed, seed1, SHA256_DIGEST_LENGTH);
    compute_sha256sum(keypair_seed, seed2, seed_length + sizeof(seed1));
    generate_deterministic_key_pair(seed2, SHA256_DIGEST_LENGTH, seckey, pubkey);
}

void compute_sha256sum(const uint8_t *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght)
{
    SHA256_CTX ctx;
    sha256_Init(&ctx);
    sha256_Update(&ctx, seed, seed_lenght);
    sha256_Final(&ctx, digest);
}

// address_size is the size of the allocated address buffer, it will be overwritten by the computed address size
void generate_base58_address_from_pubkey(const uint8_t* pubkey, char* address, size_t *size_address)
{
    uint8_t pubkey_hash[25] = {0};
    uint8_t r1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t r2[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum(pubkey, r1, 33);
    compute_sha256sum(r1, r2, sizeof(r1));
    ripemd160(r2, SHA256_DIGEST_LENGTH, pubkey_hash);
    // compute base58 address
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    pubkey_hash[20] = 0;
    compute_sha256sum(pubkey_hash, digest, 21);
    memcpy(&pubkey_hash[21], digest, 4);
    b58enc(address, size_address, pubkey_hash, sizeof(pubkey_hash));
}

void generate_bitcoin_address_from_pubkey(const uint8_t* pubkey, char* address, size_t *size_address)
{
    uint8_t b1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t b2[25] = {0};
    uint8_t h1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t b4[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum(pubkey, b1, 33);
    ripemd160(b1, SHA256_DIGEST_LENGTH, &b2[1]);
    compute_sha256sum(b2, h1, 21);
    compute_sha256sum(h1, b4, SHA256_DIGEST_LENGTH);
    memcpy(&b2[21], b4, 4);
    b58enc(address, size_address, b2, sizeof(b2));
}


void generate_bitcoin_private_address_from_seckey(const uint8_t* seckey, char* address, size_t *size_address)
{
    uint8_t b2[38] = {0};
    uint8_t h1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t b3[SHA256_DIGEST_LENGTH] = {0};
    memcpy(&b2[1], seckey, 32);
    b2[0] = 0x80;
    b2[33] = 0x01;
    compute_sha256sum(b2, h1, 34);
    compute_sha256sum(h1, b3, SHA256_DIGEST_LENGTH);
    memcpy(&b2[34], b3, 4);
    b58enc(address, size_address, b2, sizeof(b2));
}

// uses secp256k1 curve
// priv_key is a 32 byte big endian stored number
// sig is 65 bytes long array for the signature
// digest is 32 bytes of digest
// is_canonical is an optional function that checks if the signature
// conforms to additional coin-specific rules.
int ecdsa_skycoin_sign(const uint32_t nonce_value, const uint8_t *priv_key, const uint8_t *digest, uint8_t *sig)
{
	int i;
	curve_point R;
	bignum256 nonce, z, randk;
	bignum256 *s = &R.y;
	uint8_t by; // signature recovery byte

    HDNode dummy_node;
    char seed_str[256] = "dummy seed";
    create_node(seed_str, &dummy_node);
	bn_read_be(digest, &z);

	for (i = 0; i < 1; i++) {
		// generate random number nonce
		// generate_k_random(&nonce, &dummy_node.curve->params->order);
		bn_read_uint32(nonce_value, &nonce);
		// compute nonce*G
		scalar_multiply(dummy_node.curve->params, &nonce, &R);
		by = R.y.val[0] & 1;
		// r = (rx mod n)
		if (!bn_is_less(&R.x, &dummy_node.curve->params->order)) {
			bn_subtract(&R.x, &dummy_node.curve->params->order, &R.x);
			by |= 2;
		}
		// if r is zero, we retry
		if (bn_is_zero(&R.x)) {
            printf("Premature exit 1");
			continue;
		}
		bn_inverse(&nonce, &dummy_node.curve->params->order);         // (nonce*rand)^-1
		bn_read_be(priv_key, s);               // priv
		bn_multiply(&R.x, s, &dummy_node.curve->params->order);   // R.x*priv
        bn_add(s, &z);                         // R.x*priv + z
		bn_multiply(&nonce, s, &dummy_node.curve->params->order);     // (nonce*rand)^-1 (R.x*priv + z)
		bn_mod(s, &dummy_node.curve->params->order);

		// if s is zero, we retry
		if (bn_is_zero(s)) {
            printf("Premature exit 2");
			continue;
		}

		// if S > order/2 => S = -S
		if (bn_is_less(&dummy_node.curve->params->order_half, s)) {
			bn_subtract(&dummy_node.curve->params->order, s, s);
			by ^= 1;
		}
		// we are done, R.x and s is the result signature
		bn_write_be(&R.x, sig);
		bn_write_be(s, sig + 32);

        sig[64] = by;

		memset(&nonce, 0, sizeof(nonce));
		memset(&randk, 0, sizeof(randk));

		return 0;
	}

	// Too many retries without a valid signature
	// -> fail with an error
	memset(&nonce, 0, sizeof(nonce));
	memset(&randk, 0, sizeof(randk));

	return -1;
}
