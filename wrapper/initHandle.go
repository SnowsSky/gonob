package wrapper

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	alpm "github.com/Jguer/dyalpm"
)

type Alpm struct {
	handle       *alpm.Handle
	localDB      alpm.Database
	syncDB       []alpm.Database
	syncDBsCache []alpm.Database
}

var GAlpm *Alpm

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
	GAlpm = &Alpm{
		handle: &handle,
	}

	return GAlpm.handle
}
