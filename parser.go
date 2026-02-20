package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
)

var version = "1.0.0-dev-1"

func parser(args []string) {
	if len(args) == 0 {
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
		return
	}
	switch args[0] {
	case "install":
		if args[1] == "--aur" {
			aur.Install(args[2:])
		}
	case "--version", "-v":
		fmt.Println(aur.White + "gonob@" + version + aur.Reset)
	default:
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
	}
}
