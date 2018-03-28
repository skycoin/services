package model

import (
	"errors"
	"sync"

	"github.com/skycoin/services/otc/pkg/otc"
)

var ErrMissing = errors.New("orders missing")

type Lookup struct {
	sync.RWMutex

	Orders   map[string]*otc.Order
	Users    map[string]*otc.User
	Statuses map[string]*otc.User
}

func NewLookup() *Lookup {
	return &Lookup{
		Orders: make(map[string]*otc.Order),
		Users:  make(map[string]*otc.User),
	}
}

func (l *Lookup) AddStatus(user *otc.User) {
	l.Lock()
	defer l.Unlock()
	l.Statuses[string(user.Drop.Currency)+":"+user.Drop.Address] = user
}

func (l *Lookup) AddOrder(order *otc.Order) {
	l.Lock()
	defer l.Unlock()
	l.Orders[order.Id] = order
}

func (l *Lookup) AddUser(user *otc.User) {
	l.Lock()
	defer l.Unlock()
	l.Users[user.Id] = user
}

func (l *Lookup) GetStatus(id string) (*otc.User, error) {
	l.RLock()
	defer l.RUnlock()

	if l.Statuses[id] == nil {
		return nil, ErrMissing
	}

	return l.Statuses[id], nil
}

func (l *Lookup) GetOrder(id string) (*otc.Order, error) {
	l.RLock()
	defer l.RUnlock()

	if l.Orders[id] == nil {
		return nil, ErrMissing
	}

	return l.Orders[id], nil
}

func (l *Lookup) GetOrders() []*otc.Order {
	l.RLock()
	defer l.RUnlock()

	orders := make([]*otc.Order, len(l.Orders))
	for _, order := range l.Orders {
		orders = append(orders, order)
	}

	return orders
}

func (l *Lookup) GetUser(id string) (*otc.User, error) {
	l.RLock()
	defer l.RUnlock()

	if l.Users[id] == nil {
		return nil, ErrMissing
	}

	return l.Users[id], nil
}

func (l *Lookup) GetUsers() []*otc.User {
	l.RLock()
	defer l.RUnlock()

	users := make([]*otc.User, len(l.Users))
	for _, user := range l.Users {
		users = append(users, user)
	}

	return users
}
