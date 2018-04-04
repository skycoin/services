# Coin-Api server

[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/services)](https://goreportcard.com/report/github.com/skycoin/services)

Server that responsible for managing keys, addresses, deposits of multiple
cryptocurrencies.


## Install

Download service code with and install dependencies

```
git clone https://github.com/skycoin/services $GOPATH/src/github.com/skycoin/services
cd $GOPATH/src/github.com/skycoin/services/coin-api
dep ensure
```

## Configure

```
[bitcoin]
// BTC node credentials, btc node is used first for controlling account balance
NodeAddress="localhost:8332"
User="aaa"
Password="123"
// Option for disabling tls
TLS=true
CertFile="btcd.cert"
// Block explorer - any public block explorer can be used
BlockExplorer="https://api.blockchain.info"

// Sky coin node credentials
[skycoin]
Host="0.0.0.0"
Port="9341"

[server]
ListenStr="0.0.0.0:8080"
// Http server config are neccessary for all public servers
ReadTimeout = 5
WriteTimeout = 10
IdleTimeout = 120
```

## Run

Start from main

```
cd cmd
go run main.go -config config.toml
```


Start from command line

```
cd cmd/cli
go build
./cli server start <config_file_path>
```


## Command line API

### Btc command

generate keypair

```
./cli generatekeys
```

generate address

```
./cli generateaddr <publicKey>
```

check balance

```
./cli checkbalance <btcAddress>
```

## API overview

### Common server API
```
GET /api/v1/ping
ResponseHeader: 200
ResponseBody {
    "code": 0,
    "status":"ok",
    "result": {
        "message":"pong"
    }
}
```

#### GET /api/v1/list

List of supported currencies

#### Successful response:
```
ResponseHeader: 200
ResponseBody {
    "status": "ok",
    "code": 0,
    "result": {
            [
                {
                    "сid": "LTC",
                    "name":"litecoin",
                    "timestamp": "1503458909",
                    "type": "deterministic",
                    "version": "0.1"
                },
                {
                    "сid": "ETH",
                    "name":"ethereum",
                    "timestamp": "1503458909",
                    "type": "deterministic",
                    "version": "0.1"
                }
            ]
    }
}
```

#### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody: {
    "status": "error",
    "code": -4,
    "result": {
        "description": "Wallet error, no coin discovered or coin initialization error",
    }
}
```


### BTC API
#### Generate key pair
##### POST /api/v1/btc/keys
##### Successful response:

```
ResponseHeader: 201
ResponseBody {
    "status":"ok",
    "code": 0,
    "result": {
        "public":"9182b02c0004217ba9a55593f8cf0abec",
        "private":"99182b02c0004217ba9a55593f8cf0abec182b02c0004217ba9a55593f8cf0abec"
    },
}
```

##### Unsuccessful response:
```
ResponseHeader: 505
ResponseBody {
    "status": "error",
    "code": -32603,
    "result": {
        "description": "Unable to generate keypair, internal server error"
    }
}
```

#### BTC generate address based on public key
```
POST /api/v1/btc/address
```

Request

```
{
    "key":"02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc"
}
```

##### Successful response:
```
ResponseHeader: 201
ResponseBody {
    "status": "ok",
    "code": 0,
    "result": {
        "address": "9182b02c0004217ba9a55593f8cf0abecc30d041e094",
    }
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status": "error",
    "code": -5,
    "result": {
        "description": "Unable to generate address, given key not found"
    }
}
```

#### BTC check the balance (and get unspent outputs) for an address in satoshis
##### GET /api/v1/btc/address/:address
##### Successful response:
```
ResponseHeader: 200
ResponseBody {
             	"status": "Ok",
             	"code": 200,
             	"result": {
             		"address": "1M3GipkG2YyHPDMPewqTpup83jitXvBg9N",
             		"balance": 26943184,
             		"deposits": [
             			{
             				"amount": 42482,
             				"confirmations": 278,
             				"height": 514968
             			},
             			{
             				"amount": 116000,
             				"confirmations": 415
             			}
             	    ]
             	}
             }
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status": "error",
    "code": -32602,
    "result": {
        "description": "Unable to find given address"
    }
}
```

#### BTC check the status of a transaction (tracks transactions by transaction hash)
##### GET /api/v1/btc/transaction/:transid
##### Successful response:
```
ResponseHeader: 200
ResponseBody {
    "status": "ok",
    "code": 0,
    "result": {
        "transid": "7ba9a55593f8cf0abecc30d041e094",
        "status": {
           "ver":1,
           "inputs":[
              {
                 "sequence":4294967295,
                 "witness":"",
                 "prev_out":{
                    "spent":true,
                    "tx_index":335439599,
                    "type":0,
                    "addr":"1PiMMmKxbpDhoPBweEDKiu3FCEsjLwsyqc",
                    "value":293717322,
                    "n":1,
                    "script":"76a914f924d0d0959f4a45e4a0b6ff390ebd38cdff726d88ac"
                 },
                 "script":"47304402206de55e3b9f013f337cf31399d9314a5181b9eb5f02c1a9a3618f8b812901871902202ecce58957c179b33a53229bbf4b7eb9fad65dc18cdc3a07e0881008a6c858ee01210213209d0af0becd42171eb83befdc3be1e408ec4e2953e9e1d1442eae9958fb02"
              }
           ],
           "weight":900,
           "relayed_by":"0.0.0.0",
           "out":[
              {
                 "spent":false,
                 "tx_index":335447338,
                 "type":0,
                 "addr":"1PK3WGyQ7t1JPVjrCJihCH9rfGjq94q5on",
                 "value":293156277,
                 "n":0,
                 "script":"76a914f4bc586d07e7936eb743c697fa4885b9294472cf88ac"
              },
              {
                 "spent":false,
                 "tx_index":335447338,
                 "type":0,
                 "addr":"1LVNiEn7rYdGfbEhR5GgkeHqB7zu3jubQg",
                 "value":556185,
                 "n":1,
                 "script":"76a914d5c82905d98094e3bc348fa874ff9810ffc0b9f288ac"
              }
           ],
           "lock_time":0,
           "size":225,
           "double_spend":false,
           "time":1520610698,
           "tx_index":335447338,
           "vin_sz":1,
           "hash":"ac2e8c4dd81253c819824c8725f7ad359ab76a43824b3b3e9338fb1baf90b819",
           "vout_sz":2
        },
    }
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status": "error",
    "code": -32602,
    "result": {
        "description": "Unable to find given transaction"
    }
}
```

### Multicoin API

#### Generate key pair
##### POST /api/v1/:coin/keys
##### Successful response:
```
ResponseHeader: 201
ResponseBody {
    "status":"ok",
    "code": 0,
    "result": {
        "public":"9182b02c0004217ba9a55593f8cf0abec",
        "private":"99182b02c0004217ba9a55593f8cf0abec182b02c0004217ba9a55593f8cf0abec"
    }
}
```
##### Unsuccessful response:
```
ResponseHeader: 505
ResponseBody {
    "status": "error",
    "code": -32603,
    "result": {
        "description": "Unable to generate keypair, internal server error"
    }
}
```

#### Generate address from public key
##### POST /api/v1/:coin/address/:key

Request body

```
{
    "key":"032417bc1f336ad55a0686a956ccc687ac6be4c0413758c1f78bf82e29c8dcf8b9"
}
```

##### Successful response:
```
ResponseHeader: 201
ResponseBody {
    "status":"ok",
    "code": 0,
    "result": {
        "address": "9182b02c0004217ba9a55593f8cf0abecc30d041e094",
    }
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status":"error",
    "code": -5,
    "result": {
        "description": "Unable to generate address, given key not found"
    }
}
```

#### check the balance (and get unspent outputs) for an address
##### GET /api/v1/:coin/address/:address
##### Successful response:
```
ResponseHeader: 200
ResponseBody {
    "status":"ok",
    "code": 0,
    "result": {
        "address": "9182b02c0004217ba9a55593f8cf0abecc30d041e094",
        "balance": 12.07,
    }
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status":"error",
    "code":-25 //-25 -26 -27 possible codes
    "result": {
        "description":"description according to given error code"
    }
}
```

#### sign a transaction
##### POST /api/v1/:coin/transaction/:transid/sign
##### Successful response:
```
Request {
    "signid":"392900939dijdked",
    "sourceTrans":"392900939dijdked392900939dijdked",
}

ResponseHeader: 201
ResponseBody {
    "status":"ok",
    "code": 0,
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status":"error",
    "code":-25 //-25 -26 -27 possible codes
    "result": {
        "description":"description according to given error code"
    }
}
```

#### inject transaction into network
##### PUT /api/v1/:coin/transaction/:netid
##### Successful response:
```
Request {
    "transid":"392900939dijdked"
}

ResponseHeader: 201
ResponseBody {
    "status":"ok",
    "code": 0,
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "status":"error",
    "code":-25 //-25 -26 -27 possible codes
    "result": {
        "description":"description according to given error code"
    }
}
```

#### check the status of a transaction (tracks transactions by transaction hash)
##### GET /api/v1/sky/transaction/:transid
##### Successful response:
```
ResponseHeader: 200
ResponseBody {
    "status":"ok",
    "code": 0,
    "result": {
        "transid":"392900939dijdked"
        "status":"pending"
    }
}
```
##### Unsuccessful response:
```
ResponseHeader: 404
ResponseBody {
    "code":-25,
    "status":"error",
    "result": {
        "description":"any possible transaction error description"
    }
}
```
