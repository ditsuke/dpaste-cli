package main

import (
	"os"
)

func main() {
	app := getApp()
	_ = app.Run(os.Args)
}
