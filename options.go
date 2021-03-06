package watcher

import (
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// WatcherOption ...
type WatcherOption func(watcher *Watcher)

// Reconfigure ...
func (w *Watcher) Reconfigure(options ...WatcherOption) {
	for _, option := range options {
		option(w)
	}
}

// WithConfiguration ...
func WithConfiguration(config *WatcherConfig) WatcherOption {
	return func(watcher *Watcher) {
		watcher.config = config
	}
}

// WithLogger ...
func WithLogger(logger logger.ILogger) WatcherOption {
	return func(watcher *Watcher) {
		watcher.logger = logger
		watcher.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) WatcherOption {
	return func(watcher *Watcher) {
		watcher.logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) WatcherOption {
	return func(watcher *Watcher) {
		watcher.pm = mgr
	}
}

// WithQuitChannel ...
func WithQuitChannel(quit chan int) WatcherOption {
	return func(watcher *Watcher) {
		watcher.quit = quit
	}
}

// WithReloadTime ...
func WithReloadTime(reloadTime int64) WatcherOption {
	return func(watcher *Watcher) {
		watcher.reloadTime = reloadTime
	}
}

// WithEventChannel ...
func WithEventChannel(event chan *Event) WatcherOption {
	return func(watcher *Watcher) {
		watcher.event = event
	}
}
