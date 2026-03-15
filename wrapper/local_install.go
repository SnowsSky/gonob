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

func Local_Install(handle *alpm.Handle, packages []string, noconfirm bool) {
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

	for _, pkg := range packages {
		toadd, err := (*handle).LoadPackage(pkg, true, 0)
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("failed_to_get_local_package") + "\n")
			return
		}
		err = trans.AddPkg(toadd)
		if err != nil {
			scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
			scolor.BoldWhite.DisplayText(translations.Translate("failed_to_get_local_package") + "\n")
			return
		}
	}

	_, err = trans.Prepare()
	if err != nil {
		log.Fatal(err)
	}

	pkgs, err := trans.GetAdd()
	if len(pkgs) <= 0 {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	for i, pkg := range pkgs {
		pkgSizeMiB := float64(pkg.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkg.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		scolor.BoldBlue.DisplayText("(" + fmt.Sprintf("%d", i+1) + ") " + "--> ")
		scolor.BoldGreen.DisplayText("local" + ":")
		scolor.BoldWhite.DisplayText(pkg.Name() + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)\n")
	}
	scolor.BoldWhite.DisplayText("==> " + fmt.Sprint(len(pkgs)) + " " + translations.Translate("len_packages_to_add") + ".\n")
	scolor.BoldBlue.DisplayText("==> " + translations.Translate("size_to_add") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB\n")
	var response string
	if !noconfirm {
		scolor.BoldWhite.DisplayText("==> " + translations.Translate("ask_to_continue") + " [y/n] ")
		fmt.Scan(&response)
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
