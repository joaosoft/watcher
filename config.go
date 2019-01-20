package watcher

import (
	"fmt"

	manager "github.com/joaosoft/manager"
	"github.com/labstack/gommon/log"
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
func NewConfig() (*WatcherConfig, error) {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", GetEnv()), appConfig); err != nil {
		log.Error(err.Error())

		return &WatcherConfig{}, err
	}

	return appConfig.Watcher, nil
}
