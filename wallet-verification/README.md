# Scripts description
The purpose of these scripts is to validate signatures of the wallets that are hosted on the [Downloads](https://www.skycoin.net/downloads/) page of the skycoin.net.

## Script steps:
1. Download wallets from the skycoin.net
2. Download signatures from the skycoin.net
3. Download and register public key from the https://raw.githubusercontent.com/skycoin/skycoin/develop/gz-c.asc
4. Validate signatures

# Supported platforms:
1. Linux
1. macOS

# Required packages
1. curl
2. gpg

## macOS gpg (2.2.7) installation instructions:
```sh
    $ brew install gnupg 
```

## Ubuntu gpg (1.4.20) installation instructions:
```sh
$ sudo apt-get install gnupg
```

# How to install/run on osx & linux
To run validation process perform following command:
```sh
$ export VERSION=0.23.0
$ sh ./run.sh
```