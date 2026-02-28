package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"

	alpm "github.com/Jguer/dyalpm"
	"github.com/Morganamilo/go-pacmanconf"
)

func InitSyncDatabases(handle *alpm.Handle) []alpm.Database {
	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + " : " + translations.Translate("alpm_unable_to_fetch-_syncdbs") + Reset)
		return nil
	}

	for _, repo := range conf.Repos {
		db, err := (*handle).RegisterSyncDB(repo.Name, 0)
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + " : " + translations.Translate("alpm_unable_to_fetch-_syncdbs") + Reset)
			return nil
		}
		db.SetServers(repo.Servers)
	}

	syncDBs, err := (*handle).SyncDBs()
	if err != nil {
		log.Fatal(err)
	}
	return syncDBs
}
