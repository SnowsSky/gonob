package main

import (
	"gonob/aur"
)

func parser(args []string) {
	switch args[0] {
	case "install":
		if args[1] == "--aur" {
			aur.Install(args[2:])
		}

	}
}
