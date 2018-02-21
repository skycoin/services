
## API overview

### Common server API
```
GET /ping
GET /list
```

### BTC API
#### Generate key pair
```
POST /api/v1/btc/keys
```

#### BTC generate address based on public key
```
POST /api/v1/btc/address
```

json body

```
{"key":"02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc"}
```

#### BTC check the balance (and get unspent outputs) for an address
```
GET /api/v1/btc/address/:address
```

#### BTC check the status of a transaction (tracks transactions by transaction hash)
```
GET /api/v1/btc/transaction/:transid
```

### Multicoin API
#### Generate address, private keys, pubkeys from deterministic seed
```
POST /api/v1/multi/:coin/address
```

#### check the balance (and get unspent outputs) for an address
```
GET /api/v1/multi/:coin/address/:address
```

#### sign a transaction
```
POST /api/v1/multi/:coin/transaction/sign/:sign
```

#### inject transaction into network
```
POST /api/v1/multi/:coin/transaction/:netid/:transid
```

#### check the status of a transaction (tracks transactions by transaction hash)
```
GET /api/v1/multi/:coin/transaction/:transid
```
