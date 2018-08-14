package main

type Cell struct {
	X, Y int
}

type World map[Cell]struct{}

func (w World) Live(c Cell) bool {
	_, ok := w[c]
	return ok
}

func (w World) Neighbors(c Cell) (n int) {
	for y := c.Y - 1; y <= c.Y+1; y++ {
		for x := c.X - 1; x <= c.X+1; x++ {
			cell := Cell{X: x, Y: y}
			if cell == c {
				continue
			}

			if w.Live(cell) {
				n++
			}
		}
	}

	return n
}

func (w World) LiveState(c Cell) bool {
	n := w.Neighbors(c)

	if (n < 2) || (n > 3) {
		return false
	}

	return true
}

func (w World) DeadState(c Cell) bool {
	n := w.Neighbors(c)

	if n == 3 {
		return true
	}

	return false
}

func (w World) Next() World {
	next := make(World, len(w))
	for c := range w {
		if w.LiveState(c) {
			next[c] = struct{}{}
		}

		for y := c.Y - 1; y <= c.Y+1; y++ {
			for x := c.X - 1; x <= c.X+1; x++ {
				cell := Cell{X: x, Y: y}
				if cell == c {
					continue
				}

				if !w.Live(cell) && w.DeadState(cell) {
					next[cell] = struct{}{}
				}
			}
		}
	}

	return next
}
