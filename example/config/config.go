package main

import (
	"fmt"
	"io/ioutil"
)

//"github.com/olebedev/config"

const ConfigName = "porgex_server.yaml"

func main() {
	cfg, err := ioutil.ReadFile("/Users/kikimans/go/src/github.com/lucasjo/go-porgex-server/porgex_server.yaml")
	if err != nil {
		fmt.Println("err : ", err)
	}

	fmt.Println("cfg : ", cfg)
}
