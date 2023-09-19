package receiver

import (
	"sync"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
)

var once sync.Once
var instance ProductReceiver

type ProductReceiver interface {
	Receive(storeProduct datasource.StoreProduct)
}

type productReceiver struct {
}

func GetReceiver() ProductReceiver {
	once.Do(func() {
		instance = &productReceiver{}
	})
	return instance
}

// Receive implements Receiver.
func (*productReceiver) Receive(storeProduct datasource.StoreProduct) {
	panic("unimplemented")
}
