CC       ?= gcc

OPTFLAGS ?= -O3 -g -fPIC

CFLAGS   += $(OPTFLAGS) \
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
            -Werror

INC=-I$(CHECK_PATH)/src
INC+=-I$(CHECK_PATH)
INC+=-Itools/
CFLAGS += -I. $(INC)

SRCS += skycoin_crypto.c
SRCS += skycoin_check_signature.c
SRCS += skycoin_check_signature_tools.c
SRCS += $(shell ls tools/*.c)

OBJS   = $(SRCS:.c=.o)

%.o: %.c %.h
	$(CC) $(CFLAGS) -o $@ -c $<

TESTLIBS = -L/usr/local/lib/ -lm -lrt -L$(CHECK_PATH)/src -lcheck 

all: test_skycoin_crypto

skycoin-crypto: skycoin_crypto.o $(OBJS)
	$(CC) -rdynamic -shared $(CFLAGS) $(OBJS) -o libskycoin-crypto.so

libskycoin-crypto.a: skycoin_crypto.o $(OBJS)
	$(AR) rcs libskycoin-crypto.a skycoin_crypto.o $(OBJS)

test_skycoin_crypto: test_skycoin_crypto.o skycoin-crypto
	$(CC) test_skycoin_crypto.o $(OBJS) -L. -lskycoin-crypto  $(TESTLIBS) -o test_skycoin_crypto



clean:
	rm -f *.o test_skycoin_crypto
	rm -f *.so
	rm -f tools/*.o
	rm -f trezor-tools/*.o
	rm -f *.a
