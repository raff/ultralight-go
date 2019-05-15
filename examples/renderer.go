package main

import (
	"github.com/raff/ultralight-go"
)

func main() {
	c := ultralight.NewConfig()
	defer c.Destroy()

	r := ultralight.NewRenderer(c)
	defer r.Destroy()

	v := r.NewView(200, 200, false)
	defer v.Destroy()

	done := false

	v.OnFinishLoading(func() {
		r.Render()
		v.WriteToPNG("result.png")
		done = true
	})

	v.LoadHTML("<h1>Hello!</h1><p>Welcome to Ultralight!</p>")

	for !done {
		r.Update()
	}
}
