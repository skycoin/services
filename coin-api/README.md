
## API overview

### Common server API
```
GET /api/v1/ping
ResponseHeader: 200
ResponseBody {
    "status":"pong"
}
```
```
GET /api/v1/list
ResponseHeader: 200
ResponseBody {
    [
        {
            "id": "LTC",
            "name":"litecoin",
            "address": "99902999f9s99ds999s9",
            "lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
            "tm": "1503458909",
            "type": "deterministic",
            "version": "0.1"
        },
        {
            "id": "ETH",
            "name":"ethereum",
            "address": "99902999f9s99ds999s9",
            "lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
            "tm": "1503458909",
            "type": "deterministic",
            "version": "0.1"
        }
    ]
    "status":"ok"
}
```
### BTC API
#### Generate key pair
```
POST /api/v1/btc/keys

ResponseHeader: 202
ResponseBody {
    "status":"ok",
    "public":"9182b02c0004217ba9a55593f8cf0abec",
    "private":"99182b02c0004217ba9a55593f8cf0abec182b02c0004217ba9a55593f8cf0abec"

}
```

#### BTC generate address based on public key
```
POST /api/v1/btc/address/:key

ResponseHeader: 202
ResponseBody {
    "address": "9182b02c0004217ba9a55593f8cf0abecc30d041e094",
}
```

#### BTC check the balance (and get unspent outputs) for an address
```
GET /api/v1/btc/address/:address

ResponseHeader: 200
ResponseBody {
    "address": "9182b02c0004217ba9a55593f8cf0abecc30d041e094",
    "balance": 12.07,
}
```

#### BTC check the status of a transaction (tracks transactions by transaction hash)
```
GET /api/v1/btc/transaction/:transid

ResponseHeader: 200
ResponseBody {
    "transid": "7ba9a55593f8cf0abecc30d041e094",
    "status":"pending",
}
```

### Multicoin API
#### Generate address, private keys, pubkeys from deterministic seed
```
POST /api/v1/:coin/address
```

#### check the balance (and get unspent outputs) for an address
```
GET /api/v1/:coin/address/:address

ResponseHeader: 200
ResponseBody {
    "address": "9182b02c0004217ba9a55593f8cf0abecc30d041e094",
    "balance": 12.07,
}
```

#### sign a transaction
```
POST /api/v1/:coin/transaction/:transid/sign

Request {
    "signid":"392900939dijdked"
}

ResponseHeader: 202
ResponseBody {
    "status":"Ok"
}
```

#### inject transaction into network
```
PUT /api/v1/:coin/transaction/:netid

Request {
    "transid":"392900939dijdked"
}

ResponseHeader: 202
ResponseBody {
    "status":"ok"
}
```

#### check the status of a transaction (tracks transactions by transaction hash)
```
GET /api/v1/sky/transaction/:transid

ResponseHeader: 200

ResponseBody {
    "transid":"392900939dijdked"
    "status":"pending"
}
```
