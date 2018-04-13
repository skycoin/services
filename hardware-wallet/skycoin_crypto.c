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

void create_node(const char* seed_str, HDNode* node)
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


void generate_pubkey_from_seckey(const uint8_t* seckey, uint8_t* pubkey)
{
    char seed_str[256] = "dummy seed";
    HDNode dummy_node;
    create_node(seed_str, &dummy_node);
    ecdsa_get_public_key33(dummy_node.curve->params, seckey, pubkey);
}

void generate_deterministic_key_pair(const uint8_t* seed, const size_t seed_length, uint8_t* seckey, uint8_t* pubkey)
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
    generate_deterministic_key_pair(hash2, SHA256_DIGEST_LENGTH, dummy_seckey, pubkey);
    ecdh(seckey, pubkey, ecdh_key);
    memcpy(secp256k1Hash, hash, sizeof(hash));
    memcpy(&secp256k1Hash[SHA256_DIGEST_LENGTH], ecdh_key, sizeof(ecdh_key));
    compute_sha256sum((const char *)secp256k1Hash, secp256k1Hash_digest, sizeof(secp256k1Hash));
}

void generate_deterministic_key_pair_iterator(const char* seed, uint8_t* seckey, uint8_t* pubkey)
{
    size_t seed_length = 0;
    uint8_t seed1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t seed2[SHA256_DIGEST_LENGTH] = {0};
    char keypair_seed[256] = {0};
    secp256k1Hash(seed, seed1);
    seed_length = strlen(seed);
    memcpy(keypair_seed, seed, seed_length);
    memcpy(&keypair_seed[seed_length], seed1, sizeof(seed1));
    compute_sha256sum(keypair_seed, seed2, seed_length + sizeof(seed1));
    generate_deterministic_key_pair(seed2, SHA256_DIGEST_LENGTH, seckey, pubkey);
}

void compute_sha256sum(const char *seed, uint8_t* digest /*size SHA256_DIGEST_LENGTH*/, size_t seed_lenght)
{
    SHA256_CTX ctx;
    sha256_Init(&ctx);
    sha256_Update(&ctx, (const uint8_t*) seed, seed_lenght);
    sha256_Final(&ctx, digest);
}

// address_size is the size of the allocated address buffer, it will be overwritten by the computed address size
void generate_base58_address_from_pubkey(const uint8_t* pubkey, char* address, size_t *size_address)
{
    uint8_t pubkey_hash[25] = {0};
    uint8_t r1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t r2[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum((char *)pubkey, r1, 33);
    compute_sha256sum((char *)r1, r2, sizeof(r1));
    ripemd160(r2, SHA256_DIGEST_LENGTH, pubkey_hash);
    // compute base58 address
    uint8_t digest[SHA256_DIGEST_LENGTH] = {0};
    pubkey_hash[20] = 0;
    compute_sha256sum((char *)pubkey_hash, digest, 21);
    memcpy(&pubkey_hash[21], digest, 4);
    b58enc(address, size_address, pubkey_hash, sizeof(pubkey_hash));
}

void generate_bitcoin_address_from_pubkey(const uint8_t* pubkey, char* address, size_t *size_address)
{
    uint8_t b1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t b2[25] = {0};
    uint8_t h1[SHA256_DIGEST_LENGTH] = {0};
    uint8_t b4[SHA256_DIGEST_LENGTH] = {0};
    compute_sha256sum((char *)pubkey, b1, 33);
    ripemd160(b1, SHA256_DIGEST_LENGTH, &b2[1]);
    compute_sha256sum((char *)b2, h1, 21);
    compute_sha256sum((char *)h1, b4, SHA256_DIGEST_LENGTH);
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
    compute_sha256sum((char *)b2, h1, 34);
    compute_sha256sum((char *)h1, b3, SHA256_DIGEST_LENGTH);
    memcpy(&b2[34], b3, 4);
    b58enc(address, size_address, b2, sizeof(b2));
}

// Compute public key from signature and recovery id.
// returns 0 if verification succeeded
int verify_digest_recover(const ecdsa_curve *curve, uint8_t *pub_key, const uint8_t *sig, const uint8_t *digest)
{
	bignum256 r, s, e;
	curve_point cp, cp2;

	// read r and s
	bn_read_be(sig, &r);
	bn_read_be(sig + 32, &s);
	if (!bn_is_less(&r, &curve->order) || bn_is_zero(&r)) {
		return 1;
	}
	if (!bn_is_less(&s, &curve->order) || bn_is_zero(&s)) {
		return 1;
	}
    uint8_t recid = sig[64];

	// cp = R = k * G (k is secret nonce when signing)
	if (recid & 2) {
		bn_add(&r, &curve->order);
		if (!bn_is_less(&r, &curve->prime)) {
			return 1;
		}
	}

	memcpy(&cp.x, &r, sizeof(bignum256));
	// compute y from x
	uncompress_coords(curve, recid & 1, &cp.x, &cp.y);
	if (!ecdsa_validate_pubkey(curve, &cp)) {
		return 1;
	}
	// r := r^-1
	bn_inverse(&r, &curve->order);
    
	// e = -digest
	bn_read_be(digest, &e);
    while (! (uint8_t)e.val[0])
    {
        for (int i = 0; i < 8; ++i)
        {
            bn_rshift(&e);
        }
    }

    bn_multiply(&r, &e, &curve->order);
	bn_subtractmod(&curve->order, &e, &e, &curve->order);
    bn_fast_mod(&e, &curve->order);
	bn_mod(&e, &curve->order);
   
    bn_multiply(&r, &s, &curve->order);

	// cp := s * R = s * k *G
	point_multiply(curve, &s, &cp, &cp);
	// cp2 := -digest * G
	scalar_multiply(curve, &e, &cp2);
	
    // cp := (s * k - digest) * G = (r*priv) * G = r * Pub
	point_add(curve, &cp2, &cp);
	pub_key[0] = 0x04;
	bn_write_be(&cp.x, pub_key + 1);
	bn_write_be(&cp.y, pub_key + 33);
    
	return 0;
}

/*signature: 65 bytes, 
message 32 bytes, 
pubkey 33 bytes
returns 0 if signature matches and 5 if it does not*/
int recover_pubkey_from_signed_message(const char* message, const uint8_t* signature, uint8_t* pubkey)
{
    int res = -1;
    HDNode dummy_node;
    char seed_str[256] = "dummy seed";
	uint8_t long_pubkey[65];
    create_node(seed_str, &dummy_node);

	bignum256 r, s;
	bn_read_be(signature, &r);
	bn_read_be(signature + 32, &s);

    res = verify_digest_recover(dummy_node.curve->params, long_pubkey, signature, (uint8_t*)message);
    memcpy(&pubkey[1], &long_pubkey[1], 32);
    if (long_pubkey[64] % 2 == 0)
    {
        pubkey[0] = 0x02;
    }
    else
    {
        pubkey[0] = 0x03;
    }
    return res;
}

// uses secp256k1 curve
// priv_key is a 32 byte big endian stored number
// sig is 65 bytes long array for the signature
// digest is 32 bytes of digest
// is_canonical is an optional function that checks if the signature
// conforms to additional coin-specific rules.
int ecdsa_skycoin_sign(const uint32_t nonce_value, const uint8_t *priv_key, const uint8_t *digest, uint8_t *sig, uint8_t *pby)
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
//////////////////////
        printf("R.x: ");
        bn_print(&R.x);
        printf("\n");
//////////////////////
        printf("R.y: ");
        bn_print(&R.y);
        printf("\n");
//////////////////////

		// randomize operations to counter side-channel attacks
		// generate_k_random(&randk, &dummy_node.curve->params->order);
        // bn_one(&randk);
		// bn_multiply(&randk, &nonce, &dummy_node.curve->params->order); // nonce*rand
		bn_inverse(&nonce, &dummy_node.curve->params->order);         // (nonce*rand)^-1
		bn_read_be(priv_key, s);               // priv
//////////////////////
        printf("s: ");
        bn_print(s);
        printf("\n");
//////////////////////
		bn_multiply(&R.x, s, &dummy_node.curve->params->order);   // R.x*priv
//////////////////////
        printf("s: ");
        bn_print(s);
        printf("\n");
//////////////////////
        bn_add(s, &z);                         // R.x*priv + z
//////////////////////
        printf("z: ");
        bn_print(&z);
        printf("\n");
//////////////////////
		bn_multiply(&nonce, s, &dummy_node.curve->params->order);     // (nonce*rand)^-1 (R.x*priv + z)
		// bn_multiply(&randk, s, &dummy_node.curve->params->order);  // nonce^-1 (R.x*priv + z)
		bn_mod(s, &dummy_node.curve->params->order);
//////////////////////
        printf("s: ");
        bn_print(s);
        printf("\n");
//////////////////////
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

		if (pby) {
			*pby = by;
		}
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
