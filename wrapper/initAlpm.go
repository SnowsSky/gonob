package wrapper

import (
	"fmt"

	"github.com/Jguer/go-alpm/v2"
)

func InitAlpm() {
	h, er := alpm.Initialize("/", "/var/lib/pacman")
	if er != nil {
		fmt.Println(er)
		return
	}
	defer h.Release()
}
