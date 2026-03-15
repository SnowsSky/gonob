package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"

	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

func List(handle *alpm.Handle, syncDBs []alpm.Database) {
	localDB, err := (*handle).LocalDB()
	if err != nil {
		log.Fatal(err)

	}
	var i int = 0
	for _, pkg := range localDB.PkgCache().Collect() {
		var pkg_db string = "aur"
		for _, db := range syncDBs {
			if db.Pkg(pkg.Name()) != nil {
				pkg_db = db.Name()
				break
			}
		}
		scolor.BoldGreen.DisplayText("==> " + pkg_db + ":")
		scolor.BoldWhite.DisplayText(pkg.Name() + "@" + pkg.Version() + "\n")
		i++
	}
	scolor.BoldGreen.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(fmt.Sprintf("%d", i) + " " + translations.Translate("installed_packages") + "\n")
}
