package main

import (
	"gonob/wrapper"
)

func parser(args []string) {
	switch args[0] {
	case "upgrade":
		wrapper.Upgrade()
	}
}
