package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/lucasjo/go-porgex-node/models"
	lconfig "github.com/lucasjo/go-porgex-server/config"
	"github.com/lucasjo/go-porgex-server/db"
	"github.com/lucasjo/go-porgex-server/util"
	gozd "github.com/lucasjo/zero-downtime-daemon"
)

var (
	signal = flag.String("c", "", `send signal to the porgex-server
			kill - graceful shutdown
			stop - fast shutdown
			reload - reloading the configuration file
			reopen - reopening service
			start - service start`)

	configPtr = flag.String("config", "", "config option root directory filename or filepath")

	configpath = ""
	sg         = ""
)

var message = make(chan interface{})

func saveData(v interface{}) {

	dbc := db.New(lconfig.GetConfig(configpath))

	collections := db.GetColl(dbc)

	var collection *mgo.Collection

	switch v.(type) {
	case models.MemStats:
		collection = collections.MemUsageCollection
	case models.CPUStats:
		collection = collections.CpuUsageCollection
	}

	err := db.Save(collection, v)

	if err != nil {
		log.Fatalf("Usage Data Save Error : %v\n", err)
	}

	log.Printf("Usage Data Insert : %v, %v\n", reflect.TypeOf(v), v)
}

func handlerConnection(conn net.Conn) {

	for {
		go func(conn net.Conn) {
			var req models.Request
			dec := json.NewDecoder(bufio.NewReader(conn))

			for dec.More() {
				err := dec.Decode(&req)

				if err != nil {
					log.Fatalf("json decode error : %v\n", err)
					break
				}

				switch req.Service {
				case "memory":
					var memStatus models.MemStats
					err := json.Unmarshal(req.Data, &memStatus)
					if err != nil {
						log.Fatalf("Error Memory Usage Json Convert : %v\n", err)
					}
					message <- memStatus

				case "cpu":
					var cpuStatus models.CPUStats
					err := json.Unmarshal(req.Data, &cpuStatus)
					if err != nil {
						log.Fatalf("Error Cpu Usage Json Convert : %v\n", err)
					}
					message <- cpuStatus
				}

			}
		}(conn)

		select {
		case msg := <-message:
			go saveData(msg)
		}

	}

}

func serverTcp(cl chan net.Listener) {
	for v := range cl {
		defer v.Close()
		go func(l net.Listener) {
			log.Println("porgex-server-listener : ", reflect.ValueOf(l).Elem().FieldByName("Name").String())

			for {
				conn, err := l.Accept()

				if err != nil {
					log.Fatalf("Accept Error : %v\n", err)
					break
				}

				go handlerConnection(conn)

			}
		}(v)

	}
}

func main() {

	flag.Parse()

	configpath = *configPtr
	sg = *signal

	if configpath == "" {
		dir, _ := os.Getwd()
		configpath = strings.Join([]string{dir, "porgex_server.yaml"}, "/")
	} else {
		if !util.Exists(configpath) {
			fmt.Printf("no such file or directory %s\n", configpath)
			os.Exit(1)
		}
	}

	cfg := lconfig.GetConfig(configpath)

	cntxt := util.GetContext(cfg, sg)

	cl := make(chan net.Listener, 1)

	go serverTcp(cl)

	sig, err := gozd.Daemonize(cntxt, cl)

	if err != nil {
		log.Println("error : ", err)
		return
	}

	for s := range sig {
		switch s {
		case syscall.SIGHUP, syscall.SIGUSR2:
			// other job ok
			log.Println("Porgex Server Start ", time.Now())

		case syscall.SIGTERM:
			log.Println("Porgex Server Stop ", time.Now())
			cli := <-cl
			cli.Close()
		}
	}

}
