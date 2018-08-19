package main

// Cell stores coordinates for a cell in the world.
type Cell struct {
	X, Y int
}

// World stores the state of the world as a set of live cells. If a
// cell is present in the map as a key, then it is considered to be
// live.
//
// For both convenience and compatibility, the IsLive method should be
// used instead of manually checking for presence.
type World map[Cell]struct{}

// IsLive returns true if the given cell is live.
func (w World) IsLive(c Cell) bool {
	_, ok := w[c]
	return ok
}

// Neighbors returns the number of live neighbors of the given cell,
// whether or not that cell is itself live. It will always return a
// number in the range [0, 8].
func (w World) Neighbors(c Cell) (n int) {
	for y := c.Y - 1; y <= c.Y+1; y++ {
		for x := c.X - 1; x <= c.X+1; x++ {
			cell := Cell{X: x, Y: y}
			if cell == c {
				continue
			}

			if w.IsLive(cell) {
				n++
			}
		}
	}

	return n
}

// NextState returns the state that the given cell should have in the
// next world state given both its state and the state of its
// neighbors in the current world state.
func (w World) NextState(c Cell) bool {
	if w.IsLive(c) {
		return w.nextStateLive(c)
	}

	return w.nextStateDead(c)
}

// nextStateLive returns the next state of the cell at c assuming that
// the cell is currently live.
func (w World) nextStateLive(c Cell) bool {
	n := w.Neighbors(c)

	if (n < 2) || (n > 3) {
		return false
	}

	return true
}

// nextStateDead returns the next state of the cell at c assuming that
// the cell is currently dead.
func (w World) nextStateDead(c Cell) bool {
	n := w.Neighbors(c)

	if n == 3 {
		return true
	}

	return false
}

// Next returns a new World instance containing the next state
// following the state contained by w.
func (w World) Next() World {
	next := make(World, len(w))
	for c := range w {
		if w.nextStateLive(c) {
			next[c] = struct{}{}
		}

		for y := c.Y - 1; y <= c.Y+1; y++ {
			for x := c.X - 1; x <= c.X+1; x++ {
				cell := Cell{X: x, Y: y}
				if cell == c {
					continue
				}

				if !w.IsLive(cell) && w.nextStateDead(cell) {
					next[cell] = struct{}{}
				}
			}
		}
	}

	return next
}
