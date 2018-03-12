# Configure and run new coin network

Tools from this folder are used to deploy network for new coin.

## Preparing to installation

Supposedly you are in deploy-coin directory where this README.md file is located:

just execute ``make prepare-installation``, don't pay attention if you see any error messages like `File exists`
it doesn't matter at this step.

## Build

Execute ``make build-create-coin:``
and     ``make build-start-coin``

#Install 
Create coin configuration ``./cmd/create -addr 1 -file config.txt -addr`` 
for specifying count of addresses for initial coin distribution ``-file`` specify filename to store generated config

Start master node ``./cmd/start -config config.txt -master true``

<!-- * Install create-coin tool from cmd/create
* Install start-coin tool from cmd/start -->

<!-- ## Build

1. Create JSON configuration file using create-coin tool. You can override default nubmer of coin distribution addresses and trusted peers.
2. Use start-coin to launch several nodes with addresses from list of trusted peers.
3. Use start-coin to launch several slave nodes.
4. Use start-coin to launch master node, which will split initial coin volume among distribution addresses.
5. Check that block with initial transaction is broadcasted over network. -->

## Note
	Known issue - in case of interrupting node with SIGINT node cannot be recovered, delete data.db from .skycoin dir to clear everything and start it again.

In case of multi-node install - modify trusted peers and connect to peer node. 

  "port": 20101 - responsible for heartbeats between nodes
  "rpcPort": 20200 - used for rpc calls for getting data, requires -rpc option to be true
  "guiPort": 20300 - used for web interface, requires -gui option set to be true

1. Use start-coin to launch several nodes with addresses from list of trusted peers.
2. Use start-coin to launch several slave nodes.
3. Use start-coin to launch master node, which will split initial coin volume among distribution addresses.
4. Check that block with initial transaction is broadcasted over network. 

## !REMEMBER
All your nodes config is in $(GOPATH)/bin/start-coin/config.txt file or whatever you pointed out creating coin 
configuration few steps above.
