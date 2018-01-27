## How to use

### Step 1

Add import to your application

`_ "github.com/skycoin/viscript/signal"`

This is all the code you need to change.

### Step 2

And then add environment variable to the command line

```
run as a signal client(default:false)
SIGNAL_CLIENT

connect to the address of signal server (default:"localhost:7999")
SIGNAL_SERVER_ADDRESS

client id to use(default:1)
SIGNAL_CLIENT_ID
```

for example:

`~$ SIGNAL_CLIENT=true ./cx --repl`