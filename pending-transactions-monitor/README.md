# Pending transactions monitor
## Usage
To run monitor perform following command:
```sh
go run ./monitor.go
  -config="" # Path to config file (optional, must be specified only when a config file is used)
  -nodeAddress="http://127.0.0.1:6420" # Path to the Skycoin node API
  -mailHost="smtp.server" # SMTP server (smtp.gmail.com:587)
  -mailUsername="sender@gmail.com" # SMTP server user
  -mailPassword="password of sender account" #  SMTP server password
  -mailToAddress="received@email.com" # Receiver of notifications
```