package config

import (
	"testing"

	"github.com/lucasjo/go-porgex-server/config"
)

func TestGetConfig(t *testing.T) {
	//init()

	cfg := config.GetConfig("/Users/kikimans/go/src/github.com/lucasjo/go-porgex-server/porgex_server.yaml")

	if cfg == nil {
		t.FailNow()
	}

}
