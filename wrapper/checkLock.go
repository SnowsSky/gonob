package wrapper

import "os"

var lockfile = "/var/lib/pacman/db.lck"

func CheckLock() bool {
	if _, err := os.Stat(lockfile); os.IsNotExist(err) {
		// the folder does not exist.
		return false
	}
	return true
}
