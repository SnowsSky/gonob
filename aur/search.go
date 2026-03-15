package aur

import (
	"encoding/json"
	"fmt"
	"gonob/translations"
	"net/http"

	scolor "github.com/SnowsSky/scolor/pkg"
)

type AURResponse struct {
	ResultCount int         `json:"resultcount"`
	Results     []AURResult `json:"results"`
}

type AURResult struct {
	Name        string  `json:"Name"`
	Version     string  `json:"Version"`
	Maintainer  string  `json:"Maintainer"`
	Description string  `json:"Description"`
	Popularity  float64 `json:"Popularity"`
}

func InstallSearch(pkg string) (string, string, string, float64, error) {
	URL := "https://aur.archlinux.org/rpc.php?v=5&type=info&arg=" + pkg
	response, err := http.Get(URL)
	if err != nil {
		return "", "", "", 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", "", "", 0, fmt.Errorf("%s", response.StatusCode)
	}

	var aurResp AURResponse
	if err := json.NewDecoder(response.Body).Decode(&aurResp); err != nil {
		return "", "", "", 0, fmt.Errorf(err.Error())
	}
	if aurResp.ResultCount == 0 || len(aurResp.Results) == 0 {
		return "", "", "", 0, fmt.Errorf(translations.Translate("unknown_aur_package"))
	}

	result := aurResp.Results[0]
	return result.Name, result.Version, result.Maintainer, result.Popularity, nil
}

func Search(pkg string) {
	URL := "https://aur.archlinux.org/rpc/v5/search/" + pkg
	response, err := http.Get(URL)
	if err != nil {
		scolor.BoldRed.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable") + "\n")
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		scolor.BoldRed.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable") + "\n")
		return
	}
	var aurResp AURResponse
	if err := json.NewDecoder(response.Body).Decode(&aurResp); err != nil {
		scolor.BoldRed.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("error_string") + " : " + err.Error() + "\n")
		return
	}

	for _, result := range aurResp.Results {
		if result.Popularity <= 2.5 {
			scolor.BoldYellow.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(result.Name + "@" + result.Version + " [" + result.Maintainer + "] " + "[" + translations.Translate("low_popularity") + "]\n   --> " + result.Description + "\n")
		} else {
			scolor.BoldGreen.DisplayText("==> ")
			scolor.BoldWhite.DisplayText(result.Name + "@" + result.Version + " [" + result.Maintainer + "]" + "\n   --> " + result.Description + "\n")
		}

	}
	if aurResp.ResultCount <= 0 {
		scolor.BoldYellow.DisplayText("==> " + translations.Translate("warning_string") + " : ")
		scolor.BoldWhite.DisplayText(fmt.Sprint(aurResp.ResultCount) + " " + translations.Translate("search_found") + " : " + pkg + "\n")
		return
	}
	scolor.BoldGray.DisplayText("==> ")
	scolor.BoldWhite.DisplayText(fmt.Sprint(aurResp.ResultCount) + " " + translations.Translate("search_found") + " : " + pkg + "\n")
}
