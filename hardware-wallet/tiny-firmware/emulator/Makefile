EMULATOR := 1

CFLAGS += -DEMULATOR=1
CFLAGS += -Wno-pointer-to-int-cast
CFLAGS += -Wno-int-to-pointer-cast

HEADLESS ?= 0

CFLAGS += -include stdio.h
CFLAGS += -include $(TOP_DIR)emulator/emulator.h

ifeq ($(HEADLESS),1)
CFLAGS   += -DHEADLESS=1
else
CFLAGS   += -DHEADLESS=0
CFLAGS   += -I/usr/include/SDL2 -D_REENTRANT
LDLIBS   += -lSDL2
endif

OBJS += setup.o

OBJS += buttons.o
OBJS += flash.o
OBJS += oled.o
OBJS += rng.o
OBJS += timer.o
OBJS += udp.o

OBJS += strl.o

libemulator.a: $(OBJS)
	$(AR) rcs $@ $(OBJS)

include ../Makefile.include
