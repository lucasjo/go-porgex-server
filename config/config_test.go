package config

import (
	"testing"

	"github.com/lucasjo/go-porgex-server/config"
)

func TestGetConfig(t *testing.T) {
	//init()

	cfg := config.GetConfig("")

	if cfg == nil {
		t.FailNow()
	}

}
