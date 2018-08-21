# watcher
[![Build Status](https://travis-ci.org/joaosoft/watcher.svg?branch=master)](https://travis-ci.org/joaosoft/watcher) | [![codecov](https://codecov.io/gh/joaosoft/watcher/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/watcher) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/watcher)](https://goreportcard.com/report/github.com/joaosoft/watcher) | [![GoDoc](https://godoc.org/github.com/joaosoft/watcher?status.svg)](https://godoc.org/github.com/joaosoft/watcher)

A simple cross-platform file watcher

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Multi directories
* Exclusions
* Extensions

## Events
* OperationCreate, when added a new file
* OperationUpdate, when updated an existing file
* OperationDelete, when deleted a file
* OperationChanges, all files loaded and has some of the previous events

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/watcher
```

## Usage 
This examples are available in the project at [watcher/main/main.go](https://github.com/joaosoft/watcher/tree/master/main/main.go)
```
import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	github.com/joaosoft/watcher
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
```


> Configuration file
```
{
  "watcher": {
    "reload_time": 1,
    "dirs": {
      "watch":[ "examples/" ],
      "excluded":[ "examples/test_2" ],
      "extensions": [ "go", "json" ]
    },
    "log": {
      "level": "error"
    }
  },
  "manager": {
    "log": {
      "level": "error"
    }
  }
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
