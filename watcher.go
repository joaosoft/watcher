package watcher

import (
	"fmt"
	"os"
	"sync"

	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Watcher struct {
	config        *WatcherConfig
	watch         []string
	excluded      []string
	extensions    []string
	files         map[string]map[string]FileInfo
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
	logger        logger.ILogger
	reloadTime    int64
	quit          chan int
	event         chan *Event
	started       bool
}

func NewWatcher(options ...WatcherOption) *Watcher {
	config, simpleConfig, err := NewConfig()
	watcher := &Watcher{
		watch:      make([]string, 0),
		excluded:   make([]string, 0),
		extensions: make([]string, 0),
		reloadTime: 1,
		files:      make(map[string]map[string]FileInfo),
		pm:         manager.NewManager(manager.WithRunInBackground(true)),
		logger:     logger.NewLogDefault("watcher", logger.InfoLevel),
		event:      make(chan *Event),
		quit:       make(chan int),
		config:     &WatcherConfig{},
	}

	if watcher.isLogExternal {
		watcher.pm.Reconfigure(manager.WithLogger(watcher.logger))
	}

	if err != nil {
		logger.Error(err.Error())
	} else {
		watcher.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Watcher.Log.Level)
		logger.Debugf("setting log level to %s", level)
		logger.Reconfigure(logger.WithLevel(level))
	}

	// loading each configuration
	watcher.reloadTime = watcher.config.ReloadTime
	watcher.watch = append(watcher.watch, watcher.config.Dirs.Watch...)
	watcher.excluded = append(watcher.excluded, watcher.config.Dirs.Excluded...)
	watcher.extensions = append(watcher.extensions, watcher.config.Dirs.Extensions...)

	watcher.Reconfigure(options...)

	return watcher
}

func (w *Watcher) AddWatch(watchs ...string) *Watcher {
	w.watch = append(w.watch, watchs...)
	return w
}

func (w *Watcher) AddExtension(extensions ...string) *Watcher {
	w.extensions = append(w.extensions, extensions...)
	return w
}

func (w *Watcher) AddExcluded(excluded ...string) *Watcher {
	w.excluded = append(w.excluded, excluded...)
	return w
}

// execute ...
func (w *Watcher) execute() error {
	w.logger.Debugf("executing watcher for watch %+v", w.watch)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	// load
	for _, dir := range w.watch {

		_, err := os.Stat(dir)
		if err != nil {
			return err
		}

		go func() {
			for {
				select {
				case <-termChan:
					w.logger.Info("received term signal")
					return
				case <-w.quit:
					w.logger.Info("received shutdown signal")
					return
				case <-time.After(time.Duration(w.reloadTime) * time.Second):
					w.mux.Lock()
					changed := false
					w.logger.Info("reloading data")

					// copy before reloadTime files
					oldFiles := w.files[dir]
					w.files[dir] = make(map[string]FileInfo)

					if err = w.doLoad(oldFiles, dir, dir, &changed); err != nil {
						w.quit <- 1
					}

					if err = w.doRemove(dir, oldFiles, &changed); err != nil {
						w.quit <- 1
					}

					if changed {
						w.event <- &Event{
							File:      dir,
							Operation: OperationChanges,
						}
					}
					w.mux.Unlock()
				}
			}
		}()
	}

	return nil
}

// Start ...
func (w *Watcher) Start(wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}

	if err := w.pm.Start(); err != nil {
		return err
	}

	if err := w.execute(); err != nil {
		return err
	}

	w.started = true

	return nil
}

// Stop ...
func (w *Watcher) Stop(wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}

	if err := w.pm.Stop(); err != nil {
		return err
	}

	w.started = false

	return nil
}

// Started ...
func (w *Watcher) Started() bool {
	return w.started
}

func (w *Watcher) doLoad(oldFiles map[string]FileInfo, dir string, next string, changed *bool) error {
	fileInfo, err := os.Stat(next)
	if err != nil {
		return err
	}

	if fileInfo.Name() != "." && strings.HasPrefix(fileInfo.Name(), ".") {
		return nil
	}

	// if it is a directory
	if fileInfo.IsDir() {

		// exclude validation
		for _, exclude := range w.excluded {
			if strings.HasPrefix(next, exclude) {
				return nil
			}
		}

		w.logger.Debugf("loading files on directory [%s]", next)

		subDir, err := filepath.Glob(fmt.Sprintf("%s/*", next))
		if err != nil {
			w.logger.Errorf("error reading directory %s", err)
			return err
		}
		for _, nextDir := range subDir {
			w.logger.Debugf("loading subdirectory [%s]", nextDir)
			if err := w.doLoad(oldFiles, dir, nextDir, changed); err != nil {
				return err
			}
		}
		return nil
	}

	// extension validation
	if index := strings.LastIndex(next, "."); index > 0 {
		allowed := false
		for _, extension := range w.extensions {
			if extension == next[index+1:] {
				allowed = true
			}
		}

		if !allowed {
			return nil
		}
	}

	w.logger.Debugf("loading file [%s]", fileInfo.Name())

	// if it is a file
	w.files[dir][next] = FileInfo{
		FullName: next,
		Name:     fileInfo.Name(),
		Size:     fileInfo.Size(),
		ModTime:  fileInfo.ModTime(),
	}

	if oldFileInfo, ok := oldFiles[next]; !ok {
		// new file
		w.logger.Debugf("added a new file on directory [%s]", next)
		w.event <- &Event{
			File:      next,
			Operation: OperationCreate,
		}
		*changed = Changed
	} else {
		if oldFileInfo.ModTime != fileInfo.ModTime() ||
			oldFileInfo.Size != fileInfo.Size() {
			// updated file
			w.logger.Debugf("changed file on directory [%s]", next)
			w.event <- &Event{
				File:      next,
				Operation: OperationUpdate,
			}
			*changed = Changed
		}
	}

	return nil
}

func (w *Watcher) doRemove(dir string, oldFiles map[string]FileInfo, changed *bool) error {

	for fullName, _ := range oldFiles {
		if _, ok := w.files[dir][fullName]; !ok {
			delete(w.files[dir], fullName)
			w.logger.Debugf("deleted file on directory [%s]", dir)
			w.event <- &Event{
				File:      fullName,
				Operation: OperationDelete,
			}
			*changed = Changed
		}
	}

	return nil
}
