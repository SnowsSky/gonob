package main

import (
	"bufio"
	"fmt"
	"gonob/translations"
	"net/http"
	"strings"

	scolor "github.com/SnowsSky/scolor/pkg"
)

func Release_note() {
	url := "https://raw.githubusercontent.com/SnowsSky/gonob/main/patchnotes.md"

	resp, err := http.Get(url)
	if err != nil {
		scolor.BoldRed.DisplayText("==> ")
		scolor.BoldWhite.DisplayText(translations.Translate("error_string") + " : " + translations.Translate("unable_to_get_releases_notes") + "\n")
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var result strings.Builder
	capture := false
	scolor.BoldWhite.DisplayText("- gonob " + version + "\n")
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
	scolor.BoldWhite.DisplayText("https://raw.githubusercontent.com/SnowsSky/gonob/main/patchnotes.md \n")
}
