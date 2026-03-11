package wrapper

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	alpm "github.com/Jguer/dyalpm"
)

var Reset = "\033[0m"
var Red = "\033[1;31m"
var Green = "\033[1;32m"
var Yellow = "\033[1;33m"
var Blue = "\033[1;34m"
var Magenta = "\033[1;35m"
var Cyan = "\033[1;36m"
var Gray = "\033[1;37m"
var White = "\033[1;97m"
var builddest string

func StopHandle(trans alpm.Transaction) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nInterrupted, cleaning transaction...")
		trans.Release()
		os.Exit(1)
	}()
}

func InitHandle() *alpm.Handle {
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	handle.SetLogFile("/tmp/alpm.log")
	return &handle
}
