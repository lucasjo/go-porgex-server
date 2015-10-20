package main

import (
	"fmt"

	"github.com/lucasjo/go-porgex-server/config"
	"github.com/lucasjo/go-porgex-server/db"
)

func main() {
	c := db.New(config.GetConfig("/Users/kikimans/go/src/github.com/lucasjo/go-porgex-server/porgex-server.yaml"))
	fmt.Printf("connect %v\n", c)

}
