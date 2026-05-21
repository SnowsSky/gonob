package translations

import (
	"encoding/json"
	"fmt"
	"os"
)

var data map[string]string
var locale string = ""

func SetLang(lang string) {
	switch lang {
	case "fr_FR", "fr_CA", "fr_BE", "fr_CH":
		locale = "fr_FR"
	case "en_US", "en_CA", "en_GB", "en_AU":
		locale = "en_US"
	case "es_ES", "es_MX", "es_AR":
		locale = "es_ES"
	default:
		locale = lang
	}
}

func Translate(translation_type string) string {
	if len(data) <= 0 {
		file, err := os.ReadFile("/etc/gonob/translations/" + locale + ".json")
		if err != nil {
			file, err = os.ReadFile("/etc/gonob/translations/en_US.json")
			if err != nil {
				fmt.Println("Translations files at /etc/gonob/translations have been corrupted / deleted.\nPlease consider reinstalling gonob.")
				os.Exit(1)
				return ""
			}
		}

		if err := json.Unmarshal(file, &data); err != nil {
			fmt.Println("Translations files at /etc/gonob/translations have been corrupted / deleted.\nPlease consider reinstalling gonob.")
			os.Exit(1)
			return ""
		}
	}

	value, ok := data[translation_type]

	if !ok {
		file, err := os.ReadFile("/etc/gonob/translations/en_US.json")
		if err != nil {
			fmt.Println("Translations files at /etc/gonob/translations have been corrupted / deleted.\nPlease consider reinstalling gonob.")
			os.Exit(1)
			return ""
		}

		if err := json.Unmarshal(file, &data); err != nil {
			fmt.Println("Translations files at /etc/gonob/translations have been corrupted / deleted.\nPlease consider reinstalling gonob.")
			os.Exit(1)
			return ""
		}
		value, ok := data[translation_type]
		if !ok {
			fmt.Println("Translations files at /etc/gonob/translations have been corrupted / deleted.\nPlease consider reinstalling gonob.")
			return ""
		}

		return value
	}

	return value
}
