make -C vendor/libopencm3/
make -C vendor/nanopb/generator/proto/
make -C firmware/protob/
# make -C vendor/skycoin-crypto/
export MEMORY_PROTECT=0
make
make -C bootloader/ align
make -C firmware/ sign

#sudo cp 99-dev-kit.rules /etc/udev/rules.d/

cp bootloader/bootloader.bin bootloader/combine/bl.bin
cp firmware/trezor.bin bootloader/combine/fw.bin
pushd bootloader/combine/ && ./prepare.py
popd;

#st-flash erase

alias st-trezor='pushd bootloader/combine/; st-flash write combined.bin 0x08000000; popd;'
