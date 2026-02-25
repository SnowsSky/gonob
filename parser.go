package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
	"gonob/wrapper"
)

var version = "1.0.0-dev-14"

func parser(args []string) {
	if len(args) == 0 {
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
		return
	}
	handle := wrapper.InitHandle()
	defer (*handle).Release()
	switch args[0] {
	case "install", "-S":
		if args[1] == "--aur" {
			aur.Install(args[2:], handle, false)
		}
	case "--version", "-v":
		fmt.Println(aur.White + "gonob@" + version + "\nhttps://github.com/SnowsSky/gonob" + aur.Reset)
	case "search", "-Ss":
		if args[1] == "--aur" {
			aur.Search(args[2])
		}
	case "upgrade", "-Syu":
		if args[1] == "--aur" {

			aur.Update(handle)
		}
	case "list", "-Q":
		if args[1] == "--aur" {
			aur.List(handle)
		}
	case "--help", "-h":
		fmt.Println("Usage: gonob [command] [options]\n\nCommands:\n  install, -S      Install a package\n  search, -Ss      Search for a package\n  list, -Q         List installed packages\n  upgrade, -Syu    Upgrade all packages\n  --version, -v    Show version information\n  --help, -h       Show this help message\n\nOptions:\n  --aur            Assume that your query is from the AUR.")
	default:
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
	}
}
