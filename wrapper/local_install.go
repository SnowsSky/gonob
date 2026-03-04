package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"
	"os/exec"
	"strings"

	alpm "github.com/Jguer/dyalpm"
)

func Local_Install(handle *alpm.Handle, packages []string, noconfirm bool) {
	trans := alpm.NewTransaction(*handle)
	err := trans.Init(0)
	if err != nil {
		if CheckLock() {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("lock_file_found") + Reset)
		}
		return
	}
	(*handle).SetDownloadCallbackFunc(DownloadProgressCallback)
	(*handle).SetProgressCallbackFunc(InstallProgressCallback)
	defer trans.Release()

	for _, pkg := range packages {
		toadd, err := (*handle).LoadPackage(pkg, false, 0)
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("failed_to_get_local_package") + Reset)
			return
		}
		err = trans.AddPkg(toadd)
		if err != nil {
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("failed_to_get_local_package") + Reset)
			return
		}
	}

	_, err = trans.Prepare()
	if err != nil {
		log.Fatal(err)
		trans.Release()
	}

	pkgs, err := trans.GetAdd()
	if len(pkgs) <= 0 {
		trans.Release()
		return
	}
	if err != nil {
		log.Fatal(err)
		trans.Release()
	}

	for i, pkg := range pkgs {
		pkgSizeMiB := float64(pkg.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkg.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		fmt.Println(Blue + "(" + fmt.Sprintf("%d", i+1) + ") " + "--> " + Reset + Green + "local" + ":" + Reset + White + pkg.Name() + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)" + Reset)
	}
	fmt.Println(White + "==> " + fmt.Sprint(len(pkgs)) + " " + translations.Translate("len_packages_to_add") + "." + Reset)
	fmt.Println(Blue + "==> " + translations.Translate("size_to_add") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB")
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
	}

	fmt.Scan(&response)
	if strings.ToLower(response) == "n" {
		return
	}
	// Open the log file in the default editor and make the program wait until the editor is closed
	cmd := exec.Command("xdg-open", "/tmp/alpm.log")
	err = cmd.Run()

}
