package main

import (
	"github.com/raff/ultralight-go"
)

var (
	app = ultralight.NewApp()
)

func main() {
	defer app.Destroy()

	win := app.NewWindow(1024, 768, false, "Ultralight Browser")
	defer win.Destroy()

	NewUI(win)
	app.Run()
}
