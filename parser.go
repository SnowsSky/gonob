package main

import (
	"fmt"
	"gonob/aur"
	"gonob/translations"
	"gonob/wrapper"
	"os"

	"github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

var version = "2.2.0"

func useSudo() {
	if os.Geteuid() != 0 {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("warning_string") + " : " + translations.Translate("need_sudo_privileges") + "\n")
		os.Exit(1)
	}
}

func dontUseSudo() {
	if os.Geteuid() == 0 {
		scolor.BoldRed.DisplayText("==> " + translations.Translate("warning_string") + " : " + translations.Translate("don't_use_sudo") + "\n")
		os.Exit(1)
	}

}

func parser(args []string) {
	handle := wrapper.InitHandle()
	syncDBs := wrapper.InitSyncDatabases(handle)
	defer (*handle).Release()
	noconfirm := false
	fromaur := false
	toexecute := ""
	entry := []string{""}
	for _, arg := range args[0:] {
		switch arg {
		case "--noconfirm":
			noconfirm = true
		case "--aur":
			fromaur = true
		case "install", "-S":
			toexecute = "install"
		case "search", "-Ss":
			toexecute = "search"
		case "upgrade", "-Syu":
			toexecute = "upgrade"
		case "list", "-Q":
			toexecute = "list"
		case "-v", "--version":
			toexecute = "version"
		case "remove", "-R":
			toexecute = "remove"
		case "local_install", "-U":
			toexecute = "local_install"
		case "clean_cache", "-Sc":
			toexecute = "clean_cache"
		case "release_notes":
			toexecute = "release_notes"
		case "--help", "-h":
			toexecute = "help"
		default:
			entry = []string{arg}
		}
	}

	switch toexecute {
	case "install":
		if fromaur {
			dontUseSudo()
			aur.Install(entry, handle, noconfirm)
		} else {
			useSudo()
			wrapper.Install(handle, syncDBs, entry, noconfirm)
		}
	case "search":
		if fromaur {
			aur.Search(entry[0])
		} else {
			wrapper.Search(entry[0], handle, syncDBs)
		}
	case "upgrade":
		if fromaur {
			aur.Update(handle, syncDBs, noconfirm)
		}
	case "list":
		if fromaur {
			aur.List(handle, syncDBs)
		} else {
			wrapper.List(handle, syncDBs)
		}
	case "remove":
		useSudo()
		wrapper.Remove(handle, syncDBs, entry, noconfirm)
	case "local_install":
		useSudo()
		wrapper.Local_Install(handle, entry, noconfirm)
	case "clean_cache":
		useSudo()
		wrapper.Clean_cache()
	case "release_notes":
		Release_note()
	case "help":
		fmt.Println(`Usage: gonob [command] [options]

Commands:
  install, -S         Install a package
  local_install, -U   Install a local package
  remove, -R          Remove a package
  search, -Ss         Search for a package
  list, -Q            List installed packages
  upgrade, -Syu       Upgrade all packages
  clean_cache, -Sc 	  Clean the pacman cache
  release_notes       See the release notes for gonob
  --version, -v       Show version information
  --help, -h          Show this help message

Options:
  --aur               Assume that your query is from the AUR
  --noconfirm         Assume 'yes' for all confirmation prompts
`)
	case "version":
		scolor.BoldWhite.DisplayText("gonob@" + version + "\nlibalpm@" + fmt.Sprintf("%s", dyalpm.Version()) + "\nhttps://github.com/SnowsSky/gonob" + "\n")
	default:
		scolor.BoldYellow.DisplayText("==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + "\n")
	}
}
