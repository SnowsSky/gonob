package aur

import (
	"fmt"
	"log"

	alpm "github.com/Jguer/dyalpm"
	pacmanconf "github.com/Morganamilo/go-pacmanconf"
)

func Update(handle *alpm.Handle) {
	AurPackages := []string{}

	localDB, err := (*handle).LocalDB()
	if err != nil {
		log.Fatal(err)
	}

	conf, _, err := pacmanconf.ParseFile("/etc/pacman.conf")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, repo := range conf.Repos {
		db, err := (*handle).RegisterSyncDB(repo.Name, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.SetServers(repo.Servers)
	}

	syncDBs, err := (*handle).SyncDBs()
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range localDB.PkgCache().Collect() {
		found := false

		for _, db := range syncDBs {
			if db.Pkg(pkg.Name()) != nil {
				found = true
				break
			}
		}

		if !found {
			AurPackages = append(AurPackages, pkg.Name())
		}
	}
	fmt.Println(AurPackages)

}
