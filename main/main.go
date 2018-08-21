package main

import (
	"fmt"
	"watcher"
)

func main() {
	c := make(chan *watcher.Event)
	w, err := watcher.NewWatcher(watcher.WithEventChannel(c))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event := <-c:
				fmt.Printf("received event %+v\n", event)
			}
		}
	}()

	if err := w.Start(); err != nil {
		panic(err)
	}
}
