package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
	"gonob/wrapper"
	"os"
)

var version = "1.1.0"

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

	if len(args) > 1 {
		if args[1] == "--aur" && os.Geteuid() == 0 {
			fmt.Println(aur.Red + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("don't_use_sudo") + aur.Reset)
			return
		}
	}

	switch args[0] {
	case "install", "-S":
		if len(args) > 1 {
			if args[1] == "--aur" {
				if len(args) > 2 {
					if args[2] == "--noconfirm" {
						aur.Install(args[3:], handle, true)
						return
					}
				}
				aur.Install(args[2:], handle, false)
				return
			}
			if os.Geteuid() != 0 {
				fmt.Println(aur.Red + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("need_sudo_privileges") + aur.Reset)
				return
			}
			wrapper.Install(handle, syncDBs, args[1:])
		}
	case "--version", "-v":
		fmt.Println(aur.White + "gonob@" + version + "\nhttps://github.com/SnowsSky/gonob" + aur.Reset)
	case "search", "-Ss":
		if len(args) > 1 {
			if args[1] == "--aur" {
				aur.Search(args[2])
				return
			}
			wrapper.Search(args[1], handle, syncDBs)
		}

	case "upgrade", "-Syu":
		if len(args) > 1 {
			if args[1] == "--aur" {
				if len(args) > 2 {
					if args[2] == "--noconfirm" {
						aur.Update(handle, syncDBs, true)
						return
					}
				}
				aur.Update(handle, syncDBs, false)
				return

			}
		}
		return
	case "list", "-Q":
		if len(args) > 1 {
			if args[1] == "--aur" {
				aur.List(handle, syncDBs)
				return
			}
		}
		wrapper.List(handle, syncDBs)

	case "remove", "-R":
		if len(args) > 1 {
			if args[1] == "--noconfirm" {
				wrapper.Remove(handle, syncDBs, args[2:], true)
				return
			}
			wrapper.Remove(handle, syncDBs, args[1:], false)
		}
	case "local_install", "-U":
		if len(args) > 1 {
			if args[1] == "--noconfirm" {
				wrapper.Local_Install(handle, args[2:], true)
				return
			}
		}
		wrapper.Local_Install(handle, args[1:], false)
	case "release_notes":
		Release_note()
	case "--help", "-h":
		fmt.Println("Usage: gonob [command] [options]\n\nCommands:\n  install, -S      Install a package\n  local_install, -U      Install a local package\n  remove, -R		Remove a package\n  search, -Ss      Search for a package\n  list, -Q         List installed packages\n  upgrade, -Syu    Upgrade all packages\n  release_notes    See the releases notes for gonob\n  --version, -v    Show version information\n  --help, -h       Show this help message\n\nOptions:\n  --aur            Assume that your query is from the AUR.\n  --noconfirm            Assume that the response of all confirmation messages are 'yes'.")
	default:
		fmt.Println(aur.Yellow + "==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + aur.Reset)
	}
}
