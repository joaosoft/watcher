package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"watcher"
)

func main() {
	c := make(chan *watcher.Event)
	quit := make(chan int)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	w := watcher.NewWatcher(watcher.WithEventChannel(c), watcher.WithQuitChannel(quit))

	go func() {
		for {
			select {
			case <-termChan:
				quit <- 1
				return
			case event := <-c:
				fmt.Printf("received event %+v\n", event)
			}
		}
	}()

	if err := w.Start(); err != nil {
		panic(err)
	}

	<-termChan
	if err := w.Stop(); err != nil {
		panic(err)
	}
}
