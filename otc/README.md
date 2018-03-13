# otc

OTC is a daemon with the sole purpose of exposing an HTTP API that allows users to exchange a variety of currencies (BTC, ETH, etc.) for Skycoin (and others in the future).

# running (for testing)

1. start [skycoin](github.com/skycoin/skycoin) node

```
$ cd github.com/skycoin/skycoin
$ make run
```

2. start [btcwallet](https://github.com/btcsuite/btcwallet) node (requires btcd but for testing you don't need that)

```
$ ./btcwallet -u otc -P OTC 
```
**NOTES:**
If you don't have BTC wallet you should create a new one using following command
```
$ ./btcwallet -u username -P passphrase --create
```

3. start otc

**NOTES:**
1. seed - seed of the new SkyCoin wallet
2. account - passphrase of BTC wallet certificate

```
$ cd github.com/skycoin-karl/services/otc
$ go build
$ ./otc
```

# frontend

OTC's frontend is exposed as an HTTP API. 

## /api/bind

This creates a new [request](#request) in the backend and returns JSON output on the frontend. The user can then send their currency to the returned address and the process will begin.

**http request**

```json
{
	"address": "...skycoin address...",
	"drop_currency": "BTC"
}
```

* `address` is the user's skycoin address where skycoin will be delivered
* `drop_currency` determines the type of `drop_address` to generate (what currency the user wants to deposit)

**http response**

```json
{
	"drop_address": "...",
	"drop_currency": "BTC",
	"drop_value": 159900
}
```

* `drop_address` is the address of type `drop_currency` for the user to send their currency to
* `drop_currency` is the same as sent in the request
* `drop_value` is the current price of 1 SKY in terms of `drop_currency`
	* represented as satoshis, example: 159900 = 0.00159900 BTC

## /api/status

This gets the [metadata](#request) of a request and returns it to the user.

**http request**

```
{
	"drop_address": "...",
	"drop_currency": "..."
}
```

* `drop_address` denotes the address they want the status of
* `drop_currency` denotes the currency of the address they want the status of

**http response**

```
{
	"status": "...",
	"updated_at": 1519131184
}
```

* `status` is one of the following:
	* `waiting_deposit` - skycoin address is bound, no deposit seen yet 
	* `waiting_send` - deposit detected, waiting to send to user 
	* `waiting_confirm` - skycoin sent, waiting to confirm transaction 
	* `done` - skycoin transaction confirmed 
	* `expired` - drop expired
* `updated_at` is the unix time (seconds) when the request was last updated

# admin api

## /api/status

Get status of OTC.

**http request**

Nothing, just a GET.

**http response**

```json
{
	"prices": {
		"internal": 150000,
		"internal_updated": 1519131184,
		"exchange": 119833,
		"exchange_updated": 1519131184
	},
	"source": "internal",
	"paused": true
}
```

## /api/price

Set the price of `internal` source. 

**http request**

```json
{
	"price": 119833
}
```

* `price` is the satoshi value of 1 SKY. `119833` is equal to `0.00119833 BTC`.

## /api/source

Set the price source.

**http request**

```json
{
	"source": "exchange"
}
```

* `source` is either `exchange` to get pricing from an exchange, or `internal` to use the manually set `internal` price

## /api/pause

Pause state transitions.

**http request**

```json
{
	"pause": true
}
```

* `pause` is a boolean denoting whether to pause or not

## transactions

### transaction

```json
{
    "address": "2dvVgeKNU7UHdvvBUVZXbBaxoTkpemo1cmg",
    "status": "waiting_confirm",
    "txid": "acfc6bd9e5b0b8eca1dad7e393658b51f01eda0f961a6b336a130ebc752565b8",
    "rate": {
        "value": 150000,
        "source": "internal"
    },
    "drop": {
        "address": "mnum2BxQ47qERaGtyfShAG1wkzJhhJE2J2",
        "currency": "BTC",
        "amount": 50000000
    },
    "timestamps": {
        "created_at": 1519131184,
        "deposited_at": 1519131200,
        "sent_at": 1519131500,
        "updated_at": 1519131500
    }
}
```

* `address` is the destination skycoin address
* `status` is one of the status codes for the entire transaction
* `txid` is the txid of the otc->user skycoin transaction (final step)
* `rate` contains information about the rate used when skycoin was sent to the user
    * `value` is the value of 1 SKY in the drop currency (`150000` is `0.00150000 BTC` in this example)
    * `source` denotes which rate source was used (either `exchange` or `internal` for now)
* `drop` contains information about the deposit address (the "drop")
    * `address` is the address being scanned for deposit
    * `currency` denotes the currency of the drop (currently just BTC)
    * `amount` denotes the amount of drop currency received (`50000000` is `0.50000000 BTC` in this example)
* `timestamps` contains unix times of each step
    * `created_at` is when the transaction was created
    * `deposited_at` is when the deposit was detected 
    * `sent_at` is when the skycoin was sent to the user
    * `updated_at` is the last time the request was updated and saved to disk

### /api/transactions

Returns all transactions.

```
[
    {transaction},
    {transaction}...
]
```

### /api/transactions/pending

Returns uncompleted transactions.

```
[
    {transaction},
    {transaction}...
]
```

### /api/transactions/completed

Returns created transactions.

```
[
    {transaction},
    {transaction}...
]
```