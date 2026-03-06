package translations

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var data map[string]string

func Translate(translation_type string) string {
	locale := os.Getenv("LANG")
	locale = strings.Split(locale, ".")[0]
	if len(data) <= 0 {
		file, err := os.ReadFile("/etc/gonob/translations/" + locale + ".json")
		if err != nil {
			file, err = os.ReadFile("/etc/gonob/translations/us_US.json")
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
		file, err := os.ReadFile("/etc/gonob/translations/us_US.json")
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
