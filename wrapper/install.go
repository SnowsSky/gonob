package wrapper

import (
	"fmt"
	"gonob/translations"
	"log"
	"os/exec"
	"strings"

	"github.com/Jguer/dyalpm"
	alpm "github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
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
	if filled == barLen {
		bar = strings.Repeat("=", barLen)
	} else if filled > 0 {
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
		percent := float64(data.Downloaded) / float64(data.Total) * 100
		fmt.Printf("\r%s : %s (%.1f%%)", translations.Translate("downloading"), ev.Filename, percent)

	case dyalpm.DownloadCompleted:
		fmt.Printf("\r%s : %s (100%%)\n", translations.Translate("downloading"), ev.Filename)
	}
}
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
