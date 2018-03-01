package buyer

import (
	"container/list"

	"github.com/skycoin/services/otc/exchange"
	"github.com/skycoin/services/otc/types"
)

// ExchangeDeposit sends currency from OTC wallet to exchange deposit address
// so that it can later be traded for SKY on the exchange.
func (b *Buyer) ExchangeDeposit(w *types.Work, e *list.Element) {
	// get balance of drop
	amount, err := b.dropper.GetBalance(w.Request.Currency, w.Request.Drop)
	if err != nil {
		w.Return(err)
		b.work.Remove(e)
		return
	}

	// TODO: get deposit address from multiple exchanges for multiple
	// currencies
	drop, err := exchange.GetBTCDepositAddress()
	if err != nil {
		w.Return(err)
		b.work.Remove(e)
		return
	}

	// send from drop to exchange
	hash, err := b.dropper.Connections[w.Request.Currency].Send(drop, amount)
	if err != nil {
		w.Return(err)
		b.work.Remove(e)
		return
	}

	// next state
	w.Request.Metadata.BuyDrop = drop
	w.Request.Metadata.BuyStatus = types.EXCHANGE_CONFIRM
	w.Request.Metadata.BuyHashTo = hash
	w.Return(nil)
}

// ExchangeConfirm confirms that the transaction from OTC->Exchange has
// finished and the amount can be used to execute a trade.
func (b *Buyer) ExchangeConfirm(w *types.Work, e *list.Element) {
	hash := w.Request.Metadata.BuyHashTo

	// check if confirmed yet
	confirmed, err := b.dropper.Connections[w.Request.Currency].Confirmed(hash)
	if err != nil {
		w.Return(err)
		b.work.Remove(e)
		return
	}

	// if confirmed send to next step
	if confirmed {
		w.Request.Metadata.BuyStatus = types.EXCHANGE_TRADE
		w.Return(nil)
	} else {
		w.Return(nil)
		return
	}
}

// ExchangeTrade executes a trade on the exchange.
func (b *Buyer) ExchangeTrade(w *types.Work, e *list.Element) {}

// TODO
func (b *Buyer) ExchangeReturn(w *types.Work, e *list.Element) {}

// TODO
func (b *Buyer) ExchangeReturned(w *types.Work, e *list.Element) {}
