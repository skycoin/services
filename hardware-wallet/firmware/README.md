This repository had been copied and modified from [trezor-mcu](https://github.com/trezor/trezor-mcu). Please refer to the [README.md file](https://github.com/trezor/trezor-mcu/blob/master/README.md) on that repository for extra details about bootloader and firmware compilation.

Even though this repo is pushed here it is (for now) intended to be a submodule of [trezor-mcu repository](https://github.com/trezor/trezor-mcu).

Trezor-mcu contains firmware and bootloader code example for STM32 hardware.

This code aims at tranforming the cipher library from [this repository](https://github.com/skycoin/skycoin/tree/develop/src/cipher) into firmware for the STM32 hardware.


# How to compile for firmware

## 1. Clone trezor-mcu from this github account

    git clone https://github.com/mpsido/trezor-skycoin.git
    cd trezor-skycoin

You should be on a branch called skycoin

## 2. Prepare environment

#### Download GNU ARM Embedded Toolchain

    wget https://developer.arm.com/-/media/Files/downloads/gnu-rm/6-2017q2/gcc-arm-none-eabi-6-2017-q2-update-linux.tar.bz2

#### Extract crosscompiler and add it to your path

    tar xjf gcc-arm-none-eabi-6-2017-q2-update-linux.tar.bz2
    export PATH="gcc-arm-none-eabi-6-2017-q2-update/bin:$PATH"

## 3. Compile the sources

The steps are the following, you can create a make_firmware.sh file with this content.

    #this option is important, it prevents the bootloader from locking the firmware in case its signature is wrong
    export MEMORY_PROTECT=0 
    make -C vendor/libopencm3/
    make -C vendor/nanopb/generator/proto/
    make -C firmware/protob/
    make -C vendor/skycoin-crypto/
    make
    #Merge the bootloader and the firmware into one file
    make -C bootloader/ align
    make -C firmware/ sign
    cp bootloader/bootloader.bin bootloader/combine/bl.bin
    cp firmware/trezor.bin bootloader/combine/fw.bin
    pushd bootloader/combine/ && ./prepare.py
    popd;

The output binary file is combined.bin located in bootloader/combine


# How to burn the firmware in the device

## 1. Install ST-LINK

Follow the steps [here](https://github.com/texane/stlink/blob/master/doc/compiling.md).

## 2. Use ST-LINK to burn the device

You can check the device is seen by your ST-LINK using this command:

    st-info --probe

To flash the device on a microcontroller of STM32f2xx family the command is:

    st-flash write combined.bin 0x08000000; 

# Communicate with the device using the command line

## 1. Download the python-trezor repository

    cd ..
    git clone https://github.com/mpsido/python-trezor.git
    cd python-trezor

## 2. Install dependencies 

You need superuser priviledges for this step.

    sudo apt-get -y install python-dev cython libusb-1.0-0-dev libudev-dev (note: required for python-trezor testing)

## 3. Configure your usb module

We need to tell your kernel to use the [hidraw module](https://www.kernel.org/doc/Documentation/hid/hidraw.txt) to communicate with the hardware device. If you don't your kernel will treat the device as a mouse or a keyboard.

Create a file named 99-dev-kit.rules in your /etc/udev/rules.d/ folder and write this content in that file (super user priviledges are required for this step).

    # 0483:df11 STMicroelectronics STM Device in DFU Mode
    SUBSYSTEM=="usb", ATTR{idVendor}=="0483", ATTR{idProduct}=="df11", MODE="0666"
    # 0483:3748 STMicroelectronics ST-LINK/V2
    SUBSYSTEM=="usb", ATTR{idVendor}=="0483", ATTR{idProduct}=="3748", MODE="0666"
    # 0483:374b STMicroelectronics ST-LINK/V2.1 (Nucleo-F103RB)
    SUBSYSTEM=="usb", ATTR{idVendor}=="0483", ATTR{idProduct}=="374b", MODE="0666"
    # 534c:0001 SatoshiLabs Bitcoin Wallet [TREZOR]
    SUBSYSTEM=="usb", ATTR{idVendor}=="534c", ATTR{idProduct}=="0001", MODE="0666"
    KERNEL=="hidraw*", ATTRS{idVendor}=="534c", ATTRS{idProduct}=="0001", MODE="0666"
    # 1209:53c0 SatoshiLabs TREZOR v2 Bootloader
    SUBSYSTEM=="usb", ATTR{idVendor}=="1209", ATTR{idProduct}=="53c0", MODE="0666"
    KERNEL=="hidraw*", ATTRS{idVendor}=="1209", ATTRS{idProduct}=="53c0", MODE="0666"
    # 1209:53c1 SatoshiLabs TREZOR v2
    SUBSYSTEM=="usb", ATTR{idVendor}=="1209", ATTR{idProduct}=="53c1", MODE="0666"
    KERNEL=="hidraw*", ATTRS{idVendor}=="1209", ATTRS{idProduct}=="53c1", MODE="0666"

Restart your machine or force your udev kernel module to [reload the rules](https://unix.stackexchange.com/questions/39370/how-to-reload-udev-rules-without-reboot).

If when you plug the device your OS is not seeing the device, skip the warning on device's screen saying that the signature is wrong and then try [this](https://askubuntu.com/questions/645/how-do-you-reset-a-usb-device-from-the-command-line).

## 4. Generate a skycoin address from seed

    ./trezorctl skycoin_address seed

# How to compile and run tests 

## Trezor-crypto

This repository includes header files coming from [Trezor-crypto](https://github.com/trezor/trezor-crypto/) repository.

Download the repository

    git clone git@github.com:trezor/trezor-crypto.git

Then setup the TREZOR_CRYPTO_PATH environment variable:

    export TREZOR_CRYPTO_PATH=$PWD/trezor-crypto
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$TREZOR_CRYPTO_PATH


The dependancy libTrezorCrypto.a can be recompiled from sources.

Add this line to the CFLAGS: "CFLAGS += -DUSE_BN_PRINT=1"

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
