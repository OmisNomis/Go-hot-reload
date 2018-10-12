package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	config     interface{}
	configLock = new(sync.RWMutex)
)

func init() {
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	// signal.Notify(s, syscall.SIGUSR2, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGQUIT)
	// signal.Notify(s)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Printf("Reloaded using: %+v", s)
		}
	}()
}

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	var data interface{}
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	configLock.Lock()
	config = data
	configLock.Unlock()
}

func GetConfig() interface{} {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		log.Println("Config: ", GetConfig())
	}
}
