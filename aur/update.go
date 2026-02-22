package aur

import (
	"fmt"
	"gonob/translations"
	"log"
	"strings"

	alpm "github.com/Jguer/dyalpm"
	pacmanconf "github.com/Morganamilo/go-pacmanconf"
)

type AurPackage struct {
	Name    string
	Version string
}

var response string

func Update(handle *alpm.Handle) {
	AurPackages := []AurPackage{}
	ToUpdate := []string{}

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
			AurPackages = append(AurPackages, AurPackage{Name: pkg.Name(), Version: pkg.Version()})
		}
	}
	if len(AurPackages) == 0 {
		fmt.Println(Green + "==> " + translations.Translate("warning_string") + " : " + Reset + White + translations.Translate("no_aur_updates") + Reset)
		return
	}
	AurUpdates := 0
	for _, pkg := range AurPackages {
		_, aur_version, _, _, err := InstallSearch(pkg.Name)
		if err != nil {
			continue
		}
		if aur_version != pkg.Version {
			AurUpdates++
			ToUpdate = append(ToUpdate, pkg.Name)
			fmt.Println(Green + "==> " + Reset + White + pkg.Name + "@" + Reset + Yellow + pkg.Version + Reset + " --> " + Green + aur_version + Reset)
		}
	}
	if AurUpdates == 0 {
		fmt.Println(Green + "==> " + translations.Translate("warning_string") + " : " + Reset + White + translations.Translate("no_aur_updates") + Reset)
		return
	}
	fmt.Println(Yellow + "==> " + Reset + White + fmt.Sprint(AurUpdates) + " " + translations.Translate("aur_updates_available") + Reset)
	fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [Y/n] " + Reset)
	fmt.Scan(&response)
	if strings.ToLower(response) == "n" {
		fmt.Println(Red + "==> " + Reset + White + translations.Translate("canceled") + Reset)
		return
	}

	Install(ToUpdate, handle, true)

}
