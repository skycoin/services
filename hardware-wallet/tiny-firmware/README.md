# Skycoin electronic wallet firmware

This firmware had been copied and modified from [trezor-mcu](https://github.com/trezor/trezor-mcu). Please refer to the [README.md file](https://github.com/trezor/trezor-mcu/blob/master/README.md) on that repository for extra details about bootloader and firmware compilation.

This code aims at tranforming the cipher library from [this repository](https://github.com/skycoin/skycoin/tree/develop/src/cipher) into firmware for the STM32 hardware.

## 1. Prepare environment

### Download GNU ARM Embedded Toolchain

    wget https://developer.arm.com/-/media/Files/downloads/gnu-rm/6-2017q2/gcc-arm-none-eabi-6-2017-q2-update-linux.tar.bz2

### Extract crosscompiler and add it to your path

    tar xjf gcc-arm-none-eabi-6-2017-q2-update-linux.tar.bz2
    export PATH="$PWD/gcc-arm-none-eabi-6-2017-q2-update/bin:$PATH"

### Install ST-LINK

Follow the steps [here](https://github.com/texane/stlink/blob/master/doc/compiling.md).

### Install google protobuf

```
sudo apt-get install protobuf-compiler python-protobuf golang-goprotobuf-dev
```

### Configure your usb module

We need to tell your kernel to use the [hidraw module](https://www.kernel.org/doc/Documentation/hid/hidraw.txt) to communicate with the hardware device. If you don't your kernel will treat the device as a mouse or a keyboard.

Create a file named 99-dev-kit.rules in your /etc/udev/rules.d/ folder and write this content in that file (*super user priviledges are required for this step*).

    ## 0483:df11 STMicroelectronics STM Device in DFU Mode
    SUBSYSTEM=="usb", ATTR{idVendor}=="0483", ATTR{idProduct}=="df11", MODE="0666"
    ## 0483:3748 STMicroelectronics ST-LINK/V2
    SUBSYSTEM=="usb", ATTR{idVendor}=="0483", ATTR{idProduct}=="3748", MODE="0666"
    ## 0483:374b STMicroelectronics ST-LINK/V2.1 (Nucleo-F103RB)
    SUBSYSTEM=="usb", ATTR{idVendor}=="0483", ATTR{idProduct}=="374b", MODE="0666"
    ## 534c:0001 SatoshiLabs Bitcoin Wallet [TREZOR]
    SUBSYSTEM=="usb", ATTR{idVendor}=="534c", ATTR{idProduct}=="0001", MODE="0666"
    KERNEL=="hidraw*", ATTRS{idVendor}=="534c", ATTRS{idProduct}=="0001", MODE="0666"
    ## 1209:53c0 SatoshiLabs TREZOR v2 Bootloader
    SUBSYSTEM=="usb", ATTR{idVendor}=="1209", ATTR{idProduct}=="53c0", MODE="0666"
    KERNEL=="hidraw*", ATTRS{idVendor}=="1209", ATTRS{idProduct}=="53c0", MODE="0666"
    ## 1209:53c1 SatoshiLabs TREZOR v2
    SUBSYSTEM=="usb", ATTR{idVendor}=="1209", ATTR{idProduct}=="53c1", MODE="0666"
    KERNEL=="hidraw*", ATTRS{idVendor}=="1209", ATTRS{idProduct}=="53c1", MODE="0666"

Restart your machine or force your udev kernel module to [reload the rules](https://unix.stackexchange.com/questions/39370/how-to-reload-udev-rules-without-reboot).

    sudo udevadm control --reload-rules
    sudo udevadm trigger

## 2. How to compile firmware

Then for the actual compilation source make_firmware.sh script:

    cd tiny-firmware
    . make_firmware.sh


If you want to compile on PC target set EMULATOR environment variable to 1:

    cd tiny-firmware
    export EMULATOR=1
    . make_firmware.sh
    ./skycoin.elf

[Optional] If you want to compile with a inverted screen and inverted buttons set REVERSE_SCREEN and REVERSE_BUTTONS environment variable to 1 :

    export REVERSE_SCREEN=1
    export REVERSE_BUTTONS=1
    
If you get SDL errors you might want to install these:

    sudo apt-get install libsdl2-dev libsdl2-image-dev

Works also with docker if you run the script:

    ./build-emulator.sh

## 3. How to burn the firmware in the device

### Use ST-LINK to burn the device

You can check the device is seen by your ST-LINK using this command:

    st-info --probe

To flash the device on a microcontroller of STM32f2xx family the command is:

    st-flash write combined.bin 0x08000000;

If you sourced the make_firmware.sh file as recommended. You can use the alias st-trezor to burn the device.

    st-trezor

## 4. Firmware signature

### Activate the signature checking

The bootloader code and Makefile had been updated to disable the signature checking by default.

If you want to activate it set SIGNATURE_PROTECT environment variable to 1 before making the bootloader.

    cd tiny-firmware/bootloader/
    export SIGNATURE_PROTECT=1
    make clean
    make

### Skip the wrong firmware signature warning

We are compiling a modified version of trezor original firmware. The bootloader we are using is configured to detect it and warn the user about this.
We are still in developpement steps, this warning is acceptable.
The devices allows the user to skip the warning but during that time the OS of the host computer may miss the opportunity to load the driver to communicate with the device.

If when you plug the device your OS is not seeing the device, skip the warning on device's screen saying that the signature is wrong and then try [this](https://askubuntu.com/questions/645/how-do-you-reset-a-usb-device-from-the-command-line).

If you are fast enough you can also quickly click the button "accept" on the device when the warning about wrong firmware signature appears.

### How to perform a custom signature

#### Environment setup

The first time we need to prepare the script for firmware signature:

    cd tiny-firmware/bootloader
    ./prepare_signature.sh

This script will generate and copy libskycoin-crypto.so and libtrezor-crypto.so in bootloader repertory.

These libraries are required to perform signature and signature checking from [firmware_sign.py](https://github.com/skycoin/services/blob/master/hardware-wallet/tiny-firmware/bootloader/firmware_sign.py) script.

Note: you also need to be able to generate skycoin key pairs. You can for instance use the [skycoin-cli](https://github.com/skycoin/skycoin).

#### You need three signatures:

The system stores five public keys and expects three signatures issued from one of these public keys.

The public keys are hardwritten in the bootloader's source code in file [signatures.c](https://github.com/skycoin/services/blob/master/hardware-wallet/tiny-firmware/bootloader/signatures.c)

The signatures are also present in [firmware_sign.py](https://github.com/skycoin/services/blob/master/hardware-wallet/tiny-firmware/bootloader/firmware_sign.py) script, in the "pubkeys" array.

#### Get the firmware

Copy back the file fw.bin located in tiny-firwmare/bootloader/combine/ in the bootloader folder.

This file was generated by tiny-firmware/make_firmware.sh. It contains "empty signature slots".

Name it skycoin.bin.

    cp combine/fw.bin skycoin.bin

#### Use your secret key to perform signature

Run:

    ./firmware_sign.py -s -f skycoin.bin

The command line tool will ask you in which of the three slots do you want to store the signature.

The it will ask you to provide a secret key that must correspond to one of the five public keys stored in the bootloader and the script as described above.

#### Recombine the firmware and the bootloader

From the tiny-firmware/bootloader repertory.

    cp skycoin.bin combine/fw.bin
    cd combine
    ./prepare.py

Then you can re-flash the firmware for instance with st-skycoin alias.

## 5. Communicate with the device

### Use golang code examples

Check [here](https://github.com/skycoin/services/tree/master/hardware-wallet/go-api-for-hardware-wallet) for golang code example communicating with the device.

Feel free to hack [main.go](https://github.com/skycoin/services/blob/master/hardware-wallet/go-api-for-hardware-wallet/main.go) file.

You can also try the trezorctl [python based command line](https://github.com/trezor/python-trezor).

## 6. How to read the firmware's code

The communication between PC and firmware is a master/slave model where the firmware is slave.
It reacts to messages but cannot initiate a communication.
The messages are defined using google protobuf code generation tools. The same file messages.proto can be copy pasted elswhere to generate the same structures in other coding languages.

The folder [go-api-for-hardware-wallet](https://github.com/skycoin/services/tree/master/hardware-wallet/go-api-for-hardware-wallet) provides examples to communicate with the device using golang.

The firmware has two components: the [bootloader](https://github.com/skycoin/services/tree/master/hardware-wallet/tiny-firmware/bootloader) and the [firmware](https://github.com/skycoin/services/tree/master/hardware-wallet/tiny-firmware/firmware).
The bootloader main role is to check firmware's signature in order to warn user in case the detected firmware is not the official firmware distributed by skycoin.

Here is a quick presentation of most important files of the firmware:

* protob/messages.proto: defines all the messages received the firmware can receive and their structure
* firmware/fsm.c: all the messages received in the firmware correspond to a call to a function in fsm.c, the corresponding between messages and function is defined in protob/messages_map.h (generated file)
* firmware/storage.c: manages the persistent (persistent when pluging out the power supply) memory of the firmware.
* oled.c/layout.c: manage screen display
* firmware/trezor.c: main entry point

On bootloader side it is worth mentioning:

* signatures.c: checks the firmware's signature matches skycoin's public keys
