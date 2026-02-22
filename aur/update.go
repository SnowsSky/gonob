package aur

import (
	"fmt"
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func Update(handle *alpm.Handle) {
	AurPackages := []string{}

	localDB, err := (*handle).LocalDB()
	if err != nil {
		log.Fatal(err)
	}

	syncDBs, err := (*handle).SyncDBs()
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range localDB.PkgCache().Collect() {
		for _, db := range syncDBs {
			fmt.Printf("%s\n", db.Name())
			if db.Pkg(pkg.Name()) == nil {
				AurPackages = append(AurPackages, pkg.Name())
			}
		}
	}
	fmt.Println(AurPackages)

}
