package management

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/brunofjesus/pricetracker/catalog/src/store"
)

var once sync.Once
var instance Listener
var lock sync.Mutex

type Listener struct {
	ProductChannel        chan store.StoreProduct
	listening             bool
	listenerSignalChannel chan bool
}

func GetListener() *Listener {
	once.Do(func() {
		instance = Listener{
			ProductChannel:        make(chan store.StoreProduct, 10),
			listening:             false,
			listenerSignalChannel: make(chan bool),
		}
	})

	return &instance
}

func (l *Listener) Start() error {
	lock.Lock()
	defer lock.Unlock()

	if l.listening == true {
		return errors.New("already listening")
	}

	l.listening = true
	go l.listen()

	return nil
}

func (l *Listener) Stop() error {
	lock.Lock()
	defer lock.Unlock()

	if l.listening == false {
		return errors.New("already stopped")
	}

	l.listening = false
	l.listenerSignalChannel <- true

	return nil
}

func (l *Listener) listen() {
	fmt.Println("Listening!")
	for {
		select {
		case <-l.listenerSignalChannel:
			return
		case storeProduct := <-l.ProductChannel:
			fmt.Printf("%v\n", storeProduct.Name)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
