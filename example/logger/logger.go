package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		//DisableTimestamp: false,
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})
	//log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

}

func main() {
	log.Info("te")
}
