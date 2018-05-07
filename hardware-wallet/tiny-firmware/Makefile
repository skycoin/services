#libtrezor.a

OBJS += startup.o

OBJS += buttons.o
OBJS += layout.o
OBJS += oled.o
OBJS += rng.o
OBJS += serialno.o

OBJS += setup.o

OBJS += util.o
OBJS += memory.o

OBJS += timer.o

OBJS += gen/bitmaps.o
OBJS += gen/fonts.o

#end libtrezor.a
LDSCRIPT  = memory.ld

DEBUG_LINK ?= 0
CFLAGS += -Wno-sequence-point
CFLAGS += -I$(TOP_DIR)vendor/nanopb -Iprotob -DPB_FIELD_16BIT=1
CFLAGS += -DQR_MAX_VERSION=0
CFLAGS += -DDEBUG_LINK=$(DEBUG_LINK)
CFLAGS += -DDEBUG_LOG=$(DEBUG_LOG)

INC+=-Ifirmware
CFLAGS += -I. $(INC)

OBJS += $(TOP_DIR)vendor/nanopb/pb_common.o
OBJS += $(TOP_DIR)vendor/nanopb/pb_decode.o
OBJS += $(TOP_DIR)vendor/nanopb/pb_encode.o

OBJS += firmware/usb.o
OBJS += firmware/trezor.o
OBJS += firmware/layout2.o
OBJS += firmware/storage.o
OBJS += firmware/messages.o
OBJS += firmware/fsm.o
OBJS += firmware/protect.o

OBJS += protob/messages.pb.o
OBJS += protob/types.pb.o


OBJS += $(TOP_DIR)vendor/skycoin-crypto/skycoin_crypto.o
#trezor-crypto
OBJS += $(TOP_DIR)vendor/trezor-crypto/address.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/bignum.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/ecdsa.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/curves.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/secp256k1.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/nist256p1.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/rand.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/memzero.o

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
OBJS += $(TOP_DIR)vendor/trezor-crypto/sha2.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/sha3.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/blake256.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/hasher.o

OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aescrypt.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aeskey.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aestab.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/aes/aes_modes.o
#end trezor-crypto


NAME  = skycoin

.PHONY: proto libopencm3

all: libopencm3 proto $(NAME).bin

include Makefile.include

proto: 
	cd protob && make

libopencm3: 
	cd vendor/libopencm3 && make

clean:
	rm -f *.elf
	rm -f *.bin
	rm -f *.d
	rm -f gen/*.d
	rm -f firmware/*.d
	rm -f $(OBJS)
	cd protob && make clean