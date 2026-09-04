package wrapper

import (
	"fmt"
	"gonob/config"
	"gonob/translations"
	"log"
	"os"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

func Update() {
	cmd := exec.Command(config.PrivilegeEscalationTool, "pacman", "-Sy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		scolor.BoldRed.DisplayText(translations.Translate("error_string") + " ==> " + translations.Translate("sync_error") + "\n")
		return
	}
}

func Upgrade(handle *alpm.Handle, syncDBs []alpm.Database, noconfirm bool) {
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
	(*handle).SetProgressCallbackFunc(UpgradeProgressCallback)
	defer trans.Release()
	StopHandle(trans)
	err = trans.SyncSysupgrade(false)
	if err != nil {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " : ")
		scolor.BoldWhite.DisplayText(translations.Translate("unable_to_upgrade") + " " + err.Error() + "\n")
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
	// TO CHANGE TO CHECK IF THERE IS PKG TO UPGRADE, not only the size
	if TotalSizeBytes <= float64(0.0) {
		scolor.BoldBlue.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("no_packages_to_update") + "\n")
		return
	}
	scolor.BoldWhite.DisplayText("==> " + fmt.Sprint(len(pkgs)) + " " + translations.Translate("len_packages_to_upgrade") + ".\n")
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
		scolor.BoldRed.DisplayText("==> Commit échoué : " + err.Error() + "\n")
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
		}
	}

	// Open the log file in the default editor and make the program wait until the editor is closed
	cmd := exec.Command("xdg-open", "/tmp/alpm.log")
	err = cmd.Run()

}
