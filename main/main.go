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
	w := watcher.NewWatcher(watcher.WithEventChannel(c), watcher.WithQuitChannel(quit))

	go func() {
		for {
			select {
			case event := <-c:
				fmt.Printf("received event %+v\n", event)
			}
		}
	}()

	if err := w.Start(&sync.WaitGroup{}); err != nil {
		panic(err)
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	<-termChan
	quit <- 1
}
