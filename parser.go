package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
	"gonob/wrapper"
)

var version = "1.0.0-dev-7"

func parser(args []string) {
	if len(args) == 0 {
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
		return
	}
	handle := wrapper.InitHandle()
	switch args[0] {
	case "install", "-S":
		if args[1] == "--aur" {
			aur.Install(args[2:], handle)
		}
	case "--version", "-v":
		fmt.Println(aur.White + "gonob@" + version + "\nhttps://github.com/SnowsSky/gonob" + aur.Reset)
	case "search", "-Ss":
		if args[1] == "--aur" {
		}
	case "upgrade", "-Syu":
		if args[1] == "--aur" {

			aur.Update(handle)
		}
	default:
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
	}
}
