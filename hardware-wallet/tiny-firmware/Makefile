# imported from Makefile.include
TOP_DIR       := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))../firmware/
TOOLCHAIN_DIR ?= $(TOP_DIR)vendor/libopencm3
APPVER = 1.0.0

PREFIX   ?= arm-none-eabi-
CC       := $(PREFIX)gcc
LD       := $(PREFIX)gcc
OBJCOPY  := $(PREFIX)objcopy
OBJDUMP  := $(PREFIX)objdump
AR       := $(PREFIX)ar
AS       := $(PREFIX)as
OPENOCD  := openocd -f interface/stlink-v2.cfg -c "transport select hla_swd" -f target/stm32f2x.cfg

OPTFLAGS ?= -O3
DBGFLAGS ?= -g -DNDEBUG
CPUFLAGS ?= -mcpu=cortex-m3 -mthumb
FPUFLAGS ?= -msoft-float

CFLAGS   += $(OPTFLAGS) \
            $(DBGFLAGS) \
            -std=gnu99 \
            -W \
            -Wall \
            -Wextra \
            -Wimplicit-function-declaration \
            -Wredundant-decls \
            -Wstrict-prototypes \
            -Wundef \
            -Wshadow \
            -Wpointer-arith \
            -Wformat \
            -Wreturn-type \
            -Wsign-compare \
            -Wmultichar \
            -Wformat-nonliteral \
            -Winit-self \
            -Wuninitialized \
            -Wformat-security \
            -Werror \
            -fno-common \
            -fno-exceptions \
            -fvisibility=internal \
            -ffunction-sections \
            -fdata-sections \
            -fstack-protector-all \
            $(CPUFLAGS) \
            $(FPUFLAGS) \
            -DSTM32F2 \
            -DCONFIDENTIAL='__attribute__((section("confidential")))' \
            -DRAND_PLATFORM_INDEPENDENT=1 \
            -I$(TOOLCHAIN_DIR)/include \
            -I$(TOP_DIR) \
            -I$(TOP_DIR)gen \
            -I$(TOP_DIR)vendor/trezor-crypto \
            -I$(TOP_DIR)vendor/trezor-qrenc \
            -I$(TOP_DIR)vendor/skycoin-crypto

LDFLAGS  += -L$(TOP_DIR) \
            $(DBGFLAGS) \
            $(CPUFLAGS) \
            $(FPUFLAGS)



LDSCRIPT  = memory.ld

LDFLAGS  += --static \
            -Wl,--start-group \
            -lc \
            -lgcc \
            -lnosys \
            -Wl,--end-group \
            -L$(TOOLCHAIN_DIR)/lib \
            -T$(LDSCRIPT) \
            -nostartfiles \
            -Wl,--gc-sections

LDLIBS   += -ltrezor
LIBDEPS  += $(TOP_DIR)/libtrezor.a

LDLIBS   += -lopencm3_stm32f2
LIBDEPS  += $(TOOLCHAIN_DIR)/lib/libopencm3_stm32f2.a

CFLAGS += -DFASTFLASH=0
CFLAGS += -DEMULATOR=0
CFLAGS += -DMEMORY_PROTECT=0

# end Makefile.include

CFLAGS += -Wno-sequence-point
CFLAGS += -I$(TOP_DIR)vendor/nanopb -Iprotob -DPB_FIELD_16BIT=1

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
OBJS += firmware/recovery.o
OBJS += firmware/protect.o
OBJS += firmware/pinmatrix.o
OBJS += firmware/reset.o
OBJS += firmware/fastflash.o

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

all: proto $(NAME).bin

$(NAME).bin: $(NAME).elf
	$(OBJCOPY) -Obinary $(NAME).elf $(NAME).bin

$(NAME).elf: $(OBJS) $(LDSCRIPT) $(LIBDEPS)
	$(LD) -o $(NAME).elf $(OBJS) $(LDLIBS) $(LDFLAGS)

%.o: %.c Makefile
	$(CC) $(CFLAGS) -MMD -MP -o $@ -c $<

%.d: %.c Makefile
	@$(CC) $(CFLAGS) -MM -MP -MG -o $@ $<

proto:
	cd protob && make

clean:
	rm -f *.elf
	rm -f *.bin
	rm -f *.d
	rm -f firmware/*.d
	rm -f $(OBJS)
	cd protob && make clean