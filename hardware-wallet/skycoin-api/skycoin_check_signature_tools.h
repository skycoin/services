#ifndef SKYCOIN_CHECK_SIGNATURE_TOOLS_H
#define SKYCOIN_CHECK_SIGNATURE_TOOLS_H

#include <stdint.h>
#include "bignum.h"
// #include "ecdsa.h"


// curve point x and y
typedef struct {
	bignum256 x, y;
} curve_point;


typedef struct {

	bignum256 prime;       // prime order of the finite field
	curve_point G;         // initial curve point
	bignum256 order;       // order of G
	bignum256 order_half;  // order of G divided by 2
	int       a;           // coefficient 'a' of the elliptic curve
	bignum256 b;           // coefficient 'b' of the elliptic curve

} ecdsa_curve;


typedef enum {
    HASHER_SHA2,
    HASHER_BLAKE,
} HasherType;


typedef struct {
	const char *bip32_name;    // string for generating BIP32 xprv from seed
	const ecdsa_curve *params; // ecdsa curve parameters, null for ed25519
	HasherType hasher_type;    // hasher type for BIP32 and ECDSA
} curve_info;


// typedef struct {
// 	uint32_t depth;
// 	uint32_t child_num;
// 	uint8_t chain_code[32];
// 	uint8_t private_key[32];
// 	uint8_t public_key[33];
// 	const curve_info *curve;
// } HDNode;


typedef struct jacobian_curve_point {
	bignum256 x, y, z;
} jacobian_curve_point;


typedef struct {
	uint32_t depth;
	uint32_t child_num;
	uint8_t chain_code[32];
	uint8_t private_key[32];
	uint8_t public_key[33];
	const curve_info *curve;
} HNode;

void ecdsa_get_public_key33(const ecdsa_curve *curve, const uint8_t *priv_key, uint8_t *pub_key);
// const curve_info *get_curve_by_name(const char *curve_name);
// int hdnode_from_seed(const uint8_t *seed, int seed_len, const char* curve, HDNode *out);
// void hdnode_fill_public_key(HDNode *node);
void create_node(const char* seed_str, HNode* node);
void uncompress_mcoords(const ecdsa_curve *curve, uint8_t odd, const bignum256 *x, bignum256 *y);
int mecdsa_validate_pubkey(const ecdsa_curve *curve, const curve_point *pub);
void mpoint_multiply(const ecdsa_curve *curve, const bignum256 *k, const curve_point *p, curve_point *res);
void mscalar_multiply(const ecdsa_curve *curve, const bignum256 *k, curve_point *res);
void mpoint_set_infinity(curve_point *p);
int mpoint_is_infinity(const curve_point *p);
void mpoint_add(const ecdsa_curve *curve, const curve_point *cp1, curve_point *cp2);
void mpoint_copy(const curve_point *cp1, curve_point *cp2);
int mpoint_is_equal(const curve_point *p, const curve_point *q);
int mpoint_is_negative_of(const curve_point *p, const curve_point *q);
void mpoint_double(const ecdsa_curve *curve, curve_point *cp);
void mcurve_to_jacobian(const curve_point *p, jacobian_curve_point *jp, const bignum256 *prime);
void mjacobian_to_curve(const jacobian_curve_point *jp, curve_point *p, const bignum256 *prime);
void mpoint_jacobian_add(const curve_point *p1, jacobian_curve_point *p2, const ecdsa_curve *curve);
void mconditional_negate(uint32_t cond, bignum256 *a, const bignum256 *prime);
void mpoint_jacobian_double(jacobian_curve_point *p, const ecdsa_curve *curve);

#endif