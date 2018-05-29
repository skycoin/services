# Pending transactions monitor
**Please note:** Ideally, this service must be hosted on the same server where the Skycoin node, to prevent issues related with time differences. In case when node and service are hosted on different servers, these servers must have the same time set.

## Usage
To run monitor perform following command:
```sh
go run ./monitor.go
  -config="" # Path to config file (optional, must be specified only when a config file is used)
  -pendingTime=60 # Max pending transaction time (seconds, default = 60)
  -nodeAddress="http://127.0.0.1:6420" # Path to the Skycoin node API
  -mailHost="smtp.server" # SMTP server (smtp.gmail.com:587)
  -mailUsername="sender@gmail.com" # SMTP server user
  -mailPassword="password of sender account" #  SMTP server password
  -mailToAddress="received@email.com" # Receiver of notifications
```

Config example:
```
pendingTime 180
nodeAddress http://127.0.0.1:6420
mailHost smtp.server
mailUsername sender@gmail.com
mailPassword password of sender account
mailToAddress received@email.com
```

To run service using config file perform following command:
```sh
go run ./monitor.go
  -config="./config.config"
```