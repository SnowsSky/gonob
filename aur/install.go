package aur

import (
	"fmt"
	"os"
	"os/exec"
)

var builddest string

func Install(pkgs []string) {
	for _, pkg := range pkgs {
		pkg_name, pkg_version, pkg_maintainer, pkg_popularity, err := Search(pkg, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		builddest = "/tmp/" + pkg_name
		fmt.Println(pkg_name, pkg_version, pkg_maintainer, pkg_popularity)
		cmd := exec.Command("git", "clone", fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg_name), "/tmp/"+pkg_name)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("==> ERROR: Failed to clone", pkg_name)
		}
	}
}
