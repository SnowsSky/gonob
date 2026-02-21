package aur

import (
	"fmt"
	"gonob/translations"
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func Update() {
	AurPackages := []string{}
	Packages := make(map[string]struct{})
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Release()
	localDB, err := handle.LocalDB()
	if err != nil {
		return
	}
	syncDBs, err := handle.SyncDBs()
	if err != nil {
		return
	}

	for _, db := range syncDBs {
		_ = db.PkgCache().ForEach(func(pkg alpm.Package) error {
			Packages[pkg.Name()] = struct{}{}
			return nil
		})
	}
	err = localDB.PkgCache().ForEach(func(pkg alpm.Package) error {
		if _, exists := Packages[pkg.Name()]; !exists {
			AurPackages = append(AurPackages, pkg.Name(), pkg.Version())
			fmt.Println("AUR:", pkg.Name())
		}
		return nil
	})
	if err != nil {
		return
	}
	if len(AurPackages) == 0 {
		fmt.Println(Green + "==> " + Reset + White + translations.Translate("no_aur_updates") + Reset)
	}
	for AurPackage := range AurPackages {
		AurPackage = AurPackage
	}

}
