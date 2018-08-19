package main

import (
	"image"
	"image/color"
	"syscall/js"
	"time"
)

func drawWorld(d *Display, w World) {
	d.FillRect(image.White, image.Rect(0, 0, 640, 480))

	for c := range w {
		d.FillRect(color.Black, image.Rect(c.X*10, c.Y*10, (c.X*10)+10, (c.Y*10)+10))
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

			Cell{14, 15}: struct{}{},
			Cell{15, 15}: struct{}{},
			Cell{16, 15}: struct{}{},
		}
		display *Display

		stopper = make(chan chan struct{}, 1)
	)

	js.Global().Set("Conway", map[string]interface{}{
		"init": js.NewCallback(func(args []js.Value) {
			display = NewDisplay(args[0])
			drawWorld(display, world)

			args[0].Call("addEventListener", "click", js.NewEventCallback(js.PreventDefault, func(ev js.Value) {
				if len(stopper) != 0 {
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
		}),

		"start": js.NewCallback(func(args []js.Value) {
			go func() {
				done := make(chan struct{})
				select {
				case stopper <- done:
				default:
					return
				}

				fps := time.NewTicker(time.Second / 5)
				defer fps.Stop()

				for {
					world = world.Next()
					drawWorld(display, world)

					select {
					case <-done:
						return
					case <-fps.C:
					}
				}
			}()
		}),

		"stop": js.NewCallback(func(args []js.Value) {
			select {
			case done := <-stopper:
				close(done)
			default:
			}
		}),
	})

	select {}
}
