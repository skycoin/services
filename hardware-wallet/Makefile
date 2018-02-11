CC       ?= gcc

OPTFLAGS ?= -O3 -g

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


INC=-I$(TREZOR_CRYPTO_PATH)
CFLAGS += -I. $(INC)

SRCS  += skycoin_crypto.c

OBJS   = $(SRCS:.c=.o)

TESTLIBS = -lcheck -L$(TREZOR_CRYPTO_PATH) -lTrezorCrypto

all: test_skycoin_crypto


test_skycoin_crypto: test_skycoin_crypto.o skycoin_crypto.o $(OBJS)
	$(CC) test_skycoin_crypto.o $(OBJS) $(TESTLIBS) -o test_skycoin_crypto



clean:
	rm -f *.o test_skycoin_crypto
