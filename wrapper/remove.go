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

func Remove(handle *alpm.Handle, syncDBs []alpm.Database, packages []string) {
	for i, pkg := range packages {
		pkgInfo, err := SearchPackage(pkg, handle)
		if pkgInfo == nil || err != nil {
			// package is not installed.
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("package_not_installed") + Reset)
			return
		}
		pkgInfos = append(pkgInfos, pkgInfo)
		pkgSizeMiB := float64(pkgInfo.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkgInfo.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		fmt.Println(Blue + "(" + fmt.Sprintf("%d", i+1) + ") " + "--> " + Reset + White + pkg + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)" + Reset)
	}
	fmt.Println(White + "==> " + fmt.Sprint(len(packages)) + " " + translations.Translate("len_packages_to_remove") + "." + Reset)

	trans := alpm.NewTransaction(*handle)

	err := trans.Init(0)
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

	DepsToRemove, err := trans.Prepare()
	if err != nil {
		log.Fatal(err)
		trans.Release()
	}

	for _, dep := range DepsToRemove {
		fmt.Println(Green+"==> "+translations.Translate("dep_to_remove")+" : "+dep.GetDepend().GetName(), dep.GetDepend().GetVersion())
		depPkg := dep.GetDepend()
		pkgInfo, err := SearchPackage(depPkg.GetName(), handle)
		if pkgInfo == nil || err != nil {
			// package is not installed.
			trans.Release()
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("package_not_installed") + Reset)
			return
		}
		depSizeMiB := float64(pkgInfo.ISize()) / (1024 * 1024)
		fmt.Printf("    - %s %s (%.2f MiB)\n", depPkg.GetName(), depPkg.GetVersion(), depSizeMiB)
		TotalSizeBytes += float64(pkgInfo.ISize())
	}

	fmt.Println(Blue + "==> " + translations.Translate("size_to_remove") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB")

	var response string
	fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [y/n] " + Reset)
	fmt.Scan(&response)
	if strings.ToLower(response) == "n" {
		fmt.Println(Red + "==> " + Reset + White + translations.Translate("canceled") + Reset)
		trans.Release()
		return
	}
	//fmt.Println(Green + "==> " + Reset + White + translations.Translate("removing") + " [" + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(packages)) + "]\n  " + Blue + "-->" + Reset + " " + White + pkg + "..." + Reset)

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
	fmt.Print(White + "==> " + translations.Translate("ask_to_read_alpm_log") + " [y/n] " + Reset)
	fmt.Scan(&response)
	if strings.ToLower(response) == "n" {
		return
	}
	// Open the log file in the default editor and make the program wait until the editor is closed
	cmd := exec.Command("xdg-open", "/tmp/alpm.log")
	err = cmd.Run()

}
