package aur

import (
	"fmt"
	"gonob/translations"

	alpm "github.com/Jguer/dyalpm"
)

func List(handle *alpm.Handle) {
	AurPackages := DetectAURPackages(handle)
	UnknownPackages := []AurPackage{}
	FilteredPackages := []AurPackage{}
	for _, pkg := range AurPackages {
		_, _, _, _, err := InstallSearch(pkg.Name)
		if err != nil {
			UnknownPackages = append(UnknownPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		} else {
			FilteredPackages = append(FilteredPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		}
	}
	for _, pkg := range FilteredPackages {
		fmt.Println(Green + "==> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprint(len(FilteredPackages)) + " " + translations.Translate("aur_packages") + Reset)
	for _, pkg := range UnknownPackages {
		fmt.Println(Green + "\n==> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprint(len(UnknownPackages)) + " " + translations.Translate("unknown_package_source") + Reset)
}
