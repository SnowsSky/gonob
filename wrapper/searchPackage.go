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
