package main

import (
	"gonob/wrapper"
	"os"
)

func main() {
	wrapper.InitAlpm()
	parser(os.Args[1:])
}
