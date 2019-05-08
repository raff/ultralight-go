package main

import (
	"flag"
	"fmt"

	"github.com/raff/ultralight-go"
)

func main() {
	title := flag.String("title", "Hello from Go", "Window title")
	width := flag.Uint("width", 600, "Window width")
	height := flag.Uint("height", 400, "Window height")
	full := flag.Bool("fullscreen", false, "Go full screen")
	flag.Parse()

	app := ultralight.NewApp()
	defer app.Destroy()

	win := app.NewWindow(*width, *height, *full, *title)
	defer win.Destroy()

	if flag.NArg() > 0 {
		win.View().LoadURL(flag.Arg(0)) // should be a URL
	} else {
		win.View().LoadHTML(`<html>
                <p>Resize the browser window to fire the <code>resize</code> event.</p>
                <p>Window height: <span id="height"></span></p>
                <p>Window width: <span id="width"></span></p>
                <script>
                const heightOutput = document.querySelector('#height');
                const widthOutput = document.querySelector('#width');

                function reportWindowSize() {
                    heightOutput.textContent = window.innerHeight;
                    widthOutput.textContent = window.innerWidth;
                }

                window.onresize = reportWindowSize;
                </script>
                </html>`)
	}

	win.OnResize(func(width, height uint) {
		fmt.Println("resize", width, height)
		win.Resize(width, height)
	})

	win.OnClose(func() {
		fmt.Println("window is closing")
	})

	//view := win.View()

	win.View().OnBeginLoading(func() {
		fmt.Println("begin loading")
	})

	win.View().OnFinishLoading(func() {
		fmt.Println("finish loading")
	})

	win.View().OnUpdateHistory(func() {
		fmt.Println("update history")
	})

	win.View().OnDOMReady(func() {
		fmt.Println("DOM ready")
	})

	/*
		app.OnUpdate(func() {
			fmt.Println("app should update")
		})
	*/

	app.Run()
}
