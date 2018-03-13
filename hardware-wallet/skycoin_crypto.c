#include "skycoin_crypto.h"

#include <string.h>
#include <stdio.h>

#include "sha2.h"
#include "bip32.h"
#include "curves.h"
#include "ripemd160.h"
#include "base58.h"

#include "ecdsa.h"
// #include "secp256k1.h"

extern void bn_print(const bignum256 *a);
void create_node(const char* seed_str, HDNode* node);

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


void generate_bitcoin_private_address_from_pubkey(const uint8_t* seckey, char* address, size_t *size_address)
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

void SetB32(const uint8_t* input, uint32_t* output)
{
    uint32_t v = 0;
    uint8_t i,j,limb,shift;
    for (i = 0; i < 32; i++) {
        for (j = 0; j < 4; j++) {
            limb = (8*i + 2*j) / 26;
            shift = (8*i + 2*j) % 26;
            v = (uint32_t)((input[31-i] >> (2*j)) & 0x3) << shift;
            output[limb] |= v;
        }
    }
}

// Compute public key from signature and recovery id.
// returns 0 if verification succeeded
int verify_digest_recover(const ecdsa_curve *curve, uint8_t *pub_key, const uint8_t *sig, const uint8_t *digest, int recid)
{
	bignum256 r, s, e;
	curve_point cp, cp2;
    uint32_t fx[10] = {0};
    // uint32_t ifx[10] = {0};

	// read r and s
	bn_read_be(sig, &r);
	bn_read_be(sig + 32, &s);
	if (!bn_is_less(&r, &curve->order) || bn_is_zero(&r)) {
		return 1;
	}
	if (!bn_is_less(&s, &curve->order) || bn_is_zero(&s)) {
		return 1;
	}
    uint8_t get_back_r[32] = {0};
    bn_write_be(&r, get_back_r);
    SetB32(get_back_r, fx);
    /*
    printf("fx= ");
    for (int i = 0; i < 10;++i)
    {
        printf("%x ", fx[i]);
        ifx[9-i] = fx[i];
    }
    printf("\n");
    printf("ifx= ");
    for (int i = 0; i < 10; ++i)
    {
        printf("%x ", ifx[i]);
    }
    printf("\n");
    */

	// cp = R = k * G (k is secret nonce when signing)
	if (recid & 2) {
		bn_add(&r, &curve->order);
		if (!bn_is_less(&r, &curve->prime)) {
			return 1;
		}
	}

    // bn_read_be((uint8_t*)fx, &cp.x);
    // memcpy(&cp.x, ifx, sizeof(bignum256));
	memcpy(&cp.x, &r, sizeof(bignum256));
    /*
    uint8_t get_back_x[32] = {0};
    bn_write_be(&cp.x, get_back_x);
    printf("get_back_x= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", get_back_x[i]);
    }
    printf("\n");
*/
	// compute y from x
	uncompress_coords(curve, recid & 1, &cp.x, &cp.y);
	if (!ecdsa_validate_pubkey(curve, &cp)) {
		return 1;
	}
    /*
    uint8_t get_back_y[32] = {0};
    bn_write_be(&cp.y, get_back_y);
    printf("get_back_y= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", get_back_y[i]);
    }
    printf("\n");
*/


	// r := r^-1
	bn_inverse(&r, &curve->order);
    /*
    uint8_t get_back_rn[288] = {0};
    bn_write_be(&r, get_back_rn);
    printf("get_back_rn= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", get_back_rn[i]);
    }
    printf("\n");


    uint8_t curve_order[288] = {0};
    bn_write_be(&curve->order, curve_order);
    printf("curve_order= ");
    for (int i = 0; i < 32;++i)
    {
        printf("%02x", curve_order[i]);
    }
    printf("\n");
*/
    uint8_t bn1[256] = {0};
	// e = -digest
	bn_read_be(digest, &e);
    while (! (uint8_t)e.val[0])
    {
        for (int i = 0; i < 8; ++i)
        {
            bn_rshift(&e);
        }
    }
    bn_write_be(&e, bn1);
    printf("e= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", bn1[i]);
    }
    printf("\n");

    bn_multiply(&r, &e, &curve->order);
    bn_write_be(&e, bn1);
    printf("bn1= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", bn1[i]);
    }
    printf("\n");
	bn_subtractmod(&curve->order, &e, &e, &curve->order);
    bn_fast_mod(&e, &curve->order);
	bn_mod(&e, &curve->order);
    bn_write_be(&e, bn1);
    printf("bn1= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", bn1[i]);
    }
    printf("\n");

    uint8_t bn2[33] = {0};
    bn_multiply(&r, &s, &curve->order);
    bn_write_be(&s, bn2);
    printf("bn2= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", bn2[i]);
    }
    printf("\n");

	// cp := s * R = s * k *G
	point_multiply(curve, &s, &cp, &cp);
	// cp2 := -digest * G
	scalar_multiply(curve, &e, &cp2);
	
    // cp := (s * k - digest) * G = (r*priv) * G = r * Pub
	point_add(curve, &cp2, &cp);
	// cp := r^{-1} * r * Pub = Pub
	// point_multiply(curve, &r, &cp, &cp);
	pub_key[0] = 0x04;
	bn_write_be(&cp.x, pub_key + 1);
	bn_write_be(&cp.y, pub_key + 33);
    
	return 0;
}

/*signature: 65 bytes, 
message 32 bytes, 
pubkey 33 bytes
returns 0 if signature matches and 5 if it does not*/
int recover_pubkey_from_signed_message(char* message, const uint8_t* signature, uint8_t* pubkey)
{
    int res = -1;
    HDNode dummy_node;
    char seed_str[256] = "dummy seed";
	uint8_t long_pubkey[65];
    create_node(seed_str, &dummy_node);
/*
    printf("message= ");
    for (int i = 0;i < 66;++i)
    {
        printf("%02x", (uint8_t)message[i]);
    }
    printf("\n");

    printf("sign= ");
    for (int i = 0;i < 65;++i)
    {
        printf("%02x", signature[i]);
    }
    printf("\n");
    */
	bignum256 r, s;
	bn_read_be(signature, &r);
	bn_read_be(signature + 32, &s);
    uint8_t get_back_r[288] = {0};
    uint8_t get_back_s[288] = {0};
    bn_write_be(&r, get_back_r);
    bn_write_be(&s, get_back_s);
/*
    printf("r= ");
    for (int i = 0;i < 9;++i)
    {
        printf("%08x", r.val[i]);
    }
    printf("\n");
    printf("get_back_r= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", get_back_r[i]);
    }
    printf("\n");
    printf("s= ");
    for (int i = 0;i < 9;++i)
    {
        printf("%08x", s.val[i]);
    }
    printf("\n");
    printf("get_back_s= ");
    for (int i = 0;i < 32;++i)
    {
        printf("%02x", get_back_s[i]);
    }
    printf("\n");

*/
    // var recid = int(sig[64])
    res = verify_digest_recover(dummy_node.curve->params, long_pubkey, signature, (uint8_t*)message, 0);
    // res = ecdsa_verify_digest(dummy_node.curve->params, long_pubkey, signature, (uint8_t*)message);
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

#if 0
int recover_pubkey_from_signed_message_bread_board(char* message, const uint8_t* signature, uint8_t* pubkey)
{
    secp256k1_pubkey long_pubkey;
    secp256k1_context *ctx = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
    
    
	bignum256 r, s;
	bn_read_be(signature, &r);
	bn_read_be(signature + 32, &s);

    secp256k1_ecdsa_signature sign;
    bn_write_be(&r, &sign.data[0]);
    bn_write_be(&s, &sign.data[32]);

    int res = secp256k1_ecdsa_verify(ctx, &sign, (uint8_t *)message, &long_pubkey);
    memcpy(&pubkey[1], &long_pubkey.data[1], 32);
    if (long_pubkey.data[63] % 2 == 0)
    {
        pubkey[0] = 0x02;
    }
    else
    {
        pubkey[0] = 0x03;
    }
    return res;
}
#endif