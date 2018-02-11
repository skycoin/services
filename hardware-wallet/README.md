Tools to imitate the cipher library from [this repository](https://github.com/skycoin/skycoin/tree/develop/src/cipher).


# Compilation

This repository includes header files coming from [Trezor-crypto](https://github.com/trezor/trezor-crypto/) repository.

Download the repository

    git clone git@github.com:trezor/trezor-crypto.git

Then setup the TREZOR_CRYPTO_PATH environment variable:

    export TREZOR_CRYPTO_PATH=$PWD/trezor-crypto

Finally:

    make
    ./test_skycoin_crypto

# Dependencies

## Test library "check"

The source code necessary to compile libcheck.a can be downloaded from [this repository](https://github.com/libcheck/check)

    git clone git@github.com:libcheck/check.git 

## Recompile the static library libTrezorCrypto.a

The dependancy libTrezorCrypto.a can be recompiled from sources.

Add this line to the CFLAGS: "CFLAGS += -DUSE_BN_PRINT=1"
Then run :

    make 
    ar rcs libTrezorCrypto.a $(find -name "*.o")
