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

var lastPercent = make(map[string]int)

func Install(handle *alpm.Handle, syncDBs []alpm.Database, packages []string, noconfirm bool) {
	for _, pkg := range packages {
		pkgInfo, err := SearchOnSyncDatabases(pkg, handle, syncDBs)
		if pkgInfo == nil || err != nil {
			// pkg is not in DBS
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("unknown_package") + "\n")
			return
		}
		pkgInfos = append(pkgInfos, pkgInfo)

	}

	trans := alpm.NewTransaction(*handle)

	err := trans.Init(0)
	if err != nil {
		if CheckLock() {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("lock_file_found") + "\n")
		}
		return
	}
	(*handle).SetDownloadCallbackFunc(DownloadProgressCallback)
	(*handle).SetProgressCallbackFunc(InstallProgressCallback)
	defer trans.Release()
	StopHandle(trans)
	for _, pkg := range pkgInfos {
		err = trans.AddPkg(pkg)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	_, err = trans.Prepare()
	if err != nil {
		log.Fatal(err)
	}
	pkgs, err := trans.GetAdd()
	if err != nil {
		log.Fatal(err)
	}

	for i, pkg := range pkgs {
		pkgSizeMiB := float64(pkg.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkg.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		scolor.BoldBlue.DisplayText("(" + fmt.Sprintf("%d", i+1) + ") " + "--> ")
		scolor.BoldGreen.DisplayText(pkg.DB().Name() + ":")
		scolor.BoldWhite.DisplayText(pkg.Name() + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)\n")
	}
	scolor.BoldWhite.DisplayText("==> " + fmt.Sprint(len(pkgs)) + " " + translations.Translate("len_packages_to_add") + ".\n")
	scolor.BoldBlue.DisplayText("==> " + translations.Translate("size_to_add") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB\n")
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
		scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_read_alpm_log") + " [y/n] ")
		fmt.Scan(&response)
		if strings.ToLower(response) == "n" {
			return
		}
	}

	// Open the log file in the default editor and make the program wait until the editor is closed
	cmd := exec.Command("xdg-open", "/tmp/alpm.log")
	err = cmd.Run()
}
