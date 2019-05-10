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
		view := win.View()
		win.SetTitle(view.Title())
		fmt.Println("finish loading", view.URL())
	})

	win.View().OnUpdateHistory(func() {
		fmt.Println("update history")
	})

	win.View().OnDOMReady(func() {
		fmt.Println("DOM ready")

		if true {
			// test EvaluateScript and various JSValue methods

			values := []string{
				"'hello'",
				"42",
				"true",
				"undefined",
				"null",
				"{a: 1, b: 2}",
				"[1,2,3]",
				"new Date()",
			}

			for _, s := range values {
				v := win.View().EvaluateScript(s)
				fmt.Printf("%v t=%v o=%v, s=%v, N=%v, b=%v, a=%v, d=%v u=%v n=%v %q\n",
					s,
					v.Type(),
					v.IsObject(),
					v.IsString(),
					v.IsNumber(),
					v.IsBoolean(),
					v.IsArray(),
					v.IsDate(),
					v.IsUndefined(),
					v.IsNull(),
                                        v.String(),
				)
			}
		}
	})

	win.View().OnConsoleMessage(func(source ultralight.MessageSource, level ultralight.MessageLevel,
		message string, line uint, col uint, sourceId string) {
		fmt.Printf("CONSOLE source=%v level=%v id=%q line=%c col=%v %v\n",
			source, level, sourceId, line, col, message)
	})

	/*
		app.OnUpdate(func() {
			fmt.Println("app should update")
		})
	*/

	app.Run()
}
