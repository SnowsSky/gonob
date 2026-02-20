package translations

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Translate(translation_type string) string {
	locale := os.Getenv("LANG")
	locale = strings.Split(locale, ".")[0]

	file, err := os.ReadFile("/etc/gonob/translations/" + locale + ".json")
	if err != nil {
		return fmt.Sprintf("Missing translation file: %s", locale)
	}

	var data map[string]string
	if err := json.Unmarshal(file, &data); err != nil {
		return "Invalid translation file"
	}

	value, ok := data[translation_type]

	if !ok {
		return fmt.Sprintf("Missing translation key: %s", translation_type)
	}

	return value
}
