package aur

import (
	"fmt"
	"gonob/translations"
	"log"
	"strings"

	alpm "github.com/Jguer/dyalpm"
)

type AurPackage struct {
	Name    string
	Version string
}

var response string

func DetectNonOfficialPackages(handle *alpm.Handle, syncDBs []alpm.Database) []AurPackage {
	AurPackages := []AurPackage{}

	localDB, err := (*handle).LocalDB()
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
	return AurPackages
}

func Update(handle *alpm.Handle, syncDBs []alpm.Database, noconfirm bool) {
	fmt.Println(Blue + "==> " + Reset + White + translations.Translate("fetch_aur_updates") + Reset)
	AurPackages := DetectNonOfficialPackages(handle, syncDBs)
	ToUpdate := []string{}

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
	if !noconfirm {
		fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [y/n] " + Reset)
		fmt.Scan(&response)
		if strings.ToLower(response) == "n" {
			fmt.Println(Red + "==> " + Reset + White + translations.Translate("canceled") + Reset)
			return
		} else {
			Install(ToUpdate, handle, true)
		}
	} else {
		Install(ToUpdate, handle, true)
	}

}
