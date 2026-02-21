package aur

import (
	"encoding/json"
	"fmt"
	"gonob/translations"
	"net/http"
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
		return "", "", "", 0, fmt.Errorf(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable") + Reset)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", "", "", 0, fmt.Errorf(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable") + Reset)
	}

	var aurResp AURResponse
	if err := json.NewDecoder(response.Body).Decode(&aurResp); err != nil {
		return "", "", "", 0, fmt.Errorf(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + err.Error() + Reset)
	}
	if aurResp.ResultCount == 0 || len(aurResp.Results) == 0 {
		return "", "", "", 0, fmt.Errorf(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + translations.Translate("unknown_aur_package") + Reset)
	}

	result := aurResp.Results[0]
	return result.Name, result.Version, result.Maintainer, result.Popularity, nil
}

func Search(pkg string) {
	URL := "https://aur.archlinux.org/rpc/v5/search/" + pkg
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable") + Reset)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + translations.Translate("aur_unreachable") + Reset)
		return
	}
	var aurResp AURResponse
	if err := json.NewDecoder(response.Body).Decode(&aurResp); err != nil {
		fmt.Println(Red + "==> " + Reset + White + translations.Translate("error_string") + " : " + err.Error() + Reset)
		return
	}

	for _, result := range aurResp.Results {
		fmt.Println(Green + "==> " + Reset + White + result.Name + "@" + result.Version + " [" + result.Maintainer + "]\n   --> " + result.Description)
	}
	fmt.Println(Green + "==> " + Reset + White + fmt.Sprint(aurResp.ResultCount) + translations.Translate("search_found") + Reset)
}
