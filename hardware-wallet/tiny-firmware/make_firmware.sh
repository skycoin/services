make -C vendor/nanopb/generator/proto/
make
make -C bootloader/ align

cp skycoin.bin bootloader
pushd bootloader
./firmware_align.py bootloader.bin
./firmware_sign.py -f skycoin.bin

cp bootloader/bootloader.bin combine/bl.bin
cp skycoin.bin combine/fw.bin
pushd combine/ && ./prepare.py
popd; popd;

#st-flash erase

alias st-skycoin='pushd bootloader/combine/; st-flash write combined.bin 0x08000000; popd;'
