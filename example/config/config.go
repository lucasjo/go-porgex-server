package main

import (
	"fmt"

	"github.com/lucasjo/go-porgex-server/config"
	//"github.com/olebedev/config"
)

const ConfigName = "porgex-server.yaml"

func main() {
	cfg := config.GetConfig("")
	port, _ := cfg.String("development.tcp.port")
	fmt.Printf("config %v\n", port)
}
