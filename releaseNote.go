package main

import (
	"bufio"
	"fmt"
	"gonob/aur"
	"gonob/translations"
	"net/http"
	"strings"
)

func Release_note() {
	url := "https://raw.githubusercontent.com/SnowsSky/gonob/main/patchnotes.md"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(aur.Red + "==> " + aur.Reset + aur.White + translations.Translate("error_string") + " : " + translations.Translate("unable_to_get_releases_notes") + aur.Reset)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var result strings.Builder
	capture := false
	fmt.Println(aur.White + "- gonob " + version + aur.Reset)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "- gonob") {
			if strings.TrimSpace(strings.TrimPrefix(line, "- gonob")) == version {
				capture = true
				continue
			}
			if capture {
				break
			}
		}
		if capture {
			result.WriteString(line + "\n")
		}
	}
	content := result.String()
	fmt.Println(content)
	fmt.Println(aur.White + "https://raw.githubusercontent.com/SnowsSky/gonob/main/patchnotes.md" + aur.Reset)
}
