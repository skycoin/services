NAME  = bootloader

OBJS += ../vendor/skycoin-crypto/skycoin_crypto.o
#trezor-crypto
OBJS += $(TOP_DIR)vendor/trezor-crypto/address.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/curves.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/nist256p1.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/rand.o

OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/curve25519-donna-32bit.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/curve25519-donna-helpers.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/modm-donna-32bit.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/ed25519-donna-basepoint-table.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/ed25519-donna-32bit-tables.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/ed25519-donna-impl-base.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/ed25519.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/curve25519-donna-scalarmult-base.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/ed25519-sha3.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ed25519-donna/ed25519-keccak.o

OBJS += $(TOP_DIR)vendor/trezor-crypto/hmac.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/bip32.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/bip39.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/pbkdf2.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/base32.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/base58.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/segwit_addr.o

OBJS += $(TOP_DIR)vendor/trezor-crypto/ripemd160.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/sha3.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/blake256.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/hasher.o

OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aescrypt.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aeskey.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aestab.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aes_modes.o
#end trezor-crypto
#libtrezor.a

OBJS += ../startup.o

OBJS += ../buttons.o
OBJS += ../layout.o
OBJS += ../oled.o
OBJS += ../rng.o
OBJS += ../serialno.o

OBJS += ../setup.o

OBJS += ../util.o
OBJS += ../memory.o

OBJS += ../timer.o

OBJS += ../gen/bitmaps.o
OBJS += ../gen/fonts.o

#end libtrezor.a

OBJS += bootloader.o
OBJS += signatures.o
OBJS += usb.o

OBJS += ../vendor/trezor-crypto/bignum.small.o
OBJS += ../vendor/trezor-crypto/ecdsa.small.o
OBJS += ../vendor/trezor-crypto/secp256k1.small.o
OBJS += ../vendor/trezor-crypto/sha2.small.o
OBJS += ../vendor/trezor-crypto/memzero.small.o

CFLAGS += -DUSE_PRECOMPUTED_IV=0
CFLAGS += -DUSE_PRECOMPUTED_CP=0

ifeq ($(SIGNATURE_PROTECT),1)
CFLAGS += -DSIGNATURE_PROTECT=1
else
CFLAGS += -DSIGNATURE_PROTECT=0
endif
OPTFLAGS ?= -Os

LDSCRIPT = bootloader-memory.ld
LDLIBS   += -lopencm3_stm32f2

include ../Makefile.include


align: $(NAME).bin
	./firmware_align.py $(NAME).bin
