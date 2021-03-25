package goconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/zsounder/zgo/logger"
)

type FuncConfigLoader func(string) error

var ConfigLoaders map[string]FuncConfigLoader = make(map[string]FuncConfigLoader, 10)

var (
	configDir  string
	configKeys []string
	loadLock   sync.RWMutex
	inLoading  bool
)

func Init(dir string, configs []string) {
	if !strings.HasSuffix(dir, "/") {
		configDir = dir + "/"
	}
	configKeys = configs
	Load()
	go watch()
}

func ReadBegin() {
	loadLock.RLock()
}

func AddLoader(key string, loader FuncConfigLoader) {
	fmt.Printf("key:%s", key)
	ConfigLoaders[key] = loader
}

func ReadEnd() {
	loadLock.RUnlock()
}

func Load() {
	if inLoading {
		return
	}

	loadLock.Lock()
	inLoading = true
	defer func() {
		loadLock.Unlock()
		inLoading = false
	}()

	for _, config := range configKeys {
		if config == "*" {
			filepath.Walk(configDir,
				func(path string, f os.FileInfo, err error) error {
					if f == nil {
						return err
					}
					if f.IsDir() {
						logger.Fatal("dir:", path)
						return nil
					}

					if loader, ok := ConfigLoaders["*"]; ok {
						if err := loader(configDir + path); err != nil {
							logger.Fatalf("config:%s loading with error: %s", path, err.Error())
						}
					} else {
						logger.Fatal("config loader not exist:", path)
					}
					logger.Config("config load success:", path)

					return nil
				})
		} else {
			if loader, ok := ConfigLoaders[config]; ok {
				if err := loader(configDir + config); err != nil {
					logger.Fatalf("config:%s loading with error: %s", config, err.Error())
				}
			} else {
				logger.Fatal("config loader not exist:", config)
			}
			logger.Config("config load success:", config)
		}
	}
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					Load()
				}
			case err := <-watcher.Errors:
				logger.Error("error:", err)
			}
		}
	}()

	if err = watcher.Add(configDir); err != nil {
		logger.Fatal(err)
	}

	<-done
}
