package main

import (
	"github.com/raff/ultralight-go"
)

func main() {
	app := ultralight.NewApp()
	defer app.Destroy()

	window := app.NewWindow(800, 600, false, "Window", ultralight.WindowFlagBorderless|ultralight.WindowFlagResizable)
	defer window.Destroy()

	window.View().LoadHTML(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<div>
			<p>This window has no border!</p>
			<p onclick="close()">click me to close the window!</p>
		</div>
	</body>
	</html>`)

	window.View().JSContext().GlobalObject().SetPropertyValue("close", func(f, this *ultralight.JSObject, args ...*ultralight.JSValue) *ultralight.JSValue {
		app.Quit()
		return nil
	})

	window.Focus()
	app.Run()
}
