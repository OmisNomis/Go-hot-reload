package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
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
	ln, err := net.Listen("unix", "/tmp/go.sock")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go echoServer(fd)
	}
}

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Writing client error: ", err)
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
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		log.Println("Config: ", GetConfig())
	}
}
