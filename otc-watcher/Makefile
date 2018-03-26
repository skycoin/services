run:
	go build
	./otc-watcher -rpc_node="localhost:8332" -rpc_user="aaa" -rpc_pass="123" -wallet_account="1" -wallet_pass="1234"
format:
	goimports -w -local github.com/skycoin/services/otc-watcher ./pkg/api
	goimports -w -local github.com/skycoin/services/otc-watcher ./pkg/currency
	goimports -w -local github.com/skycoin/services/otc-watcher ./pkg/scanner