package main

import (
	"net"
	log "github.com/Sirupsen/logrus"
	)

type TCP struct {
	exit     chan bool
	listener *net.TCPListener
	//active int32
	//errors uint32

}

func NewTCP() *TCP {
	return &TCP{
		exit: make(chan bool),
	}
}

func (t *TCP) stop() {
	close(t.exit)
}

// listen bind :port data out 
func (t *TCP) listen(addr *net.TCPAddr) {

	var err error
	t.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	go func(){
		for{
			select {
			case <- t.exit :
				t.listener.Close()
				return
			}
		}
	}

	handler := t.rquesthandler

	go func(){
		defer t.listener.Close()

		for{
			conn, err := t.listener.Accept()

			if err := nil {
				log.Fatal(...)
				break
			}
		}
	}
	

}
