make -C ../vendor/trezor-crypto clean
make -C ../vendor/trezor-crypto 
make -C ../../skycoin-api clean
make -C ../../skycoin-api
cp ../../skycoin-api/libskycoin-crypto.so .
make -C ../../skycoin-api clean
make -C ../vendor/trezor-crypto clean
ln -sf ../../skycoin-api/skycoin_crypto.py .