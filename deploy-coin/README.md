# Configure and run new coin network

Tools from this folder are used to deploy network for new coin.

## Installation

* Install create-coin tool from cmd/create
* Install start-coin tool from cmd/start

## Run

1. Create JSON configuration file using create-coin tool. You can override default nubmer of coin distribution addresses and trusted peers.
2. Use start-coin to launch several nodes with addresses from list of trusted peers.
3. Use start-coin to launch several slave nodes.
4. Use start-coin to launch master node, which will split initial coin volume among distribution addresses.
5. Check that block with initial transaction is broadcasted over network.


