# Ethereum tokens holders balance extractor

## To start extraction process from scratch run following command:
```sh
eth-scan extractWallets 
http://url_of_eth_node_api  "0xaddress_of_token_smart_contract" "0xcode_of_transfer_method" ./path_to_snapshots_folder block_number_where_smart_contract_was_deployed threads_count
```

1. address_of_token_smart_contract - address of smart contract wallet
2. code_of_transfer_method - to find this code find any token transfer transaction (see [this](https://etherscan.io/tx/0x7767e8e4710bde871ecc2081fabb412f242d75b272e29c94c750ee444016f934) example: Input data - MethodID, "0xa9059cbb"). 
3. block_number_where_smart_contract_was_deployed - can be found in the transaction where smart contract was deployed ([example](https://etherscan.io/address/0xf230b790e05390fc8295f4d3f60332c93bed42e2) Contract Creator at tnx -> Block Height: 4212165)
4. threads_count - number of threads that extract blocks and transaction from the Ethereum node API.

## To continue extraction process run following command:
```sh
./eth-scan continueExtraction ./path_to_snapshot.csv http://url_of_eth_node_api  "0xaddress_of_token_smart_contract" "0xcode_of_transfer_method" ./path_to_snapshots_folder block_number_where_smart_contract_was_deployed threads_count
```

path_to_snapshot - path to any snapshot created before 