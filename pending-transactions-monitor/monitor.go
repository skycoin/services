package main

import (
	"fmt"
	"time"

	"github.com/namsral/flag"
	transactionsMonitor "github.com/skycoin/services/pending-transactions-monitor/src"
)

func main() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	nodeAddress := flag.String("nodeAddress", "http://127.0.0.1:6420", "Path to the Skycoin node API")
	mailHost := flag.String("mailHost", "smtp.gmail.com:587", "SMTP server")
	mailToAddress := flag.String("mailToAddress", "what.if.do.best@gmail.com", "From address for email")
	mailUsername := flag.String("mailUsername", "testtempmail07@gmail.com", "SMTP server user")
	mailPassword := flag.String("mailPassword", "sdksdk123", "SMTP server password")

	flag.Parse()

	mailer := transactionsMonitor.NewMailer(*mailHost, *mailUsername, *mailPassword, *mailToAddress)

	monitor := transactionsMonitor.NewMonitor(*nodeAddress)
	transactions, err := monitor.GetPendingTransactions()
	if err != nil {
		fmt.Println(err)
	}

	pendigTransactionTxs := ""
	for _, t := range transactions {
		if time.Now().After(t.Received.Add(time.Minute)) {
			fmt.Println(t.Received)
			pendigTransactionTxs = pendigTransactionTxs + "<br/>" + t.Transaction.Txid
		}
	}

	toAddr := *mailToAddress
	emailBody :=
		fmt.Sprintf("txids of pending transactions: <br/><br/>%s", pendigTransactionTxs)
	mailer.SendMail(&transactionsMonitor.Letter{
		Body:    emailBody,
		Subject: "Skycoin > Pending transactions",
		To:      toAddr,
	})
}
