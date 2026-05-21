package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

var pkgInfos []alpm.Package
var TotalSizeBytes float64 = 0
var TotalSizeMiB float64 = 0

func Remove(handle *alpm.Handle, syncDBs []alpm.Database, packages []string, noconfirm bool) {
	localDB, err := (*handle).LocalDB()
	conflict := false
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, pkg := range packages {
		pkgInfo, err := SearchPackage(pkg, handle)
		if pkgInfo == nil || err != nil {
			// package is not installed.
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("package_not_installed") + "\n")
			return
		}

		for _, p := range localDB.PkgCache().Collect() {
			for _, dep := range p.Depends() {
				if dep.Name == pkg {
					conflict = true
					scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
					scolor.BoldWhite.DisplayText(translations.Translate("conflict") + " : " + fmt.Sprintf("%s\n", p.Name()))
				}
			}
		}

		if conflict {
			return
		}

		pkgInfos = append(pkgInfos, pkgInfo)

	}

	trans := alpm.NewTransaction(*handle)

	flags := alpm.TransFlagRecurse | alpm.TransFlagNoSave
	err = trans.Init(flags)
	if err != nil {
		if CheckLock() {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("lock_file_found") + "\n")
		}
		return
	}
	(*handle).SetProgressCallbackFunc(ProgressBarCallback)
	defer trans.Release()
	StopHandle(trans)

	for _, pkg := range pkgInfos {
		err = trans.RemovePkg(pkg)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	_, err = trans.Prepare()
	if err != nil {
		log.Fatal(err)
		return
	}

	DepsToRemove, err := trans.GetRemove()

	for i, pkg := range DepsToRemove {
		pkgSizeMiB := float64(pkg.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkg.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		scolor.BoldBlue.DisplayText("(" + fmt.Sprintf("%d", i+1) + ") " + "--> ")
		scolor.BoldWhite.DisplayText(pkg.Name() + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)\n")
	}
	scolor.BoldWhite.DisplayText("==> " + fmt.Sprint(len(DepsToRemove)) + " " + translations.Translate("len_packages_to_remove") + ".\n")
	scolor.BoldBlue.DisplayText("==> " + translations.Translate("size_to_remove") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB\n")
	var response string
	if !noconfirm {
		scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_continue") + " [Y/n] ")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "n" {
			scolor.BoldRed.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(translations.Translate("canceled") + "\n")
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
	scolor.BoldGreen.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(translations.Translate("sucess") + "\n")
	if !noconfirm {
		scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_read_alpm_log") + " [Y/n] ")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "n" {
			return
		} else {
			// Open the log file in the default editor and make the program wait until the editor is closed
			cmd := exec.Command("xdg-open", "/tmp/alpm.log")
			err = cmd.Run()
		}
	}

}
