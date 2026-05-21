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
	scolor "github.com/SnowsSky/scolor/pkg"
	git "github.com/go-git/go-git/v6"
)

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
			scolor.BoldRed.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable"))
			fmt.Println(err)
			return
		}
		_, err = wrapper.SearchPackage(pkg_name, handle)
		if err != nil {
			scolor.BoldGreen.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(translations.Translate("reinstalling") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(pkgs)) + "]\n  ")
			scolor.BoldBlue.DisplayText("--> ")
			scolor.BoldWhite.DisplayText(pkg_name + "@" + pkg_version + "...\n")

		} else {
			scolor.BoldGreen.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(translations.Translate("installing") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(pkgs)) + "]\n  ")
			scolor.BoldBlue.DisplayText("--> ")
			scolor.BoldWhite.DisplayText(pkg_name + "@" + pkg_version + "...\n")
		}

		builddest = "/tmp/" + pkg_name
		if !noconfirm && pkg_popularity <= 2.5 {
			var response string
			scolor.BoldYellow.DisplayText("==> " + translations.Translate("warning_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("low_popularity") + "\n")
			scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_continue") + " [Y/n] ")
			fmt.Scanln(&response)
			if strings.ToLower(response) == "n" {
				scolor.BoldRed.DisplayText("==> ")
				scolor.BoldWhite.DisplayText(translations.Translate("canceled") + "\n")
				return
			}
		}
		fmt.Println(pkg_name, pkg_version, pkg_maintainer, pkg_popularity)

		if !CheckPkgFolder() {
			// Clone the given repository to the given directory
			_, err := git.PlainClone(builddest, &git.CloneOptions{
				URL:      fmt.Sprintf("https://aur.archlinux.org/%s.git", pkg_name),
				Progress: os.Stdout,
			})
			if err != nil {
				scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
				scolor.BoldWhite.DisplayText(translations.Translate("clone_error") + fmt.Sprintf("\n %s", err) + "\n")
			}
		} else {
			scolor.BoldYellow.DisplayText("==> " + translations.Translate("warning_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("folder_already_exists") + "\n")
		}
		if !noconfirm {
			scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_read_pkgbuild") + " [Y/n] ")
			fmt.Scanln(&response)
			if strings.ToLower(response) != "n" {
				// Open the PKGBUILD file in the default editor and make the program wait until the editor is closed
				cmd := exec.Command("xdg-open", builddest+"/PKGBUILD")
				err = cmd.Run()
				scolor.BoldWhite.DisplayText("==> " + translations.Translate("press_any_key_to_continue") + " : ")
				fmt.Scanln(&response)
			}
		}

		cmd := exec.Command("makepkg", "-s", "-f", "--noconfirm")
		cmd.Dir = builddest
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
			scolor.BoldWhite.DisplayText(translations.Translate("build_error") + "\n")
			return
		}
		scolor.BoldGreen.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("build_success") + "\n")

		var pkgPath string
		files, err := os.ReadDir(builddest)
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
			scolor.BoldWhite.DisplayText(translations.Translate("build_error") + "\n")
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
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
			scolor.BoldWhite.DisplayText(" " + translations.Translate("build_error") + "\n")
			return
		}

	}

}
