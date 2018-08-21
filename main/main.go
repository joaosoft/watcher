package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"watcher"
)

func main() {
	quit := make(chan int)
	c := make(chan *watcher.Event)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	w := watcher.NewWatcher(watcher.WithEventChannel(c), watcher.WithQuitChannel(quit))

	go func() {
		for {
			select {
			case <-termChan:
				return
			case <-quit:
				return
			case event := <-c:
				fmt.Printf("received event %+v\n", event)
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := w.Start(&wg); err != nil {
		panic(err)
	}

	<-termChan
	quit <- 1
}
