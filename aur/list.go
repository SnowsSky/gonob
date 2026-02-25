package aur

import (
	"fmt"
	"gonob/translations"
	"os"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
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
			fmt.Println(Red+"==> "+translations.Translate("error_string")+Reset+White+translations.Translate("get_aur_package_list_error")+Reset, err)
		}
		cmd = exec.Command("gunzip", "-f", "/tmp/packages.gz")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(Red+"==> "+translations.Translate("error_string")+Reset+White+translations.Translate("get_aur_package_list_error")+Reset, err)
		}

	}
}

func List(handle *alpm.Handle) {
	Packages := DetectNonOfficialPackages(handle)
	UnknownPackages, AurPackages := FilterPackages(Packages)
	for _, pkg := range AurPackages {
		fmt.Println(Green + "--> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprint(len(AurPackages)) + " " + translations.Translate("aur_packages") + Reset)
	for _, pkg := range UnknownPackages {
		fmt.Println(Yellow + "--> " + Reset + White + pkg.Name + "@" + pkg.Version + Reset)
	}
	fmt.Println(Yellow + "==> " + Reset + White + fmt.Sprint(len(UnknownPackages)) + " " + translations.Translate("unknown_package_source") + Reset)
}

func FilterPackages(pkgs []AurPackage) ([]AurPackage, []AurPackage) {
	UnknownPackages := []AurPackage{}
	AurPackages := []AurPackage{}
	GetAurPackagesList()
	data, err := os.ReadFile("/tmp/packages")
	content := string(data)
	if err != nil {
		fmt.Println(Red+"==> "+translations.Translate("error_string")+Reset+White+translations.Translate("get_aur_package_list_error")+Reset, err)
	}
	for _, pkg := range pkgs {
		if strings.Contains(content, pkg.Name) {
			AurPackages = append(AurPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		} else {
			UnknownPackages = append(UnknownPackages, AurPackage{Name: pkg.Name, Version: pkg.Version})
		}
	}
	return UnknownPackages, AurPackages

}
