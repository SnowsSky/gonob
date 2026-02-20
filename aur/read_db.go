package aur

import (
	"fmt"
	"os"

	"github.com/Jguer/go-alpm/v2"
)

func Read_db() []string {
	h, er := alpm.Initialize("/", "/var/lib/pacman")
	if er != nil {
		print(er, "\n")
		os.Exit(1)
	}

	db, er := h.LocalDB()
	if er != nil {
		fmt.Println(er)
		os.Exit(1)
	}

	if h.Release() != nil {
		os.Exit(1)
	}

	pkgs := db.PkgCache().Slice()
	names := make([]string, len(pkgs))
	for i, p := range pkgs {
		names[i] = p.Name()
	}
	return names
}
