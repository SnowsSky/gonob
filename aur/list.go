package aur

import (
	"fmt"
	"gonob/translations"

	alpm "github.com/Jguer/dyalpm"
)

func List(handle *alpm.Handle) {
	Packages := DetectAURPackages(handle)
	UnknownPackages := []AurPackage{}
	AurPackages := []AurPackage{}
	for _, pkg := range Packages {
		_, _, _, _, err := InstallSearch(pkg.Name)
		if err != nil {
			UnknownPackages = append(UnknownPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		} else {
			AurPackages = append(AurPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		}
	}
	for _, pkg := range AurPackages {
		fmt.Println(Green + "--> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprint(len(AurPackages)) + " " + translations.Translate("aur_packages") + Reset)
	for _, pkg := range UnknownPackages {
		fmt.Println(Yellow + "--> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Yellow + "==> " + Reset + White + fmt.Sprint(len(UnknownPackages)) + " " + translations.Translate("unknown_package_source") + Reset)
}
