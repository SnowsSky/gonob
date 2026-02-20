package aur

import (
	"fmt"
	"os"
	"os/exec"
)

var builddest string

func CheckPkgFolder() bool {
	if _, err := os.Stat(builddest); os.IsNotExist(err) {
		// the folder does not exist.
		return false
	}
	return true
}

func Install(pkgs []string) {
	for _, pkg := range pkgs {
		pkg_name, pkg_version, pkg_maintainer, pkg_popularity, err := Search(pkg, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		builddest = "/tmp/" + pkg_name
		fmt.Println(pkg_name, pkg_version, pkg_maintainer, pkg_popularity)

		if !CheckPkgFolder() {
			fmt.Println("==> Cloning", pkg_name, "'s repository...")
			cmd := exec.Command("git", "clone", fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg_name), "/tmp/"+pkg_name)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				fmt.Println("==> ERROR: Failed to clone", pkg_name)
			}
		} else {
			fmt.Println("==> Warning:", pkg_name+"'s folder already exists, skipping...")
		}

		fmt.Println("==> Building", pkg_name+"...")
		cmd := exec.Command("makepkg", "-si", "--noconfirm")
		cmd.Dir = builddest
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("==> ERROR: Failed to build", pkg_name)
			return
		}

		fmt.Println("==> Package", pkg_name, "successfully built and installed.")
	}
}
