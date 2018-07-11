#!/bin/sh

# First launch a file server to "GET peers.txt" from it
/app/gohttpserver --root=/app/public &
sleep 10
dns2peerlist -service skystack_skycoin-node -format text -output /app/public/peers.txt

# From here its a copy of https://github.com/skycoin/skycoin/blob/develop/docker_launcher.sh
COMMAND="skycoin --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 skycoin

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R skycoin:skycoin $( realpath $DATA_DIR )
chown -R skycoin:skycoin $( realpath $WALLET_DIR )

su skycoin -c "$COMMAND"