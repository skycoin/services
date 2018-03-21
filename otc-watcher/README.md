# otc-watcher

Watches addresses on blockchain and stores output information.

*NOTE*: Currently only BTC address watching is supported. Adding support for other currencies in the future will be easy, thanks to the [currency connection interface](pkg/currency/currency.go).

## running

Right now the only dependency is to have `btcwallet` running. Then you can start `otc-watcher` with the following command (assuming you ran `go build`):

```
./otc-watcher -rpc_node="localhost:8332" \
              -rpc_user="username"       \
              -rpc_pass="password"       \
              -wallet_account="account"  \
              -wallet_pass="password"
```

* `rpc_node` is the `btcwallet` rpc server listening address
* `rpc_user` is the `btcwallet` rpc username
* `rpc_pass` is the `btcwallet` rpc password
* `wallet_account` is the `btcwallet` account name (will be created if missing)
* `wallet_pass` is the `btcwallet` passphrase used when creating the wallet

## http api

### /register

Registers an address to be watched. OTC needs to watch addresses indefinitely so there's no need for an "unregister" call right now.

#### request

```js
{
	// address to be watched
	"address": "1Hz96kJKF2HLPGY15JWLB5m9qGNxvt8tHJ",
	// currency type of address (must be supported by otc-watcher)
	"currency": "BTC"
}
```

#### response

`200 OK` if good, anything else and there was an error.

### /outputs

Gets the outputs for an address that was previously registered.

#### request

```js
{
	// address to get the outputs of
	"address": "1Hz96kJKF2HLPGY15JWLB5m9qGNxvt8tHJ",
	// currency type of address (must be supported by otc-watcher)
	"currency": "BTC"
}
```

#### response

```js
{
	// transaction hash
	"e0ba30a518d5c52504d84446a645d8865513e4fd7a4db53b507705eb43812ed0": {
		// output index	
		"1": {
			// satoshi amount
			"amount": 684830048,
			// number of confirmations (always > 1)
			"confirmations": 2,
			// block height of this output
			"height": 514553
		}
	}
}
```

The transaction hash and output index can then be used to create unique "deposit" ids for use throughout OTC.