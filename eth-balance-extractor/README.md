# Ethereum tokens holders balance extractor

## events-scanner
Required software:
Node >10.0.0

### Configuration
There is the 'Parameters' section in the index.js that contains configuration of the scanner.
For the first run install npm packages using ```npm i``` command from the ```events-scanner``` folder.

### Scanning process
To start scanning process use following command from the ```events-scanner``` folder:
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
start_block # first block of the pub keys extraction process
```

## Structure of Ethereum transactions
1. blockHash - hash of the block that transaction belongs to.
2. blockNumber - number of the block that transaction belongs to.
3. from - address of the sender wallet
4. gas - amount of gas that is used by transaction
5. gasPrice - amount of Ether for every unit of gas
6. hash - hash of the transaction
7. input - payload of transaction (may contain parameters of smart contract function in case if transaction preforms a call to smart contract) 
8. nonce - value found by miner 
9. to - receiver of the transaction
10. transactionIndex - transaction position in the block
11. value - Ether amount that was sent in the transaction
12. v, r, s - values of transaction signature (must be used for public key recovery)
