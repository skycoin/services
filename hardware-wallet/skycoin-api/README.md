
This code aims at tranforming the cipher library from [this repository](https://github.com/skycoin/skycoin/tree/develop/src/cipher) into firmware for the STM32 hardware.

# How to compile and run tests 

## Trezor-crypto

This repository includes header files coming from [Trezor-crypto](https://github.com/trezor/trezor-crypto/) repository.

These files are however available on this repository as a dependency of the firmware [here](https://github.com/skycoin/services/tree/hardware-wallet/hardware-wallet/firmware/vendor/trezor-crypto).

Setup the TREZOR_CRYPTO_PATH environment variable:

    export TREZOR_CRYPTO_PATH=$PWD/services/hardware-wallet/firmware/vendor/trezor-crypto
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$TREZOR_CRYPTO_PATH


Warning: If the files were compiled for arm target you may need to clean them.

Add this line to the CFLAGS of this [Makefile](https://github.com/skycoin/services/blob/hardware-wallet/hardware-wallet/firmware/vendor/trezor-crypto/Makefile): "CFLAGS += -DUSE_BN_PRINT=1" 
    
    cd $TREZOR_CRYPTO_PATH
    make clean
    make


If you prefer compiling a static library libTrezorCrypto.a from sources:

Then run :

    make 
    ar rcs libTrezorCrypto.a $(find -name "*.o")

## Check

The source code necessary to compile libcheck.a can be downloaded from [this repository](https://github.com/libcheck/check)
Download the repository

    git clone git@github.com:libcheck/check.git

Then setup the TREZOR_CRYPTO_PATH environment variable:

    export CHECK_PATH=$PWD/check
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$CHECK_PATH

## Make and run !

    make
    ./test_skycoin_crypto