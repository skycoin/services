It will have api for
- generate address for coin
- check balance of list of addresses, get unspent output set
- sign transaction
- track status of transaction

There is one interface. Then you intialize an implementation for each coin.

- need to have list of coins the interface is implemented for
- there will be choice of implementations.

For bitcoin we will have
- scanning wallet implementation
- btcd direct connect implementation
- blockchain.info json api implementation
- skycoin website/exchange api thin client json api implementation

- the command line program must be cabable of listing which coins it supports. And which implementations are available for each coin