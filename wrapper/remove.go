package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
)

var pkgInfos []alpm.Package
var TotalSizeBytes float64 = 0
var TotalSizeMiB float64 = 0

func ProgressBarCallback(progress int32, pkg string, percent int, howmany uint64, current uint64) {
	barLen := 30
	filled := int(float64(percent) / 100.0 * float64(barLen))

	bar := ""
	if filled > 0 {
		bar += strings.Repeat("=", filled-1)
		bar += ">"
	}
	bar += strings.Repeat(" ", barLen-filled)

	fmt.Printf("\r(%d/%d) %s : %s [%s] %3d%% ", current, howmany, translations.Translate("removing"), pkg, bar, percent)
	if percent == 100 {
		fmt.Println()
	}
}

func Remove(handle *alpm.Handle, syncDBs []alpm.Database, packages []string, noconfirm bool) {
	for _, pkg := range packages {
		pkgInfo, err := SearchPackage(pkg, handle)
		if pkgInfo == nil || err != nil {
			// package is not installed.
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("package_not_installed") + Reset)
			return
		}
		if syncpkg, _ := SearchOnSyncDatabases(pkg, handle, syncDBs); syncpkg != nil {
			if syncpkg.DB().Name() == "core" {
				fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("can't_remove_core_packages") + Reset)
				return
			}
		}

		pkgInfos = append(pkgInfos, pkgInfo)

	}

	trans := alpm.NewTransaction(*handle)

	flags := alpm.TransFlagRecurse | alpm.TransFlagNoSave
	err := trans.Init(flags)
	if err != nil {
		if CheckLock() {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("lock_file_found") + Reset)
		}
		return
	}
	(*handle).SetProgressCallbackFunc(ProgressBarCallback)
	defer trans.Release()

	for _, pkg := range pkgInfos {
		err = trans.RemovePkg(pkg)
		if err != nil {
			log.Fatal(err)
			trans.Release()
			return
		}
	}

	_, err = trans.Prepare()
	if err != nil {
		log.Fatal(err)
		trans.Release()
		return
	}

	DepsToRemove, err := trans.GetRemove()

	for i, pkg := range DepsToRemove {
		pkgSizeMiB := float64(pkg.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkg.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		fmt.Println(Blue + "(" + fmt.Sprintf("%d", i+1) + ") " + "--> " + Reset + White + pkg.Name() + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)" + Reset)
	}
	fmt.Println(White + "==> " + fmt.Sprint(len(packages)) + " " + translations.Translate("len_packages_to_remove") + "." + Reset)
	fmt.Println(Blue + "==> " + translations.Translate("size_to_remove") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB")
	var response string
	if !noconfirm {

		fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [y/n] " + Reset)
		fmt.Scan(&response)
		if strings.ToLower(response) == "n" {
			fmt.Println(Red + "==> " + Reset + White + translations.Translate("canceled") + Reset)
			trans.Release()
			return
		}
	}
	// Commit the transaction
	conflicts, err := trans.Commit()
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(conflicts) > 0 {
		fmt.Println("File conflicts detected!")
		return
	}
	fmt.Println(Green + "==> " + Reset + White + translations.Translate("sucess") + Reset)
	if !noconfirm {
		fmt.Print(White + "==> " + translations.Translate("ask_to_read_alpm_log") + " [y/n] " + Reset)
		fmt.Scan(&response)
		if strings.ToLower(response) == "n" {
			return
		}
		// Open the log file in the default editor and make the program wait until the editor is closed
		cmd := exec.Command("xdg-open", "/tmp/alpm.log")
		err = cmd.Run()
	}

}
