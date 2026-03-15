package aur

import (
	"fmt"
	"gonob/translations"
	"log"
	"strings"

	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
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
	scolor.BoldBlue.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(translations.Translate("fetch_aur_updates") + "\n")
	AurPackages := DetectNonOfficialPackages(handle, syncDBs)
	ToUpdate := []string{}

	if len(AurPackages) == 0 {
		scolor.BoldGreen.DisplayText("==> " + translations.Translate("warning_string") + " : ")
		scolor.BoldWhite.DisplayText(translations.Translate("no_aur_updates") + "\n")
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
			scolor.BoldGreen.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(pkg.Name + "@")
			scolor.BoldYellow.DisplayText(pkg.Version)
			scolor.BoldWhite.DisplayText(" --> ")
			scolor.BoldGreen.DisplayText(aur_version + "\n")
		}
	}
	if AurUpdates == 0 {
		scolor.BoldGreen.DisplayText("==> " + translations.Translate("warning_string") + " : ")
		scolor.BoldWhite.DisplayText(translations.Translate("no_aur_updates") + "\n")
		return
	}
	scolor.BoldYellow.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(fmt.Sprint(AurUpdates) + " " + translations.Translate("aur_updates_available") + "\n")
	if !noconfirm {
		scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_continue") + " [y/n] ")
		fmt.Scan(&response)
		if strings.ToLower(response) == "n" {
			scolor.BoldRed.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(translations.Translate("canceled") + "\n")
			return
		} else {
			Install(ToUpdate, handle, true)
		}
	} else {
		Install(ToUpdate, handle, true)
	}

}
