package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	config     interface{}
	configLock = new(sync.RWMutex)
)

func init() {
	loadConfig(true)
}

func start() {
	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')

		if strings.Split(message, "\n")[0] == "refresh" {
			log.Println("Reloaded")
			loadConfig(true)
		}
	}
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
	go start()

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		log.Println("Config: ", GetConfig())
	}
}
