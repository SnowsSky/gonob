package aur

import (
	"fmt"
	"gonob/translations"
	"os"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
)

var Reset = "\033[0m"
var Red = "\033[1;31m"
var Green = "\033[1;32m"
var Yellow = "\033[1;33m"
var Blue = "\033[1;34m"
var Magenta = "\033[1;35m"
var Cyan = "\033[1;36m"
var Gray = "\033[1;37m"
var White = "\033[1;97m"
var builddest string

func CheckPkgFolder() bool {
	if _, err := os.Stat(builddest); os.IsNotExist(err) {
		// the folder does not exist.
		return false
	}
	return true
}

func Read_db(pkg_name string) error {
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		return err
	}
	defer handle.Release()

	// Get local database
	localDB, err := handle.LocalDB()
	if err != nil {
		return err
	}

	// Get a package
	pkg := localDB.Pkg(pkg_name)
	if pkg == nil {
		return err
	}
	return nil
}

func Install(pkgs []string) {
	for i, pkg := range pkgs {
		pkg_name, pkg_version, pkg_maintainer, pkg_popularity, err := InstallSearch(pkg)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = Read_db(pkg_name)
		if err != nil {
			fmt.Println(Green + "==> " + Reset + White + translations.Translate("installing") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(pkgs)) + "]\n  " + Blue + "-->" + Reset + " " + White + pkg_name + "@" + pkg_version + "..." + Reset)
		} else {
			fmt.Println(Green + "==> " + Reset + White + translations.Translate("reinstalling") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(pkgs)) + "]\n  " + Blue + "-->" + Reset + " " + White + pkg_name + "@" + pkg_version + "..." + Reset)

		}
		builddest = "/tmp/" + pkg_name
		if pkg_popularity <= 2.5 {
			var response string
			fmt.Println(Yellow + "==> " + translations.Translate("warning_string") + " : " + Reset + White + translations.Translate("low_popularity") + Reset)
			fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [Y/n] " + Reset)
			fmt.Scan(&response)
			if strings.ToLower(response) == "n" {
				fmt.Println(Red + "==> " + Reset + White + translations.Translate("canceled") + Reset)
				return
			}
		}
		fmt.Println(pkg_name, pkg_version, pkg_maintainer, pkg_popularity)

		if !CheckPkgFolder() {
			cmd := exec.Command("git", "clone", fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg_name), "/tmp/"+pkg_name)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + translations.Translate("clone_error") + Reset)
			}
		} else {
			fmt.Println(Yellow + "==> " + translations.Translate("warning_string") + " : " + Reset + White + translations.Translate("folder_already_exists") + Reset)
		}

		cmd := exec.Command("makepkg", "-si", "--noconfirm")
		cmd.Dir = builddest
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + translations.Translate("build_error") + Reset)
			return
		}
		fmt.Println(Green + "==> " + Reset + White + translations.Translate("build_success") + Reset)
	}
}
