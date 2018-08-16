# Skycoin hardware wallet

This folder provides a firmware implementing skycoin features, and tools to test it.

The firmware itself is under [tiny-firmware](https://github.com/skycoin/services/tree/master/hardware-wallet/tiny-firmware) folder.
The firmware had been copied and modified from [this repository](https://github.com/trezor/trezor-mcu).

The [skycoin-api](https://github.com/skycoin/services/tree/master/hardware-wallet/skycoin-api) folder contains the definition of the functions implementing the skycoin features.

The [go-api-for-hardware-wallet](https://github.com/skycoin/services/tree/master/hardware-wallet/go-api-for-hardware-wallet) defines functions that act as code example to communicate with the firmware using a golang code.

Follow up [the wiki](https://github.com/skycoin/services/wiki/Hardware-wallet-project-advancement) to keep track of project advancement.
See also the wiki about integration with skycoin web app [here](https://github.com/skycoin/services/wiki/Hardware-wallet-integration-with-skycoin-web-wallet).
