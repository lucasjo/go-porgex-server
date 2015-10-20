package logger

import (
	"io"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/howeyc/fsnotify"
)

var std = NewFileLogger()

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	signalChan := make(chan os.Signal, 16)
	signal.Notify(signalChan, syscall.SIGHUP)

	go func() {
		for {
			select {
			case <-signalChan:
				err := std.Reopen()
				log.Infof("HUP collect server, reopen log %#v", std.Filename())
				if err != nil {
					log.Errorf("Reopen %v failed, %s", std.Filename(), err.Error())

				}
			}
		}
	}()

}

type FileLogger struct {
	sync.RWMutex
	filename  string
	fd        *os.File
	watchDone chan bool
}

//instance FileLogger
func NewFileLogger() *FileLogger {
	return &FileLogger{
		filename:  "",
		fd:        nil,
		watchDone: nil,
	}
}

func (fl *FileLogger) Open(filename string) error {
	fl.Lock()
	fl.filename = filename
	fl.Unlock()

	reopenErr := fl.Reopen()

	if fl.watchDone != nil {
		close(fl.watchDone)
	}

	fl.watchDone = make(chan bool)
	fl.fsWatch(fl.filename, fl.watchDone)

	return reopenErr
}

func (fl *FileLogger) fsWatch(filename string, q chan bool) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Warningf("fsnotify.NewWatcher() : %v\n", err)
		return
	}

	if filename == "" {
		log.Errorf("fsWatch Func parameter Not Found filename %v", filename)
		return
	}

	subscribe := func() {
		if err := watcher.WatchFlags(filename, fsnotify.FSN_CREATE|fsnotify.FSN_DELETE|fsnotify.FSN_RENAME); err != nil {
			log.Warningf("fsnotify watcher watch(%s) : %s\n", filename, err)
		}
	}

	subscribe()

	go func() {
		defer watcher.Close()

		for {
			select {
			case <-watcher.Event:
				fl.Reopen()
				subscribe()

				log.Infof("Reopen log %v by fsnotify event \n", std.Filename())
				if err != nil {
					log.Errorf("Reopen log %v failed : %s\n", std.Filename(), err)
				}

			case <-q:
				return
			}
		}
	}()
}

func (fl *FileLogger) Reopen() error {
	fl.Lock()
	defer fl.Unlock()

	var newFd *os.File
	var err error

	if fl.filename != "" {
		newFd, err = os.OpenFile(fl.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			return err
		}
	} else {
		newFd = nil
	}

	oldFd := fl.fd
	fl.fd = newFd

	var loggerOut io.Writer

	if fl.fd != nil {
		loggerOut = fl.fd
	} else {
		loggerOut = os.Stderr
	}

	log.SetOutput(loggerOut)

	if oldFd != nil {
		oldFd.Close()
	}

	return nil

}

func (fl *FileLogger) Filename() string {
	fl.RLock()
	defer fl.RUnlock()

	return fl.filename
}

func SetFile(filename string) error {
	return std.Open(filename)
}

func SetLevel(lv string) error {
	level, err := log.ParseLevel(lv)
	if err != nil {
		return err
	}

	log.SetLevel(level)

	return nil
}

func PrepareFile(filename string, owner *user.User) error {
	if filename == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if fd != nil {
		fd.Close()
	}

	if err != nil {
		return err
	}

	if err := os.Chmod(filename, 0644); err != nil {
		return err
	}

	if owner != nil {
		uid, err := strconv.ParseInt(owner.Uid, 10, 00)

		if err != nil {
			return err
		}

		gid, err := strconv.ParseInt(owner.Gid, 10, 0)
		if err != nil {
			return err
		}

		if err := os.Chown(filename, int(uid), int(gid)); err != nil {
			return err
		}
	}

	return nil
}
