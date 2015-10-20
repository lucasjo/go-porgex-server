package util

import (
	"os"
	"path"
	"syscall"

	"github.com/olebedev/config"
	gozd "github.com/tomasen/zero-downtime-daemon"
)

//zero-downtime-daemon context setting
func GetContext(cfg *config.Config, cmd string) gozd.Context {

	hash := cfg.UString("development.daemon.hash", "porgex-server-tcp-daemon")
	var command string

	if cmd != nil {
		command = cmd
	} else {
		command = cfg.UString("development.daemon.command", "start")
	}

	maxfds := syscall.Rlimit{Cur: 32677, Max: 32677}

	user := cfg.UString("development.daemon.user", "porgex")
	group := cfg.UString("development.daemon.group", "porgex")

	//logfile setting
	logfile := cfg.UString("development.daemon.logfile", "porgex-server-daemon.log")
	logfile = path.Join(os.TempDir(), logfile)

	directives := make(map[string]gozd.Server)

	network := cfg.UString("development.daemon.directives.network", "tcp")
	address := cfg.UString("development.daemon.directives.address", "127.0.0.1:3001")

	directives["port1"] = gozd.Server{
		Network: network,
		Address: address,
	}

	return gozd.Context{
		Hash:       hash,
		Command:    command,
		Maxfds:     maxfds,
		User:       user,
		Group:      group,
		Logfile:    logfile,
		Directives: directives,
	}

}
