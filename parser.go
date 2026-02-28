package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
	"gonob/wrapper"
	"os"
)

var version = "1.0.0-dev-16"

func parser(args []string) {
	if len(args) == 0 {
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
		return
	}
	handle := wrapper.InitHandle()
	syncDBs := wrapper.InitSyncDatabases(handle)
	defer (*handle).Release()

	for _, arg := range os.Args[0:] {
		if arg == "remove" || arg == "-R" {
			if os.Geteuid() != 0 {
				fmt.Println(aur.Red + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("need_sudo_privileges") + aur.Reset)
				return
			}
		}
	}

	switch args[0] {
	case "install", "-S":
		if args[1] == "--aur" {
			if len(args) > 2 {
				if args[2] == "--noconfirm" {
					aur.Install(args[3:], handle, true)
					return
				}
			}

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
			if len(args) > 2 {
				if args[2] == "--noconfirm" {
					aur.Update(handle, syncDBs, true)
					return
				}
			}
			aur.Update(handle, syncDBs, false)

		}
	case "list", "-Q":
		if args[1] == "--aur" {
			aur.List(handle, syncDBs)
		}
	case "remove", "-R":
		wrapper.Remove(handle, syncDBs, args[1:])
	case "--help", "-h":
		fmt.Println("Usage: gonob [command] [options]\n\nCommands:\n  install, -S      Install a package\n  remove, -R		Remove a package\n  search, -Ss      Search for a package\n  list, -Q         List installed packages\n  upgrade, -Syu    Upgrade all packages\n  --version, -v    Show version information\n  --help, -h       Show this help message\n\nOptions:\n  --aur            Assume that your query is from the AUR.\n  --noconfirm            Assume that the response of all confirmation messages are 'yes'.")
	default:
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
	}
}
