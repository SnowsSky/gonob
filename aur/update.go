package aur

import (
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func Update() {
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Release()
	localDB, err := handle.LocalDB()
	if err != nil {
		return
	}
	AurPackages := []string{}
	for _, pkg := range localDB.PkgCache() {
		if pkg.Origin() != alpm.PkgFromSyncDB {
			AurPackages = append(AurPackages, pkg.Name())
		}
	}
}
