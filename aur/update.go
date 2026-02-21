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
	syncDBs, err := handle.SyncDBs()
	if err != nil {
		log.Fatal(err)
	}

	// Collecte de tous les noms de paquets officiels
	officialPackages := make(map[string]bool)
	for _, db := range syncDBs {
		for _, pkg := range db.PkgCache().Collect() {
			officialPackages[pkg.Name()] = true
		}
	}

	// Vérification de chaque paquet local
	for _, pkg := range localDB.PkgCache().Collect() {
		if !officialPackages[pkg.Name()] {
			// Si le paquet n'est pas dans les repos officiels → c'est probablement un AUR
			AurPackages = append(AurPackages, pkg.Name())
		}
	}
	fmt.Println(AurPackages)

}
