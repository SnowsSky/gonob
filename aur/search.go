package aur

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AURResponse struct {
	ResultCount int         `json:"resultcount"`
	Results     []AURResult `json:"results"`
}

type AURResult struct {
	Name       string  `json:"Name"`
	Version    string  `json:"Version"`
	Maintainer *string `json:"Maintainer"`
	Popularity float64 `json:"Popularity"`
}

func Search(pkg string, to_install bool) (string, string, *string, float64, error) {
	if to_install {
		URL := "https://aur.archlinux.org/rpc.php?v=5&type=info&arg=" + pkg
		response, err := http.Get(URL)
		if err != nil {
			return "", "", nil, 0, fmt.Errorf("==> HTTP ERROR: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return "", "", nil, 0, fmt.Errorf("==> HTTP STATUS: %s", response.Status)
		}

		var aurResp AURResponse
		if err := json.NewDecoder(response.Body).Decode(&aurResp); err != nil {
			return "", "", nil, 0, fmt.Errorf("==> JSON ERROR: %v", err)
		}

		if aurResp.ResultCount == 0 || len(aurResp.Results) == 0 {
			return "", "", nil, 0, fmt.Errorf("package not found")
		}

		result := aurResp.Results[0]
		return result.Name, result.Version, result.Maintainer, result.Popularity, nil

	}
	return "", "", nil, 0, nil
}
