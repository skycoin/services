NAME  = bootloader

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

include ../Makefile.include


LDLIBS   += -L$(TOP_DIR)../firmware/ -ltrezor
LDSCRIPT  = bootloader-memory.ld
LDLIBS   += -lopencm3_stm32f2

align: $(NAME).bin
	./firmware_align.py $(NAME).bin

clean:
	rm -f *.elf
	rm -f *.bin
	rm -f *.d
	rm -f *.o
