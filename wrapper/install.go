package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"
	"os/exec"
	"strings"

	"github.com/Jguer/dyalpm"
	alpm "github.com/Jguer/dyalpm"
)

var lastPercent = make(map[string]int)

func InstallProgressCallback(progress int32, pkg string, percent int, howmany uint64, current uint64) {
	if pkg == "" {
		return
	}

	if last, ok := lastPercent[pkg]; ok && last == percent {
		return
	}
	lastPercent[pkg] = percent

	barLen := 30
	filled := int(float64(percent) / 100.0 * float64(barLen))

	bar := ""
	if filled > 0 {
		bar += strings.Repeat("=", filled-1)
		bar += ">"
	}
	bar += strings.Repeat(" ", barLen-filled)

	fmt.Printf("\r(%d/%d) %s : %s [%s] %3d%% ", current, howmany, translations.Translate("installing"), pkg, bar, percent)
	if percent == 100 {
		fmt.Println()
	}
}

func DownloadProgressCallback(ev dyalpm.DownloadEvent) {
	switch ev.Type {
	case dyalpm.DownloadProgress:
		data := ev.Data.(dyalpm.DownloadProgressData)
		fmt.Printf("\r%s : %s (%dB/%dB)", string(translations.Translate("downloading")), ev.Filename, data.Downloaded, data.Total)
	case dyalpm.DownloadCompleted:
		fmt.Printf("\r%s : %s (100%%)", translations.Translate("downloading"), ev.Filename)
	}
}
func Install(handle *alpm.Handle, syncDBs []alpm.Database, packages []string) {
	for _, pkg := range packages {
		pkgInfo, err := SearchOnSyncDatabases(pkg, handle, syncDBs)
		if pkgInfo == nil || err != nil {
			// pkg is not in DBS
			fmt.Println(Red + "==> " + translations.Translate("error_string") + " : " + Reset + White + translations.Translate("unknown_package") + Reset)
			return
		}
		pkgInfos = append(pkgInfos, pkgInfo)

	}

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

	for _, pkg := range pkgInfos {
		err = trans.AddPkg(pkg)
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
	}
	pkgs, err := trans.GetAdd()
	if err != nil {
		log.Fatal(err)
		trans.Release()
	}

	for i, pkg := range pkgs {
		pkgSizeMiB := float64(pkg.ISize()) / (1024 * 1024)
		TotalSizeBytes += float64(pkg.ISize())
		TotalSizeMiB = float64(TotalSizeBytes) / (1024 * 1024)
		fmt.Println(Blue + "(" + fmt.Sprintf("%d", i+1) + ") " + "--> " + Reset + Green + pkg.DB().Name() + ":" + Reset + White + pkg.Name() + " (" + fmt.Sprintf("%.2f", pkgSizeMiB) + " MiB)" + Reset)
	}
	fmt.Println(White + "==> " + fmt.Sprint(len(pkgs)) + " " + translations.Translate("len_packages_to_add") + "." + Reset)
	fmt.Println(Blue + "==> " + translations.Translate("size_to_add") + " : " + fmt.Sprintf("%.2f", TotalSizeMiB) + "MiB")
	var response string
	fmt.Print(White + "==> " + translations.Translate("ask_to_continue") + " [y/n] " + Reset)
	fmt.Scan(&response)
	if strings.ToLower(response) == "n" {
		fmt.Println(Red + "==> " + Reset + White + translations.Translate("canceled") + Reset)
		trans.Release()
		return
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
	fmt.Print(White + "==> " + translations.Translate("ask_to_read_alpm_log") + " [y/n] " + Reset)
	fmt.Scan(&response)
	if strings.ToLower(response) == "n" {
		return
	}
	// Open the log file in the default editor and make the program wait until the editor is closed
	cmd := exec.Command("xdg-open", "/tmp/alpm.log")
	err = cmd.Run()
}
