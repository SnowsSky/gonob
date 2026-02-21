package aur

import (
	"fmt"
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func Update() {
	AurPackages := []string{}

	// Initialisation ALPM
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Release()

	// Base locale
	localDB, err := handle.LocalDB()
	if err != nil {
		log.Fatal(err)
	}

	// DB distantes
	/*syncDBs, err := handle.SyncDBs()
	if err != nil {
		log.Fatal(err)
	}*/

	for _, pkg := range localDB.PkgCache().Collect() {
		fmt.Println(pkg.Name(), pkg.Reason())
	}
	fmt.Println(AurPackages)

}
