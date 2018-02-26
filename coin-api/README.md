
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
                    "address": "99902999f9s99ds999s9",
                    "lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
                    "tm": "1503458909",
                    "type": "deterministic",
                    "version": "0.1"
                },
                {
                    "сid": "ETH",
                    "name":"ethereum",
                    "address": "99902999f9s99ds999s9",
                    "lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
                    "tm": "1503458909",
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
ResponseHeader: 202
ResponseBody {
    "status":"ok",
    "code": 0,
    "result": {
        "public":"9182b02c0004217ba9a55593f8cf0abec",
        "private":"99182b02c0004217ba9a55593f8cf0abec182b02c0004217ba9a55593f8cf0abec"
    },
}
```
##### Unsuccesful response:
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

##### POST /api/v1/btc/address/:key
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

#### BTC check the balance (and get unspent outputs) for an address
##### GET /api/v1/btc/address/:address
##### Successful response:
```
ResponseHeader: 200
ResponseBody {
    "status": "ok",
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
        "status":"pending",
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
    "signid":"392900939dijdked"
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
