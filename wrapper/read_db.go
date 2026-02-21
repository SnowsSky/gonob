package wrapper

import (
	"fmt"
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func CheckPackageAvailabilityOnSyncDatabases() {
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Release()

	syncDBs, err := handle.SyncDBs()
	if err != nil {
		return
	}
	for _, db := range syncDBs {
		for pkg := range db.Packages() {
			fmt.Println(pkg)
		}
	}

}
