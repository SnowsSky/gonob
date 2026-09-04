package main

import (
	"fmt"
	"gonob/aur"
	"gonob/config"
	"gonob/translations"
	"gonob/wrapper"
	"os"
	"os/exec"
	"strings"

	"github.com/Jguer/dyalpm"
	scolor "github.com/SnowsSky/scolor/pkg"
)

var version = "2.5.0"

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
	lang := ""
	entry := []string{""}
	for i, arg := range args[0:] {
		switch arg {
		case "--noconfirm":
			noconfirm = true
		case "--lang":
			if len(args) <= 1 {
			} else {
				lang = args[i+1]
			}
		case "--aur":
			fromaur = true
		case "install", "-S":
			toexecute = "install"
		case "search", "-Ss":
			toexecute = "search"
		case "update", "-Sy":
			toexecute = "update"
		case "upgrade", "-Su":
			toexecute = "upgrade"
		case "-Syu":
			toexecute = "upgrade+update"
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
		case "check_cache", "-Sch":
			toexecute = "check_cache"
		case "release_notes":
			toexecute = "release_notes"
		case "--help", "-h", "help":
			toexecute = "help"
		default:
			entry = []string{arg}
		}
	}
	if lang != "" {
		translations.SetLang(lang)
	} else {
		l := os.Getenv("LANG")
		l = strings.Split(l, ".")[0]
		translations.SetLang(l)
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
	case "update":
		wrapper.Update()
	case "upgrade":
		if fromaur {
			aur.Update(handle, syncDBs, noconfirm)
		} else {
			useSudo()
			wrapper.Upgrade(handle, syncDBs, noconfirm)
		}
	case "upgrade+update":
		dontUseSudo()
		wrapper.Update()
		cmd := exec.Command(config.PrivilegeEscalationTool, "gonob", "-Su")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		aur.Update(handle, syncDBs, noconfirm)
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
	case "check_cache":
		wrapper.Check_cache()
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
  update -Sy 		  Sync databases [CMD NOT AVAILABLE FOR --aur, use gonob -Su --aur or -Syu]
  upgrade, -Su        Upgrade all packages
  -Syu 				  Sync databases & Upgrade all packages
  clean_cache, -Sc 	  Clean the pacman cache
  check_cache, -Sch   Check the pacman cache size
  release_notes       See the release notes for gonob
  --version, -v       Show version information
  --help, -h          Show this help message

Options:
  --aur               Assume that your query is from the AUR
  --noconfirm         Assume 'yes' for all confirmation prompts
  --lang              Change the lang of gonob
`)
	case "version":
		scolor.BoldWhite.DisplayText("gonob@" + version + "\nlibalpm@" + fmt.Sprintf("%s", dyalpm.Version()) + "\nhttps://github.com/SnowsSky/gonob" + "\n")
	default:
		scolor.BoldYellow.DisplayText("==> " + translations.Translate("warning_string") + " : " + translations.Translate("unknown_command") + "\n")
	}
}
