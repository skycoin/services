# imported from Makefile.include
TOP_DIR       := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))../firmware/
TOOLCHAIN_DIR ?= $(TOP_DIR)vendor/libopencm3

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

CFLAGS+= -DFASTFLASH=0
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
OBJS += protob/messages.pb.o
OBJS += protob/types.pb.o

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