package wrapper

import (
	"gonob/translations"
	"log"

	alpm "github.com/Jguer/dyalpm"
	"github.com/Morganamilo/go-pacmanconf"
	scolor "github.com/SnowsSky/scolor/pkg"
)

func InitSyncDatabases(handle *alpm.Handle) []alpm.Database {
	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string"))
		scolor.BoldWhite.DisplayText(" : " + translations.Translate("alpm_unable_to_fetch-_syncdbs") + "\n")
		return nil
	}

	for _, repo := range conf.Repos {
		db, err := (*handle).RegisterSyncDB(repo.Name, 0)
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string"))
			scolor.BoldWhite.DisplayText(" : " + translations.Translate("alpm_unable_to_fetch-_syncdbs") + "\n")
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
