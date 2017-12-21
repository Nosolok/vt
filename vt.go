package main

import (
	"os"

	"github.com/Nosolok/vt/modules"
)

func main() {
	path, _ := os.Getwd()
	modules.Find(path)
}
