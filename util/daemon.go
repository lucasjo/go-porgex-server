package util

import (
	"fmt"
	"log"
	"os/user"
	"path"
	"syscall"

	gozd "github.com/lucasjo/zero-downtime-daemon"
	"github.com/olebedev/config"
)

//zero-downtime-daemon context setting
func GetContext(cfg *config.Config, cmd string) gozd.Context {

	log.Println("porgex server exec : ", cmd)

	hash := cfg.UString("development.daemon.hash", "porgex-server-tcp-daemon")

	var command = ""

	if cmd != "" {
		command = cmd
	} else {
		command = cfg.UString("development.daemon.command", "start")
	}

	maxfds := syscall.Rlimit{Cur: 1000, Max: 1000}

	userNm := cfg.UString("development.daemon.user", "porgex")
	group := cfg.UString("development.daemon.group", "porgex")

	//logfile setting
	logfile := cfg.UString("development.daemon.logfile", "porgex-server-daemon.log")

	usr, err := user.Current()

	if err != nil {
		fmt.Println("err: ", err)
	}

	logfile = path.Join(usr.HomeDir, "logs", logfile)

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
		User:       userNm,
		Group:      group,
		Logfile:    logfile,
		Directives: directives,
	}

}
