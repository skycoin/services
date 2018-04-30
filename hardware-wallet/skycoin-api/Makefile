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


INC=-I$(TREZOR_CRYPTO_PATH) -I$(CHECK_PATH)/src
CFLAGS += -I. $(INC)

SRCS  += skycoin_crypto.c

OBJS   = $(SRCS:.c=.o)

%.o: %.c %.h
	$(CC) $(CFLAGS) -o $@ -c $<

TESTLIBS = -L/usr/local/lib/ -lm -lrt -L$(CHECK_PATH)/src -lcheck  
CRYPTOLIBS = -L$(TREZOR_CRYPTO_PATH) -ltrezor-crypto

all: test_skycoin_crypto

skycoin-crypto: skycoin_crypto.o $(OBJS)
	$(CC) -rdynamic -shared $(CFLAGS) $(OBJS) $(CRYPTOLIBS) -o libskycoin-crypto.so

test_skycoin_crypto: test_skycoin_crypto.o skycoin-crypto
	$(CC) test_skycoin_crypto.o $(OBJS) -L. -lskycoin-crypto  $(TESTLIBS) $(CRYPTOLIBS) -o test_skycoin_crypto



clean:
	rm -f *.o test_skycoin_crypto
	rm -f *.so
