make
cp skycoin.bin ../firmware/bootloader
pushd ../firmware/bootloader
./firmware_align.py bootloader.bin
./firmware_sign.py -f skycoin.bin
cp skycoin.bin combine/fw.bin
pushd combine/
./prepare.py
popd; popd