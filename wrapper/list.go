package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"

	alpm "github.com/Jguer/dyalpm"
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
		fmt.Println(Green + "==> " + Reset + Green + pkg_db + ":" + Reset + White + pkg.Name() + "@" + pkg.Version() + Reset)
		i++
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprintf("%d", i) + " " + translations.Translate("installed_packages") + Reset)
}
