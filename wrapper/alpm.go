package wrapper

import (
	"fmt"
	"gonob/translations"
	"strings"

	"github.com/Jguer/dyalpm"
)

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
