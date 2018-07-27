make -C vendor/nanopb/generator/proto/
make -C protob/

if [ -e "bootloader/skycoin_crypto.py" ]; then
    pushd bootloader
    ./prepare_signature.sh
    popd
fi

if [ -z "$EMULATOR" ]; then
    EMULATOR=0
fi

if [ "$EMULATOR" == "0" ]; then
    make -C vendor/libopencm3/
    make -C bootloader/ align
    make sign

    cp bootloader/bootloader.bin bootloader/combine/bl.bin
    cp skycoin.bin bootloader/combine/fw.bin
    pushd bootloader/combine/ && ./prepare.py
    popd;

    #st-flash erase
    alias st-skycoin='pushd bootloader/combine/; st-flash write combined.bin 0x08000000; popd;'
else
    make -C emulator/
    make
fi
