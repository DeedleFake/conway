package main

import (
	"fmt"
	"image"
	"image/color"
	"syscall/js"
)

// Display is a minimally featured wrapper around a canvas DOM element
// instance.
type Display struct {
	canvas js.Value
	ctx    js.Value
}

// NewDisplay returns a new Display instance that wraps the given
// canvas element. It does not check to make sure that it has actually
// been passed a canvas.
func NewDisplay(canvas js.Value) *Display {
	return &Display{
		canvas: canvas,
		ctx:    canvas.Call("getContext", "2d"),
	}
}

// FillRect fills the rectangle r with the color c.
func (d *Display) FillRect(c color.Color, r image.Rectangle) {
	d.ctx.Set("fillStyle", htmlColor(c))
	d.ctx.Call("fillRect", r.Min.X, r.Min.Y, r.Dx(), r.Dy())
}

func htmlColor(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%v, %v, %v, %v)",
		r*255/0xFFFF,
		g*255/0xFFFF,
		b*255/0xFFFF,
		float64(a)/0xFFFF,
	)
}
