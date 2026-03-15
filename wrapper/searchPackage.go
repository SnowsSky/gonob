package wrapper

import (
	"gonob/translations"

	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

func SearchPackage(pkg_name string, handle *alpm.Handle) (alpm.Package, error) {
	// Get local database
	localDB, err := (*handle).LocalDB()
	if err != nil {
		return nil, err
	}

	// Get a package
	pkg := localDB.Pkg(pkg_name)
	if pkg == nil {
		return nil, err
	}
	return pkg, nil
}

func Search(pkg_name string, handle *alpm.Handle, syncDBs []alpm.Database) {
	pkg, _ := SearchOnSyncDatabases(pkg_name, handle, syncDBs)
	if pkg == nil {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
		scolor.BoldWhite.DisplayText(translations.Translate("unknown_package") + "\n")
		return
	}
	scolor.BoldGreen.DisplayText("==> " + pkg.DB().Name())
	scolor.BoldWhite.DisplayText(":" + pkg.Name() + "@" + pkg.Version() + "-" + pkg.Architecture() + "\n" + pkg.Description() + "\n")
}

var found bool

func SearchOnSyncDatabases(pkg_name string, handle *alpm.Handle, syncDBs []alpm.Database) (alpm.Package, error) {
	for _, db := range syncDBs {
		pkg := db.Pkg(pkg_name)
		if pkg != nil {
			found = true
			return pkg, nil
		}
	}
	return nil, nil

}
