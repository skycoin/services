package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/namsral/flag"
	transactionsMonitor "github.com/skycoin/services/pending-transactions-monitor/logic"
)

func main() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	nodeAddress := flag.String("nodeAddress", "http://127.0.0.1:6420", "Path to the Skycoin node API")
	mailHost := flag.String("mailHost", "smtp.gmail.com:587", "SMTP server")
	mailToAddress := flag.String("mailToAddress", "what.if.do.best@gmail.com", "From address for email")
	mailUsername := flag.String("mailUsername", "testtempmail07@gmail.com", "SMTP server user")
	mailPassword := flag.String("mailPassword", "sdksdk123", "SMTP server password")
	pendingTime := flag.String("pendingTime", "60", "Max pending transaction time (seconds, default = 60)")

	flag.Parse()

	toAddr := *mailToAddress
	pendingTimeParsed, err := strconv.Atoi(*pendingTime)
	if err != nil {
		fmt.Println("Invalid pendingTime parameter: ", *pendingTime, "\n", err)
	}

	mailer := transactionsMonitor.NewMailer(*mailHost, *mailUsername, *mailPassword, *mailToAddress)

	monitor := transactionsMonitor.NewMonitor(*nodeAddress)
	transactions, err := monitor.GetPendingTransactions()
	if err != nil {
		fmt.Println("main > Error (monitor.GetPendingTransactions):", err)

		emailBody :=
			fmt.Sprintf("Error has occurred in the transactions monitor: <br/> main > Error (monitor.GetPendingTransactions): <br/>%s", err)
		err = mailer.SendMail(&transactionsMonitor.Letter{
			Body:    emailBody,
			Subject: "[ALERT] Pending transactions > Service error",
			To:      toAddr,
		})
		if err != nil {
			fmt.Println("main > Error (mailer.SendMail):", err)
		}
		return
	}

	pendigTransactionTxs := ""
	pendigTransactionsCount := 0
	oldestTimestamp := time.Time{}
	for _, t := range transactions {
		if time.Now().After(t.Received.Add(time.Duration(pendingTimeParsed) * time.Second)) {
			if pendigTransactionsCount < 100 {
				pendigTransactionTxs = pendigTransactionTxs + "<br/>" + t.Transaction.Txid
			}
			pendigTransactionsCount = pendigTransactionsCount + 1
			if oldestTimestamp.Before(t.Received) {
				oldestTimestamp = t.Received
			}
		}
	}

	if pendigTransactionTxs == "" {
		fmt.Println("There are no pending transactions")
		return
	}

	emailBody :=
		fmt.Sprintf("Total count transactions count: %d. <br/> Oldest timestamp: %s (%s) <br/>txids of pending transactions: <br/><br/>%s",
			pendigTransactionsCount,
			oldestTimestamp,
			time.Now().Sub(oldestTimestamp).Round(time.Second),
			pendigTransactionTxs)
	err = mailer.SendMail(&transactionsMonitor.Letter{
		Body:    emailBody,
		Subject: "[ALERT] Pending transactions",
		To:      toAddr,
	})
	if err != nil {
		fmt.Println("main > Error (mailer.SendMail):", err)
	}
}
