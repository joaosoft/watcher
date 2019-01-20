package watcher

import (
	"fmt"
	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Watcher *WatcherConfig `json:"watcher"`
}

// WatcherConfig ...
type WatcherConfig struct {
	ReloadTime int64 `json:"reload_time"`
	Dirs       struct {
		Watch      []string `json:"watch"`
		Excluded   []string `json:"excluded"`
		Extensions []string `json:"extensions"`
	} `json:"dirs"`
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}