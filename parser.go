package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
)

var version = "1.0.0-dev-2"

func parser(args []string) {
	if len(args) == 0 {
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
		return
	}
	switch args[0] {
	case "install", "-S":
		if args[1] == "--aur" {
			aur.Install(args[2:])
		}
	case "--version", "-v":
		fmt.Println(aur.White + "gonob@" + version + "\nhttps://github.com/SnowsSky/gonob" + aur.Reset)
	case "search", "-Ss":
		aur.Search(args[1])
	default:
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
	}
}
