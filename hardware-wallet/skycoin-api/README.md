# How to compile and run tests

This code aims at tranforming the cipher library from [this repository](https://github.com/skycoin/skycoin/tree/develop/src/cipher) into firmware for the STM32 hardware.

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