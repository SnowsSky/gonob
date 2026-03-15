package aur

import (
	"fmt"
	"gonob/translations"
	"os"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

var dest = "/tmp/packages"

func CheckPackageList() bool {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		// the folder does not exist.
		return false
	}
	return true
}

func GetAurPackagesList() {
	if !CheckPackageList() {
		cmd := exec.Command("curl", "--retry", "3", "-s", "-o", "/tmp/packages.gz", "https://aur.archlinux.org/packages.gz")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
			scolor.BoldWhite.DisplayText(translations.Translate("get_aur_package_list_error"))
			fmt.Printf("%s\n", err)
		}
		cmd = exec.Command("gunzip", "-f", "/tmp/packages.gz")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
			scolor.BoldWhite.DisplayText(translations.Translate("get_aur_package_list_error"))
			fmt.Printf("%s\n", err)
		}

	}
}

func List(handle *alpm.Handle, syncDBs []alpm.Database) {
	Packages := DetectNonOfficialPackages(handle, syncDBs)
	UnknownPackages, AurPackages := FilterPackages(Packages)
	for _, pkg := range AurPackages {
		scolor.BoldGreen.DisplayText("--> ")
		scolor.BoldWhite.DisplayText(pkg.Name + "@" + pkg.Version + "\n")
	}
	scolor.BoldGreen.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(fmt.Sprint(len(AurPackages)) + " " + translations.Translate("aur_packages") + "\n")
	for _, pkg := range UnknownPackages {
		scolor.BoldYellow.DisplayText("--> ")
		scolor.BoldWhite.DisplayText(pkg.Name + "@" + pkg.Version + "\n")
	}
	if len(UnknownPackages) >= 1 {
		scolor.BoldYellow.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(fmt.Sprint(len(UnknownPackages)) + " " + translations.Translate("unknown_package_source") + "\n")
	}

}

func FilterPackages(pkgs []AurPackage) ([]AurPackage, []AurPackage) {
	UnknownPackages := []AurPackage{}
	AurPackages := []AurPackage{}
	GetAurPackagesList()
	data, err := os.ReadFile("/tmp/packages")
	content := string(data)
	if err != nil {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
		scolor.BoldWhite.DisplayText(translations.Translate("get_aur_package_list_error") + " ")
		fmt.Printf("%s\n", err)
	}
	for _, pkg := range pkgs {
		if strings.Contains(content, pkg.Name) {
			AurPackages = append(AurPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		} else {
			if !strings.Contains(pkg.Name, "-debug") {
				UnknownPackages = append(UnknownPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
			} else {
				AurPackages = append(AurPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
			}

		}
	}
	return UnknownPackages, AurPackages

}
