# Hardware wallet project advancement

This document aims at describing the features related to skycoin hardware wallet and follow up their advancement.

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" -->

## Firmware's features

This is the firmware for the skycoin device which is intended to safely store a "seed" corresponding to a skycoin wallet.

It can generate addresses, sign messages or transactions and check signatures.

The hardware wallet has two buttons: Ok and Cancel. The user has to press one of these buttons when the hardware wallet requires user confirmation.

It communicates with a host PC using a USB wire. Please use [Skycoin's web application](https://github.com/skycoin/skycoin) or [Skycoin Command line interface](https://github.com/skycoin/skycoin/tree/develop/cmd/cli) tools to communicate with this firmware.

More informations in tiny-firmware's [README.md](https://github.com/skycoin/services/blob/master/hardware-wallet/tiny-firmware/README.md) section 6. How to read the firmware's code.

### Firmware's security

#### No message to get private key

The skycoin wallet is able to use a private key to sign any message and output the signature to the PC. But there is no way to have the firmware output the private key used to sign the message.

#### After backup: no message to get the seed

The seed used to generate key pairs is stored in the firmware's memory.

This seed is important it represents the wallet itself.

For safety it is strongly recommended that the user keeps a backup of that seed handwritten in a paper stored somewhere safe.

This seed can be very useful to recover the wallet's money in case the skycoin hardware wallet is lost.

#### Backup the seed

When the hardware wallet is freshly configured with a seed. The screen displays a **NEEDS BACKUP** message. This means that you can send a backup message to the hardware wallet to enter backup mode.

If a pin code was set (see [PIN code configuration section](#pin-code-protection)), it is required to enter backup mode.

The backup mode will display every word of the seed one by one and wait for the user to press the Ok button between each word. The user is supposed to copy these words on a paper (the order matters).

After a first round the words are displayed a second time one by one as well. The user is supposed to check that he did not mispelled any of these words.

Warning 1: once the backup is finished the **NEEDS BACKUP** disappears from the hardware wallet's screen and there is no way to do the backup again. If you feel you did not backup your seed properly better generate a new one and discard this one before you invested Skycoins on the wallet corresponding to that seed.

Warning 2: It is strongly recommended to do the backup short after the wallet creation, and before you invested Skycoins in it. If you loose a wallet that has an open door to do a backup, the person who finds it can use this backup to get the seed out of it. Especially if you did not [configure a PIN code](#pin-code-protection).

#### Don't erase your seed

At the time this document is written the hardware wallet is only able to store one seed. **TODO ? new feature store more than one seed ?**

If the user sends a seed setting message, the hardware wallet's screen asks the user if he wants to write the seed. If the user presses hardware wallet's Ok button. The new seed is stored and if there was an other seed before, it is **gone forever**.

So don't configure a new seed on a hardware wallet that is representing a wallet you are still using (see [backup section](#backup-the-seed) to avoid this problem).

#### PIN code protection

You can configure a PIN code on the hardware wallet. Check [this documentation](https://doc.satoshilabs.com/trezor-user/enteringyourpin.html) to see how to use the PIN code feature.

You can modify an existing PIN code. But the previous PIN code will be asked.

If you are not able to input a correct PIN code there is no way to change it apart from [wiping the device](#wipe-the-device).

#### PIN code cache

The PIN code is required for 
* address generation (can be cached)
* check signature (can be cached)
* signature
* device backup

If the user inputs a correct PIN code once, the PIN code is cached. When the PIN code is cached the operations where PIN code cache is authorised do not ask PIN code again to perform properly.

The PIN code stays cached until the device is plugged off.

The PIN code has to be input every time no matter what for the operation that do not allow PIN code cache.

#### PIN code brute force protection

If the user enters a wrong PIN code, the next time he asks for an operation requiring to enter a PIN code he has to wait extra seconds before he can enter it.

The amount of time he has to wait before he can try again increases everytime he enters a wrong PIN code.

#### Wipe the device

A message exist to wipe the device. It erases seed and PIN code.

When the device receives a wipe message it prompts the user to confirm by pressing Ok button.

There is no way back after a wipe. All the stored data is lost.

#### Passpharse protection

**TODO ?** [check this issue](https://github.com/skycoin/services/issues/134)
1) When to ask passphrase ?
2) Has impacts on web wallet integration

#### Memory encryption

**TODO ?**
**Use passphrase as a key for encryption ?**

### Firmware dependencies to external code

All the dependencies to external code are located in [tiny-firmware/vendor](https://github.com/skycoin/services/blob/master/hardware-wallet/tiny-firmware/vendor) directory.

It is worth mentioning 
* [libopencm3](https://github.com/libopencm3/libopencm3) which is a library managing low level interface with STM32 microchip.
* nanopb contains few source files pb_common.c, pb_encode.c, pb_decode.c. They are low level interface to decode google protobuf messages used in the communication with the PC.

### Factory test mode

**TODO** [check this issue](https://github.com/skycoin/services/issues/133)

### Other TODOs

Known bugs and possible improvements:
* The first time the device is plugged in it displays a "Storage failure detected message" that disapears the next time the device is connected. The text message has to be changed for "Storage not initialized, please restart".
* The Storage structure in storage.c still contain fields copied from trezor code base that we are not using anymore.
* the "MAGIC" four letters used by the bootloader to recognize skycoin firmware are still: TRZR, should be changed for SKCN or SKYN
* When the device is waiting for a Pin code it waits forever until the PIN code message arrives [see this issue](https://github.com/skycoin/services/issues/135)

## Integration with the skycoin web wallet

[Skycoin's web application](https://github.com/skycoin/skycoin)

Check this [pull request](https://github.com/skycoin/skycoin/pull/1686) to see most up to date code.

### Backup

**TODO** [check this issue](https://github.com/skycoin/skycoin/issues/1708)

### Pin code

**TODO** [check this issue for PIN code integration](https://github.com/skycoin/skycoin/issues/1765)

**TODO** [check this issue for PIN code configuration from the web wallet](https://github.com/skycoin/skycoin/issues/1768)

The PIN code can be configured from skycoin-cli using command:

    skycoin-cli deviceSetPinCode

### Use many connected devices at the same time

**TODO** [check this issue](https://github.com/skycoin/skycoin/issues/1709)

### Wipe a device from the web wallet

**TODO** [check this issue](https://github.com/skycoin/skycoin/issues/1769)


### Web wallet dependencies to external code

Here is the external code that was added in the project for hardware wallet integration:
* Low level interface to decode google protobuf messages used in the communication with the PC under [github.com/golang/protobuf](https://github.com/mpsido/skycoin/tree/develop-hardware-wallet/vendor/github.com/golang/protobuf).
* To send/receive messages from/to usb wire : [wire](https://github.com/mpsido/skycoin/tree/develop-hardware-wallet/src/device-wallet/wire), [usb](https://github.com/mpsido/skycoin/tree/develop-hardware-wallet/src/device-wallet/usb) and [usbhid](https://github.com/mpsido/skycoin/tree/develop-hardware-wallet/src/device-wallet/usbhid)

### Other TODOs

Known bugs and code improvements:
* Golang functions communicating with the device need to timeout if the device is not ansering. [issue #1771](https://github.com/skycoin/skycoin/issues/1771)
* Code would be more easy to work on if there was a "device-wallet" factory (very useful when we will fix issue [#1709](https://github.com/skycoin/skycoin/issues/1709))
* When the wallet is waiting for a PIN code we might have to create a callback to be called when the pin code arrives (current code makes the GUI to say what the PIN code is sent for, and the parameters need to be repeated).
