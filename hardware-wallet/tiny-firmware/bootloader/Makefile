NAME  = bootloader

#libtrezor.a
OBJS += ../startup.o
OBJS += ../buttons.o
OBJS += ../layout.o
OBJS += ../oled.o
OBJS += ../rng.o
OBJS += ../serialno.o
OBJS += ../setup.o
OBJS += ../supervise.o
OBJS += ../timer.o
OBJS += ../util.o
OBJS += ../memory.o
OBJS += ../gen/bitmaps.o
OBJS += ../gen/fonts.o
#end libtrezor.a

#skycoin-crypto
OBJS += $(TOP_DIR)vendor/skycoin-crypto/skycoin_check_signature.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/skycoin_check_signature_tools.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/curves.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/hmac.o
# OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/bip32.o

#end skycoin-crypto

OBJS += bootloader.o
OBJS += signatures.o
OBJS += usb.o

OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/bignum.small.o
# OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/ecdsa.small.o
# OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/secp256k1.small.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/sha2.small.o
OBJS += $(TOP_DIR)vendor/skycoin-crypto/tools/memzero.small.o

CFLAGS += -DUSE_PRECOMPUTED_IV=0
CFLAGS += -DUSE_PRECOMPUTED_CP=0

ifeq ($(SIGNATURE_PROTECT),1)
CFLAGS += -DSIGNATURE_PROTECT=1
else
CFLAGS += -DSIGNATURE_PROTECT=0
endif

ifeq ($(SIGNATURE_DEBUG),1)
CFLAGS += -DSIGNATURE_DEBUG=1
else
CFLAGS += -DSIGNATURE_DEBUG=0
endif
OPTFLAGS ?= -Os

LDSCRIPT = bootloader-memory.ld
LDLIBS   += -lopencm3_stm32f2

include ../Makefile.include


align: $(NAME).bin
	./firmware_align.py $(NAME).bin
