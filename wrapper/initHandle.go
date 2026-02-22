package wrapper

import (
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func InitHandle() *alpm.Handle {
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}

	return &handle
}
