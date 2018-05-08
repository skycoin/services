NAME  = bootloader

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
OPTFLAGS ?= -Os

LDSCRIPT = bootloader-memory.ld
LDLIBS   += -lopencm3_stm32f2

include ../Makefile.include


align: $(NAME).bin
	./firmware_align.py $(NAME).bin
