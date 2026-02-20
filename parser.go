package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
)

func parser(args []string) {
	switch args[0] {
	case "install":
		if args[1] == "--aur" {
			aur.Install(args[2:])
		}
	default:
		fmt.Println(translations.Translate("unknown_command"))
	}
}
