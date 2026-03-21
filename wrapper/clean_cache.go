package wrapper

import (
	"fmt"
	"gonob/translations"
	"math"
	"os"
	"strings"

	scolor "github.com/SnowsSky/scolor/pkg"
)

func Clean_cache() {
	result, err := os.ReadDir(CacheDIR)
	if err != nil {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string"))
		fmt.Println(err)
	}
	totalSize := float64(0)
	ratio := math.Pow(10, float64(2))

	for _, f := range result {

		fInfo, err := os.Stat(CacheDIR + "/" + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		fSize := float64(fInfo.Size()) / (1024 * 1024)
		fSize = math.Round(fSize*ratio) / ratio
		totalSize += float64(fSize)
		if !strings.HasSuffix(f.Name(), ".sig") {
			scolor.BoldBlue.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(translations.Translate("removing") + f.Name() + " [" + fmt.Sprintf("%v", fSize) + " MiB]" + "\n")
		}
		os.Remove(CacheDIR + "/" + f.Name())

	}
	scolor.BoldBlue.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(translations.Translate("size_to_remove") + " : " + fmt.Sprintf("%.2f", totalSize) + " MiB\n")

}
