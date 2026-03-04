package aur

import (
	"fmt"
	"gonob/translations"
	"gonob/wrapper"
	"os"
	"os/exec"
	"path/filepath"
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

func Install(pkgs []string, handle *alpm.Handle, noconfirm bool) {
	for i, pkg := range pkgs {
		pkg_name, pkg_version, pkg_maintainer, pkg_popularity, err := InstallSearch(pkg)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = wrapper.SearchPackage(pkg_name, handle)
		if err != nil {
			fmt.Println(Green + "==> " + Reset + White + translations.Translate("reinstalling") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(pkgs)) + "]\n  " + Blue + "-->" + Reset + " " + White + pkg_name + "@" + pkg_version + "..." + Reset)
		} else {
			fmt.Println(Green + "==> " + Reset + White + translations.Translate("installing") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(pkgs)) + "]\n  " + Blue + "-->" + Reset + " " + White + pkg_name + "@" + pkg_version + "..." + Reset)

		}

		builddest = "/tmp/" + pkg_name
		if !noconfirm && pkg_popularity <= 2.5 {
			var response string
			fmt.Println(Yellow + "==> " + translations.Translate("warning_string") + " : " + Reset + White + translations.Translate("low_popularity") + Reset)
			fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [y/n] " + Reset)
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
				fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("clone_error") + Reset)
			}
		} else {
			fmt.Println(Yellow + "==> " + translations.Translate("warning_string") + " : " + Reset + White + translations.Translate("folder_already_exists") + Reset)
		}
		if !noconfirm {
			fmt.Print(White + "==> " + translations.Translate("ask_to_read_pkgbuild") + " [y/n] " + Reset)
			fmt.Scan(&response)
			if strings.ToLower(response) != "n" {
				// Open the PKGBUILD file in the default editor and make the program wait until the editor is closed
				cmd := exec.Command("xdg-open", builddest+"/PKGBUILD")
				err = cmd.Run()
				fmt.Print(White + "==> " + translations.Translate("press_any_key_to_continue") + " : " + Reset)
				fmt.Scan(&response)
			}
		}

		cmd := exec.Command("makepkg", "-s", "-f", "--noconfirm")
		cmd.Dir = builddest
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + " " + translations.Translate("build_error") + Reset)
			return
		}
		fmt.Println(Green + "==> " + Reset + White + translations.Translate("build_success") + Reset)

		var pkgPath string
		files, err := os.ReadDir(builddest)
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + " " + translations.Translate("build_error") + Reset)
			return
		}
		for _, f := range files {
			name := f.Name()

			if strings.HasPrefix(name, pkg_name+"-") &&
				strings.HasSuffix(name, ".pkg.tar.zst") &&
				!strings.Contains(name, "-debug-") {
				pkgPath = filepath.Join(builddest, name)
				break
			}
		}
		cmd = exec.Command("sudo", "gonob", "-U", "--noconfirm", pkgPath)
		cmd.Dir = builddest
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + Reset + White + " " + translations.Translate("build_error") + Reset)
			return
		}

	}

}
