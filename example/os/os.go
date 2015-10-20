package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

func main() {
	fmt.Println("tempdir: ", os.TempDir())

	usr, err := user.Current()

	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("homedir: ", usr.HomeDir)
	var li syscall.Rlimit

	err1 := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &li)

	if err1 != nil {
		fmt.Println(err)
	}

	fmt.Println("lim: ", li)
}
