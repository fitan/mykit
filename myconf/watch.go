package myconf

import (
	"sync"
	"time"
)

type Watch struct {
	l      []chan struct{}
	t      *time.Ticker
	m      sync.Mutex
	reload bool
}

func (w *Watch) watch() {
	for {
		select {
		case <-w.t.C:
			w.m.Lock()
			if w.reload {
				for _, c := range w.l {
					c <- struct{}{}
				}
			}
			w.reload = false
			w.m.Unlock()
		}
	}
}

func (w *Watch) GetSignal() chan struct{} {
	w.m.Lock()
	defer w.m.Unlock()
	c := make(chan struct{})
	w.l = append(w.l, c)
	return c
}

func (w *Watch) Send() {
	w.m.Lock()
	defer w.m.Unlock()
	w.reload = true
}

func NewWatch(delayTime time.Duration) *Watch {
	w := &Watch{
		make([]chan struct{}, 0),
		time.NewTicker(delayTime),
		sync.Mutex{},
		false,
	}
	go w.watch()
	return w
}
