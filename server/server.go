package server

import (
	"bufio"
	"encoding/json"
	"net"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/lucasjo/go-porgex-server/config"
	"github.com/lucasjo/go-porgex-server/db"
	"github.com/lucasjo/go-porgex-server/models"
)

type Server struct {
	port          string
	host          string
	exit          chan bool
	joins         chan net.Conn
	clients       map[net.Conn]int
	clientCount   int
	message       chan interface{}
	removeclients chan net.Conn
}

type Metrix struct {
}

func NewServer() *Server {
	cfg := config.GetConfig("")
	_port, _ := cfg.String("development.tcp.port")
	_host, _ := cfg.String("development.tcp.host")
	return &Server{
		port:          _port,
		host:          _host,
		exit:          make(chan bool),
		joins:         make(chan net.Conn),
		clients:       make(map[net.Conn]int),
		clientCount:   0,
		message:       make(chan interface{}),
		removeclients: make(chan net.Conn),
	}
}

//server listen
func (s *Server) Listen() {
	log.Infof("listen")
	var err error
	var serveraddr *net.TCPAddr

	var addr = s.host + s.port
	serveraddr, err = net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		log.Fatalf("server ResolveTcp Fail\n", err)
	}

	listener, err := net.Listen("tcp", s.port)

	if err != nil {
		log.Fatalf("listnTcp fail port %s , %v\n", serveraddr, err)
	}

	go func() {
		select {
		case <-s.exit:
			listener.Close()
		}
	}()

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {

				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				log.Warningf("[collect-server] Failed to accept connection %s\n", err)
				continue
			}

			s.joins <- conn
		}

	}()

	for {

		//Handler 1) joins, 2) remove connection, 3) message send

		select {

		// Accept new Join
		case conn := <-s.joins:

			log.Infof("Accept new client, #%d", s.clientCount)
			s.clients[conn] = s.clientCount
			s.clientCount += 1

			go func(conn net.Conn, clientId int) {

				var req models.Request

				dec := json.NewDecoder(bufio.NewReader(conn))

				for dec.More() {
					err := dec.Decode(&req)

					if err != nil {
						log.Errorf("json decode error : %v\n", err)
						break
					}

					switch req.Service {
					case "memory":
						var memstatus models.MemStats
						err := json.Unmarshal(req.Data, &memstatus)
						if err != nil {
							log.Fatalf("error Memory Usage Json Convert: ", err)
						}

						s.message <- memstatus
					case "cpu":
						var cpustatus models.CPUStats
						err := json.Unmarshal(req.Data, &cpustatus)
						if err != nil {
							log.Fatalf("Error Cpu Usage Json Convert", err)
						}

						s.message <- cpustatus
					}
				}

			}(conn, s.clients[conn])

		case msg := <-s.message:

			log.Infof("msg %v , %v\n", msg, reflect.TypeOf(msg))
			go sendData(msg)

		case rConn := <-s.removeclients:
			rConn.Close()

		case <-s.exit:
			break

		}
	}

}

func sendData(v interface{}) {

	c := db.New(config.GetConfig(""))
	coll := db.GetColl(c)

	log.Infof("coll : %v\n", coll)
	log.Infof("v : %v\n", v)

	//vType := reflect.TypeOf(v)

	var collection *mgo.Collection

	switch v.(type) {
	case models.MemStats:
		log.Infof("models.MemStats")
		collection = coll.MemUsageCollection
	case models.CPUStats:
		log.Infof("models.CpuStats")
		collection = coll.CpuUsageCollection
	}

	err := db.Save(collection, v)

	if err != nil {
		log.Errorf("Memory Data Insert Error : ", err)
	}
	log.Infof("MemoryDataInsert : %v\n", v)

}

func (s *Server) Stop() {

	log.Infof("Server Stop %s:%s\n", s.host, s.port)
	close(s.exit)
}

func (s *Server) checkUp() {
	log.Infof("server alive")
}

/*
func sendData(line []byte) {
	var m Metrix

	err := json.Unmarshal(line, &m)
	if err != nil {
		log.Errorf("json convert fail %v\n", err)
	}

	log.Infof("message Data %v\n", &m)

}
*/
