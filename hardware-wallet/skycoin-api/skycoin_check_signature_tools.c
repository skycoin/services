#include "skycoin_check_signature_tools.h"

#include <string.h>
#include <assert.h>
#include "curves.h"
#include "memzero.h"
#include "hmac.h"
#include "rand.h"
// #include "secp256k1.h"


// generate random K for signing/side-channel noise
static void generate_k_random(bignum256 *k, const bignum256 *prime) {
	do {
		int i;
		for (i = 0; i < 8; i++) {
			k->val[i] = random32() & 0x3FFFFFFF;
		}
		k->val[8] = random32() & 0xFFFF;
		// check that k is in range and not zero.
	} while (bn_is_zero(k) || !bn_is_less(k, prime));
}

const ecdsa_curve msecp256k1 = {
	/* .prime */ {
		/*.val =*/ {0x3ffffc2f, 0x3ffffffb, 0x3fffffff, 0x3fffffff, 0x3fffffff, 0x3fffffff, 0x3fffffff, 0x3fffffff, 0xffff}
	},

	/* G */ {
		/*.x =*/{/*.val =*/{0x16f81798, 0x27ca056c, 0x1ce28d95, 0x26ff36cb, 0x70b0702, 0x18a573a, 0xbbac55a, 0x199fbe77, 0x79be}},
		/*.y =*/{/*.val =*/{0x3b10d4b8, 0x311f423f, 0x28554199, 0x5ed1229, 0x1108a8fd, 0x13eff038, 0x3c4655da, 0x369dc9a8, 0x483a}}
	},

	/* order */ {
		/*.val =*/{0x10364141, 0x3f497a33, 0x348a03bb, 0x2bb739ab, 0x3ffffeba, 0x3fffffff, 0x3fffffff, 0x3fffffff, 0xffff}
	},

	/* order_half */ {
		/*.val =*/{0x281b20a0, 0x3fa4bd19, 0x3a4501dd, 0x15db9cd5, 0x3fffff5d, 0x3fffffff, 0x3fffffff, 0x3fffffff, 0x7fff}
	},

	/* a */	0,

	/* b */ {
		/*.val =*/{7}
	}

// #if USE_PRECOMPUTED_CP
// 	,
// 	/* cp */ {
// #include "secp256k1.table"
// 	}
// #endif
};

const curve_info secp256k1_minfo = {
	.bip32_name = "Bitcoin seed",
	.params = &msecp256k1,
	.hasher_type = HASHER_SHA2,
};

static int from_seed(const uint8_t *seed, int seed_len, HNode *out)
{
	static CONFIDENTIAL uint8_t I[32 + 32];
	memset(out, 0, sizeof(HNode));
	out->depth = 0;
	out->child_num = 0;
	out->curve = &secp256k1_minfo;
	if (out->curve == 0) {
		return 0;
	}
	static CONFIDENTIAL HMAC_SHA512_CTX ctx;
	hmac_sha512_Init(&ctx, (const uint8_t*) out->curve->bip32_name, strlen(out->curve->bip32_name));
	hmac_sha512_Update(&ctx, seed, seed_len);
	hmac_sha512_Final(&ctx, I);

	if (out->curve->params) {
		bignum256 a;
		while (true) {
			bn_read_be(I, &a);
			if (!bn_is_zero(&a) // != 0
				&& bn_is_less(&a, &out->curve->params->order)) { // < order
				break;
			}
			hmac_sha512_Init(&ctx, (const uint8_t*) out->curve->bip32_name, strlen(out->curve->bip32_name));
			hmac_sha512_Update(&ctx, I, sizeof(I));
			hmac_sha512_Final(&ctx, I);
		}
		memzero(&a, sizeof(a));
	}
	memcpy(out->private_key, I, 32);
	memcpy(out->chain_code, I + 32, 32);
	memzero(out->public_key, sizeof(out->public_key));
	memzero(I, sizeof(I));
	return 1;
}


void mecdsa_get_public_key33(const ecdsa_curve *curve, const uint8_t *priv_key, uint8_t *pub_key)
{
	curve_point R;
	bignum256 k;

	bn_read_be(priv_key, &k);
	// compute k*G
	mscalar_multiply(curve, &k, &R);
	pub_key[0] = 0x02 | (R.y.val[0] & 0x01);
	bn_write_be(&R.x, pub_key + 1);
	memzero(&R, sizeof(R));
	memzero(&k, sizeof(k));
}

static void fill_public_key(HNode *node)
{
	if (node->public_key[0] != 0)
		return;
	if (node->curve->params) {
		mecdsa_get_public_key33(node->curve->params, node->private_key, node->public_key);
	} 
}

void create_node(const char* seed_str, HNode* node)
{
    from_seed((const uint8_t *)seed_str, strlen(seed_str), node);
    fill_public_key(node);
}


void uncompress_mcoords(const ecdsa_curve *curve, uint8_t odd, const bignum256 *x, bignum256 *y)
{
	// y^2 = x^3 + a*x + b
	memcpy(y, x, sizeof(bignum256));         // y is x
	bn_multiply(x, y, &curve->prime);        // y is x^2
	bn_subi(y, -curve->a, &curve->prime);    // y is x^2 + a
	bn_multiply(x, y, &curve->prime);        // y is x^3 + ax
	bn_add(y, &curve->b);                    // y is x^3 + ax + b
	bn_sqrt(y, &curve->prime);               // y = sqrt(y)
	if ((odd & 0x01) != (y->val[0] & 1)) {
		bn_subtract(&curve->prime, y, y);   // y = -y
	}
}

// Verifies that:
//   - pub is not the point at infinity.
//   - pub->x and pub->y are in range [0,p-1].
//   - pub is on the curve.

int mecdsa_validate_pubkey(const ecdsa_curve *curve, const curve_point *pub)
{
	bignum256 y_2, x3_ax_b;

	if (mpoint_is_infinity(pub)) {
		return 0;
	}

	if (!bn_is_less(&(pub->x), &curve->prime) || !bn_is_less(&(pub->y), &curve->prime)) {
		return 0;
	}

	memcpy(&y_2, &(pub->y), sizeof(bignum256));
	memcpy(&x3_ax_b, &(pub->x), sizeof(bignum256));

	// y^2
	bn_multiply(&(pub->y), &y_2, &curve->prime);
	bn_mod(&y_2, &curve->prime);

	// x^3 + ax + b
	bn_multiply(&(pub->x), &x3_ax_b, &curve->prime);  // x^2
	bn_subi(&x3_ax_b, -curve->a, &curve->prime);      // x^2 + a
	bn_multiply(&(pub->x), &x3_ax_b, &curve->prime);  // x^3 + ax
	bn_addmod(&x3_ax_b, &curve->b, &curve->prime);    // x^3 + ax + b
	bn_mod(&x3_ax_b, &curve->prime);

	if (!bn_is_equal(&x3_ax_b, &y_2)) {
		return 0;
	}

	return 1;
}


// res = k * p
void mpoint_multiply(const ecdsa_curve *curve, const bignum256 *k, const curve_point *p, curve_point *res)
{
	// this algorithm is loosely based on
	//  Katsuyuki Okeya and Tsuyoshi Takagi, The Width-w NAF Method Provides
	//  Small Memory and Fast Elliptic Scalar Multiplications Secure against
	//  Side Channel Attacks.
	assert (bn_is_less(k, &curve->order));

	int i, j;
	static CONFIDENTIAL bignum256 a;
	uint32_t *aptr;
	uint32_t abits;
	int ashift;
	uint32_t is_even = (k->val[0] & 1) - 1;
	uint32_t bits, sign, nsign;
	static CONFIDENTIAL jacobian_curve_point jres;
	curve_point pmult[8];
	const bignum256 *prime = &curve->prime;

	// is_even = 0xffffffff if k is even, 0 otherwise.

	// add 2^256.
	// make number odd: subtract curve->order if even
	uint32_t tmp = 1;
	uint32_t is_non_zero = 0;
	for (j = 0; j < 8; j++) {
		is_non_zero |= k->val[j];
		tmp += 0x3fffffff + k->val[j] - (curve->order.val[j] & is_even);
		a.val[j] = tmp & 0x3fffffff;
		tmp >>= 30;
	}
	is_non_zero |= k->val[j];
	a.val[j] = tmp + 0xffff + k->val[j] - (curve->order.val[j] & is_even);
	assert((a.val[0] & 1) != 0);

	// special case 0*p:  just return zero. We don't care about constant time.
	if (!is_non_zero) {
		mpoint_set_infinity(res);
		return;
	}

	// Now a = k + 2^256 (mod curve->order) and a is odd.
	//
	// The idea is to bring the new a into the form.
	// sum_{i=0..64} a[i] 16^i,  where |a[i]| < 16 and a[i] is odd.
	// a[0] is odd, since a is odd.  If a[i] would be even, we can
	// add 1 to it and subtract 16 from a[i-1].  Afterwards,
	// a[64] = 1, which is the 2^256 that we added before.
	//
	// Since k = a - 2^256 (mod curve->order), we can compute
	//   k*p = sum_{i=0..63} a[i] 16^i * p
	//
	// We compute |a[i]| * p in advance for all possible
	// values of |a[i]| * p.  pmult[i] = (2*i+1) * p
	// We compute p, 3*p, ..., 15*p and store it in the table pmult.
	// store p^2 temporarily in pmult[7]
	pmult[7] = *p;
	mpoint_double(curve, &pmult[7]);
	// compute 3*p, etc by repeatedly adding p^2.
	pmult[0] = *p;
	for (i = 1; i < 8; i++) {
		pmult[i] = pmult[7];
		mpoint_add(curve, &pmult[i-1], &pmult[i]);
	}

	// now compute  res = sum_{i=0..63} a[i] * 16^i * p step by step,
	// starting with i = 63.
	// initialize jres = |a[63]| * p.
	// Note that a[i] = a>>(4*i) & 0xf if (a&0x10) != 0
	// and - (16 - (a>>(4*i) & 0xf)) otherwise.   We can compute this as
	//   ((a ^ (((a >> 4) & 1) - 1)) & 0xf) >> 1
	// since a is odd.
	aptr = &a.val[8];
	abits = *aptr;
	ashift = 12;
	bits = abits >> ashift;
	sign = (bits >> 4) - 1;
	bits ^= sign;
	bits &= 15;
	mcurve_to_jacobian(&pmult[bits>>1], &jres, prime);
	for (i = 62; i >= 0; i--) {
		// sign = sign(a[i+1])  (0xffffffff for negative, 0 for positive)
		// invariant jres = (-1)^sign sum_{j=i+1..63} (a[j] * 16^{j-i-1} * p)
		// abits >> (ashift - 4) = lowbits(a >> (i*4))

		mpoint_jacobian_double(&jres, curve);
		mpoint_jacobian_double(&jres, curve);
		mpoint_jacobian_double(&jres, curve);
		mpoint_jacobian_double(&jres, curve);

		// get lowest 5 bits of a >> (i*4).
		ashift -= 4;
		if (ashift < 0) {
			// the condition only depends on the iteration number and
			// leaks no private information to a side-channel.
			bits = abits << (-ashift);
			abits = *(--aptr);
			ashift += 30;
			bits |= abits >> ashift;
		} else {
			bits = abits >> ashift;
		}
		bits &= 31;
		nsign = (bits >> 4) - 1;
		bits ^= nsign;
		bits &= 15;

		// negate last result to make signs of this round and the
		// last round equal.
		mconditional_negate(sign ^ nsign, &jres.z, prime);

		// add odd factor
		mpoint_jacobian_add(&pmult[bits >> 1], &jres, curve);
		sign = nsign;
	}
	mconditional_negate(sign, &jres.z, prime);
	mjacobian_to_curve(&jres, res, prime);
	memzero(&a, sizeof(a));
	memzero(&jres, sizeof(jres));
}

void mscalar_multiply(const ecdsa_curve *curve, const bignum256 *k, curve_point *res)
{
	mpoint_multiply(curve, k, &curve->G, res);
}


// set point to internal representation of point at infinity
void mpoint_set_infinity(curve_point *p)
{
	bn_zero(&(p->x));
	bn_zero(&(p->y));
}

// return true iff p represent point at infinity
// both coords are zero in internal representation
int mpoint_is_infinity(const curve_point *p)
{
	return bn_is_zero(&(p->x)) && bn_is_zero(&(p->y));
}


// cp2 = cp1 + cp2
void mpoint_add(const ecdsa_curve *curve, const curve_point *cp1, curve_point *cp2)
{
	bignum256 lambda, inv, xr, yr;

	if (mpoint_is_infinity(cp1)) {
		return;
	}
	if (mpoint_is_infinity(cp2)) {
		mpoint_copy(cp1, cp2);
		return;
	}
	if (mpoint_is_equal(cp1, cp2)) {
		mpoint_double(curve, cp2);
		return;
	}
	if (mpoint_is_negative_of(cp1, cp2)) {
		mpoint_set_infinity(cp2);
		return;
	}

	bn_subtractmod(&(cp2->x), &(cp1->x), &inv, &curve->prime);
	bn_inverse(&inv, &curve->prime);
	bn_subtractmod(&(cp2->y), &(cp1->y), &lambda, &curve->prime);
	bn_multiply(&inv, &lambda, &curve->prime);

	// xr = lambda^2 - x1 - x2
	xr = lambda;
	bn_multiply(&xr, &xr, &curve->prime);
	yr = cp1->x;
	bn_addmod(&yr, &(cp2->x), &curve->prime);
	bn_subtractmod(&xr, &yr, &xr, &curve->prime);
	bn_fast_mod(&xr, &curve->prime);
	bn_mod(&xr, &curve->prime);

	// yr = lambda (x1 - xr) - y1
	bn_subtractmod(&(cp1->x), &xr, &yr, &curve->prime);
	bn_multiply(&lambda, &yr, &curve->prime);
	bn_subtractmod(&yr, &(cp1->y), &yr, &curve->prime);
	bn_fast_mod(&yr, &curve->prime);
	bn_mod(&yr, &curve->prime);

	cp2->x = xr;
	cp2->y = yr;
}


// Set cp2 = cp1
void mpoint_copy(const curve_point *cp1, curve_point *cp2)
{
	*cp2 = *cp1;
}


// return true iff both points are equal
int mpoint_is_equal(const curve_point *p, const curve_point *q)
{
	return bn_is_equal(&(p->x), &(q->x)) && bn_is_equal(&(p->y), &(q->y));
}

// returns true iff p == -q
// expects p and q be valid points on curve other than point at infinity
int mpoint_is_negative_of(const curve_point *p, const curve_point *q)
{
	// if P == (x, y), then -P would be (x, -y) on this curve
	if (!bn_is_equal(&(p->x), &(q->x))) {
		return 0;
	}

	// we shouldn't hit this for a valid point
	if (bn_is_zero(&(p->y))) {
		return 0;
	}

	return !bn_is_equal(&(p->y), &(q->y));
}


// cp = cp + cp
void mpoint_double(const ecdsa_curve *curve, curve_point *cp)
{
	bignum256 lambda, xr, yr;

	if (mpoint_is_infinity(cp)) {
		return;
	}
	if (bn_is_zero(&(cp->y))) {
		mpoint_set_infinity(cp);
		return;
	}

	// lambda = (3 x^2 + a) / (2 y)
	lambda = cp->y;
	bn_mult_k(&lambda, 2, &curve->prime);
	bn_inverse(&lambda, &curve->prime);

	xr = cp->x;
	bn_multiply(&xr, &xr, &curve->prime);
	bn_mult_k(&xr, 3, &curve->prime);
	bn_subi(&xr, -curve->a, &curve->prime);
	bn_multiply(&xr, &lambda, &curve->prime);

	// xr = lambda^2 - 2*x
	xr = lambda;
	bn_multiply(&xr, &xr, &curve->prime);
	yr = cp->x;
	bn_lshift(&yr);
	bn_subtractmod(&xr, &yr, &xr, &curve->prime);
	bn_fast_mod(&xr, &curve->prime);
	bn_mod(&xr, &curve->prime);

	// yr = lambda (x - xr) - y
	bn_subtractmod(&(cp->x), &xr, &yr, &curve->prime);
	bn_multiply(&lambda, &yr, &curve->prime);
	bn_subtractmod(&yr, &(cp->y), &yr, &curve->prime);
	bn_fast_mod(&yr, &curve->prime);
	bn_mod(&yr, &curve->prime);

	cp->x = xr;
	cp->y = yr;
}

void mcurve_to_jacobian(const curve_point *p, jacobian_curve_point *jp, const bignum256 *prime) {
	// randomize z coordinate
	generate_k_random(&jp->z, prime);

	jp->x = jp->z;
	bn_multiply(&jp->z, &jp->x, prime);
	// x = z^2
	jp->y = jp->x;
	bn_multiply(&jp->z, &jp->y, prime);
	// y = z^3

	bn_multiply(&p->x, &jp->x, prime);
	bn_multiply(&p->y, &jp->y, prime);
}

void mjacobian_to_curve(const jacobian_curve_point *jp, curve_point *p, const bignum256 *prime) {
	p->y = jp->z;
	bn_inverse(&p->y, prime);
	// p->y = z^-1
	p->x = p->y;
	bn_multiply(&p->x, &p->x, prime);
	// p->x = z^-2
	bn_multiply(&p->x, &p->y, prime);
	// p->y = z^-3
	bn_multiply(&jp->x, &p->x, prime);
	// p->x = jp->x * z^-2
	bn_multiply(&jp->y, &p->y, prime);
	// p->y = jp->y * z^-3
	bn_mod(&p->x, prime);
	bn_mod(&p->y, prime);
}

void mpoint_jacobian_add(const curve_point *p1, jacobian_curve_point *p2, const ecdsa_curve *curve) {
	bignum256 r, h, r2;
	bignum256 hcby, hsqx;
	bignum256 xz, yz, az;
	int is_doubling;
	const bignum256 *prime = &curve->prime;
	int a = curve->a;

	assert (-3 <= a && a <= 0);

	/* First we bring p1 to the same denominator:
	 * x1' := x1 * z2^2
	 * y1' := y1 * z2^3
	 */
	/*
	 * lambda  = ((y1' - y2)/z2^3) / ((x1' - x2)/z2^2)
	 *         = (y1' - y2) / (x1' - x2) z2
	 * x3/z3^2 = lambda^2 - (x1' + x2)/z2^2
	 * y3/z3^3 = 1/2 lambda * (2x3/z3^2 - (x1' + x2)/z2^2) + (y1'+y2)/z2^3
	 *
	 * For the special case x1=x2, y1=y2 (doubling) we have
	 * lambda = 3/2 ((x2/z2^2)^2 + a) / (y2/z2^3)
	 *        = 3/2 (x2^2 + a*z2^4) / y2*z2)
	 *
	 * to get rid of fraction we write lambda as
	 * lambda = r / (h*z2)
	 * with  r = is_doubling ? 3/2 x2^2 + az2^4 : (y1 - y2)
	 *       h = is_doubling ?      y1+y2       : (x1 - x2)
	 *
	 * With z3 = h*z2  (the denominator of lambda)
	 * we get x3 = lambda^2*z3^2 - (x1' + x2)/z2^2*z3^2
	 *           = r^2 - h^2 * (x1' + x2)
	 *    and y3 = 1/2 r * (2x3 - h^2*(x1' + x2)) + h^3*(y1' + y2)
	 */

	/* h = x1 - x2
	 * r = y1 - y2
	 * x3 = r^2 - h^3 - 2*h^2*x2
	 * y3 = r*(h^2*x2 - x3) - h^3*y2
	 * z3 = h*z2
	 */

	xz = p2->z;
	bn_multiply(&xz, &xz, prime); // xz = z2^2
	yz = p2->z;
	bn_multiply(&xz, &yz, prime); // yz = z2^3
	
	if (a != 0) {
		az  = xz;
		bn_multiply(&az, &az, prime);   // az = z2^4
		bn_mult_k(&az, -a, prime);      // az = -az2^4
	}
	
	bn_multiply(&p1->x, &xz, prime);        // xz = x1' = x1*z2^2;
	h = xz;
	bn_subtractmod(&h, &p2->x, &h, prime);
	bn_fast_mod(&h, prime);
	// h = x1' - x2;

	bn_add(&xz, &p2->x);
	// xz = x1' + x2

	// check for h == 0 % prime.  Note that h never normalizes to
	// zero, since h = x1' + 2*prime - x2 > 0 and a positive
	// multiple of prime is always normalized to prime by
	// bn_fast_mod.
	is_doubling = bn_is_equal(&h, prime);

	bn_multiply(&p1->y, &yz, prime);        // yz = y1' = y1*z2^3;
	bn_subtractmod(&yz, &p2->y, &r, prime);
	// r = y1' - y2;

	bn_add(&yz, &p2->y);
	// yz = y1' + y2

	r2 = p2->x;
	bn_multiply(&r2, &r2, prime);
	bn_mult_k(&r2, 3, prime);
	
	if (a != 0) {
		// subtract -a z2^4, i.e, add a z2^4
		bn_subtractmod(&r2, &az, &r2, prime);
	}
	bn_cmov(&r, is_doubling, &r2, &r);
	bn_cmov(&h, is_doubling, &yz, &h);
	

	// hsqx = h^2
	hsqx = h;
	bn_multiply(&hsqx, &hsqx, prime);

	// hcby = h^3
	hcby = h;
	bn_multiply(&hsqx, &hcby, prime);

	// hsqx = h^2 * (x1 + x2)
	bn_multiply(&xz, &hsqx, prime);

	// hcby = h^3 * (y1 + y2)
	bn_multiply(&yz, &hcby, prime);

	// z3 = h*z2
	bn_multiply(&h, &p2->z, prime);

	// x3 = r^2 - h^2 (x1 + x2)
	p2->x = r;
	bn_multiply(&p2->x, &p2->x, prime);
	bn_subtractmod(&p2->x, &hsqx, &p2->x, prime);
	bn_fast_mod(&p2->x, prime);

	// y3 = 1/2 (r*(h^2 (x1 + x2) - 2x3) - h^3 (y1 + y2))
	bn_subtractmod(&hsqx, &p2->x, &p2->y, prime);
	bn_subtractmod(&p2->y, &p2->x, &p2->y, prime);
	bn_multiply(&r, &p2->y, prime);
	bn_subtractmod(&p2->y, &hcby, &p2->y, prime);
	bn_mult_half(&p2->y, prime);
	bn_fast_mod(&p2->y, prime);
}

// Negate a (modulo prime) if cond is 0xffffffff, keep it if cond is 0.
// The timing of this function does not depend on cond.
void mconditional_negate(uint32_t cond, bignum256 *a, const bignum256 *prime)
{
	int j;
	uint32_t tmp = 1;
	assert(a->val[8] < 0x20000);
	for (j = 0; j < 8; j++) {
		tmp += 0x3fffffff + 2*prime->val[j] - a->val[j];
		a->val[j] = ((tmp & 0x3fffffff) & cond) | (a->val[j] & ~cond);
		tmp >>= 30;
	}
	tmp += 0x3fffffff + 2*prime->val[j] - a->val[j];
	a->val[j] = ((tmp & 0x3fffffff) & cond) | (a->val[j] & ~cond);
	assert(a->val[8] < 0x20000);
}


void mpoint_jacobian_double(jacobian_curve_point *p, const ecdsa_curve *curve) {
	bignum256 az4, m, msq, ysq, xysq;
	const bignum256 *prime = &curve->prime;

	assert (-3 <= curve->a && curve->a <= 0);
	/* usual algorithm:
	 *
	 * lambda  = (3((x/z^2)^2 + a) / 2y/z^3) = (3x^2 + az^4)/2yz
	 * x3/z3^2 = lambda^2 - 2x/z^2
	 * y3/z3^3 = lambda * (x/z^2 - x3/z3^2) - y/z^3
	 *
	 * to get rid of fraction we set
	 *  m = (3 x^2 + az^4) / 2
	 * Hence,
	 *  lambda = m / yz = m / z3
	 *
	 * With z3 = yz  (the denominator of lambda)
	 * we get x3 = lambda^2*z3^2 - 2*x/z^2*z3^2
	 *           = m^2 - 2*xy^2
	 *    and y3 = (lambda * (x/z^2 - x3/z3^2) - y/z^3) * z3^3
	 *           = m * (xy^2 - x3) - y^4
	 */

	/* m = (3*x^2 + a z^4) / 2
	 * x3 = m^2 - 2*xy^2
	 * y3 = m*(xy^2 - x3) - 8y^4
	 * z3 = y*z
	 */

	m = p->x;
	bn_multiply(&m, &m, prime);
	bn_mult_k(&m, 3, prime);

	az4 = p->z;
	bn_multiply(&az4, &az4, prime);
	bn_multiply(&az4, &az4, prime);
	bn_mult_k(&az4, curve->a, prime);
	bn_subtractmod(&m, &az4, &m, prime);
	bn_mult_half(&m, prime);

	// msq = m^2
	msq = m;
	bn_multiply(&msq, &msq, prime);
	// ysq = y^2
	ysq = p->y;
	bn_multiply(&ysq, &ysq, prime);
	// xysq = xy^2
	xysq = p->x;
	bn_multiply(&ysq, &xysq, prime);

	// z3 = yz
	bn_multiply(&p->y, &p->z, prime);

	// x3 = m^2 - 2*xy^2
	p->x = xysq;
	bn_lshift(&p->x);
	bn_fast_mod(&p->x, prime);
	bn_subtractmod(&msq, &p->x, &p->x, prime);
	bn_fast_mod(&p->x, prime);

	// y3 = m*(xy^2 - x3) - y^4
	bn_subtractmod(&xysq, &p->x, &p->y, prime);
	bn_multiply(&m, &p->y, prime);
	bn_multiply(&ysq, &ysq, prime);
	bn_subtractmod(&p->y, &ysq, &p->y, prime);
	bn_fast_mod(&p->y, prime);
}
