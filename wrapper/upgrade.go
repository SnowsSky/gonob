package wrapper

import (
	"fmt"

	"github.com/Jguer/go-alpm/v2"
	"github.com/Morganamilo/go-pacmanconf"
)

func Upgrade() {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := h.Release(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, repo := range conf.Repos {
		fmt.Println(repo)
		db, err := h.RegisterSyncDB(repo.Name, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.SetServers(repo.Servers)
	}
}
