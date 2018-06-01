#!/bin/bash
set -e

IMAGE=trezor-mcu-build-emulator64
TAG=${1:-master}
ELFFILE=build/trezor-emulator64-$TAG

docker build -f Dockerfile.emulator -t $IMAGE .
docker run -t -v $(pwd)/build:/build:z $IMAGE /bin/sh -c "\
	git clone https://github.com/skycoin/services.git && \
	cd services/hardware-wallet/tiny-firmware && \
	git checkout $TAG && \
	export EMULATOR=1
	export HEADLESS=1
	make -C vendor/nanopb/generator/proto/ && \
	make -C protob/ && \
	make -C emulator/ && \
	make

	cp skycoin.elf /$ELFFILE"
