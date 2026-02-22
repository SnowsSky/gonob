package aur

import (
	"fmt"
	"gonob/translations"

	alpm "github.com/Jguer/dyalpm"
)

func List(handle *alpm.Handle) {
	AurPackages := DetectAURPackages(handle)
	for _, pkg := range AurPackages {
		fmt.Println(Green + "==> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprint(len(AurPackages)) + " " + translations.Translate("aur_packages") + Reset)
}
