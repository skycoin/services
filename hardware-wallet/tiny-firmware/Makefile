APPVER = 1.0.0
CFLAGS += -DAPPVER=$(APPVER)

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
LDLIBS   += -lopencm3_stm32f2
LIBDEPS  += $(TOOLCHAIN_DIR)/lib/libopencm3_stm32f2.a

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

#skycoin-crypto
OBJS += $(TOP_DIR)vendor/skycoin-crypto/skycoin_crypto.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/base58.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/blake256.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/hasher.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/rand.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/sha2.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/bignum.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/curves.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/hmac.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/ripemd160.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/bip32.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/ecdsa.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/memzero.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/secp256k1.o
#end skycoin-crypto

#trezor-crypto
OBJS += $(TOP_DIR)vendor/trezor-crypto/bip39.o
OBJS += $(TOP_DIR)vendor/trezor-crypto/pbkdf2.o
#end trezor-crypto

NAME  = skycoin

.PHONY: proto libopencm3

all: libopencm3 proto $(NAME).bin

proto: 
	cd protob && make

libopencm3: 
	cd vendor/libopencm3 && make

sign: $(NAME).bin
	bootloader/firmware_sign.py -f $(NAME).bin

include Makefile.include

clean::
	rm -f $(OBJS)
	rm -f *.a
	rm -f *.bin
	rm -f *.d
	rm -f *.elf
	rm -f *.hex
	rm -f *.list
	rm -f *.log
	rm -f *.srec
	cd protob && make clean