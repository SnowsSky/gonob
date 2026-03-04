package wrapper

import alpm "github.com/Jguer/dyalpm"

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
