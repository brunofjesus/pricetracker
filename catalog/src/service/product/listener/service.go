package listener

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/datasource"
	"github.com/brunofjesus/pricetracker/catalog/src/service/product/receiver"
)

var once sync.Once
var instance ProductListener
var lock sync.Mutex

type ProductListener interface {
	Listen() error
	Stop() error
	Channel() chan datasource.StoreProduct
}

type productListener struct {
	productChannel        chan datasource.StoreProduct
	listening             bool
	listenerSignalChannel chan bool
	productReceiver       receiver.ProductReceiver
}

func GetListener() ProductListener {
	once.Do(func() {
		instance = &productListener{
			productChannel:        make(chan datasource.StoreProduct, 10),
			listening:             false,
			listenerSignalChannel: make(chan bool),
			productReceiver:       receiver.GetProductReceiver(),
		}
	})

	return instance
}

func (l *productListener) Channel() chan datasource.StoreProduct {
	return l.productChannel
}

func (l *productListener) Listen() error {
	lock.Lock()
	defer lock.Unlock()

	if l.listening {
		return errors.New("already listening")
	}

	l.listening = true
	go l.listen()

	return nil
}

func (l *productListener) Stop() error {
	lock.Lock()
	defer lock.Unlock()

	if !l.listening {
		return errors.New("already stopped")
	}

	l.listening = false
	l.listenerSignalChannel <- true

	return nil
}

func (l *productListener) listen() {
	fmt.Println("Listening!")
	for {
		select {
		case <-l.listenerSignalChannel:
			return
		case storeProduct := <-l.productChannel:
			l.productReceiver.Receive(storeProduct)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
