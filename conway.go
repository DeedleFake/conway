package main

import (
	"context"
	"image"
	"syscall/js"
	"time"
)

func drawWorld(d *Display, w World) {
	d.FillRect(image.White, image.Rect(0, 0, 640, 480))

	for c := range w {
		d.FillRect(image.Black, image.Rect(c.X*10, c.Y*10, (c.X*10)+10, (c.Y*10)+10))
	}
}

func main() {
	var (
		world = World{
			Cell{4, 5}: struct{}{},
			Cell{5, 6}: struct{}{},
			Cell{6, 6}: struct{}{},
			Cell{6, 5}: struct{}{},
			Cell{6, 4}: struct{}{},
		}
		display *Display

		stopper func()
	)

	js.Global().Set("initConway", js.NewCallback(func(args []js.Value) {
		display = NewDisplay(args[0])
		drawWorld(display, world)

		args[0].Call("addEventListener", "click", js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
			if stopper != nil {
				return
			}

			bounds := ev.Get("target").Call("getBoundingClientRect")
			x, y := ev.Get("clientX").Int()-bounds.Get("x").Int(), ev.Get("clientY").Int()-bounds.Get("y").Int()

			cell := Cell{x / 10, y / 10}
			switch _, ok := world[cell]; ok {
			case true:
				delete(world, cell)
			case false:
				world[cell] = struct{}{}
			}

			drawWorld(display, world)
		}))
	}))

	js.Global().Set("startConway", js.NewCallback(func(args []js.Value) {
		if stopper != nil {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		stopper = cancel

		fps := time.NewTicker(time.Second / 5)
		defer fps.Stop()

		for {
			drawWorld(display, world)
			world = world.Next()

			select {
			case <-ctx.Done():
				return
			case <-fps.C:
			}
		}
	}))

	js.Global().Set("stopConway", js.NewCallback(func(args []js.Value) {
		if stopper != nil {
			stopper()
			stopper = nil
		}
	}))

	select {}
}
