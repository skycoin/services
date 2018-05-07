make -C vendor/nanopb/generator/proto/
make
make -C bootloader/ align

cp bootloader/bootloader.bin bootloader/combine/bl.bin
cp skycoin.bin bootloader/combine/fw.bin
pushd bootloader/combine/ && ./prepare.py
popd;

#st-flash erase

alias st-skycoin='pushd bootloader/combine/; st-flash write combined.bin 0x08000000; popd;'