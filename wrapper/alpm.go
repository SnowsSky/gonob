package wrapper

import (
	"fmt"
	"gonob/translations"
	"strings"

	"github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

var CacheDIR string = "/var/cache/pacman/pkg/"

func ProgressBarCallback(progress int32, pkg string, percent int, howmany uint64, current uint64) {
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

	fmt.Printf("\r(%d/%d) %s : %s [%s] %3d%% ", current, howmany, translations.Translate("removing"), pkg, bar, percent)
	if percent == 100 {
		fmt.Println()
	}
}

var last_event uint = 0
var devMode bool = false
var printed bool = false

func EventCallback(event dyalpm.Event) {
	if last_event == uint(event.Type) {
		fmt.Printf("\r\033[2K")
	} else {
		if printed {
			fmt.Print("\n")
			printed = false
		}

	}
	switch event.Type {
	// SCRIPTLET EVENT.
	case 17:
		scolor.BoldBlue.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("scriptlet"))
		printed = true

	default:
		if devMode {
			scolor.BoldYellow.DisplayText("==> ", translations.Translate("warning_string")+" : ")
			scolor.BoldWhite.DisplayText(translations.Translate("unknown_event") + fmt.Sprintf(": %d", event.Type) + "\n")
			printed = true
		}
	}
	last_event = uint(event.Type)
}

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

func UpgradeProgressCallback(progress int32, pkg string, percent int, howmany uint64, current uint64) {
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

	fmt.Printf("\r(%d/%d) %s : %s [%s] %3d%% ", current, howmany, translations.Translate("updating"), pkg, bar, percent)
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
