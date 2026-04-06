package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"
	"os"
	"os/signal"
	"syscall"

	alpm "github.com/Jguer/dyalpm"
	"github.com/Morganamilo/go-pacmanconf"
	scolor "github.com/SnowsSky/scolor/pkg"
)

type Alpm struct {
	handle *alpm.Handle
	/*localDB      alpm.Database*/
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
	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string"))
		scolor.BoldWhite.DisplayText(err.Error() + "\n")
		return nil
	}
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	handle.SetCheckSpace(conf.CheckSpace)
	handle.SetCacheDirs(conf.CacheDir)
	handle.SetParallelDownloads(5)
	handle.SetLogFile("/tmp/alpm.log")
	handle.SetEventCallbackFunc(EventCallback)
	GAlpm = &Alpm{
		handle: &handle,
	}

	return GAlpm.handle
}
