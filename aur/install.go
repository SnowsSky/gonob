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
		fmt.Println("==> Cloning", pkg_name, "'s repository...")
		cmd := exec.Command("git", "clone", fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg_name), "/tmp/"+pkg_name)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("==> ERROR: Failed to clone", pkg_name)
		}

		fmt.Println("==> Building", pkg_name+"...")
		cmd = exec.Command("makepkg", "-si", "--noconfirm")
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
