# Ethereum tokens holders balance extractor

## events-scanner
Required software:
Node >10.0.0

### Configuration
There is the 'Parameters' section in the index.js that contains configuration of the scanner.

### Scanning process
To start scanning process use following command:
```sh
npm run run
```

Result of scanning process will be 3 files:
1. negative_balance_wallets.csv
2. positive_balance_wallets.csv
3. zero_wallets.csv

These files have following structure:
1. Wallet address
2. Hash of any transaction, performed by wallet (that can be used for restoring wallet public key)
3. Balance
4. Transactions count

## eth-public-keys-extractor:
To start extraction process run following command:
```sh
eth-public-keys-extractor extractWalletsKeys http://url_of_eth_node_api ./path_to_wallets.csv # path to the file that contains wallets obtained by events-scanner utility
./path_to_dest_folder # path to folder where result file wallets.csv will be saved
```