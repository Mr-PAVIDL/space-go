package packer

import (
	"cmp"
	"slices"
	"space-go/internal/model"
)

func Pack(piecesOfGarbage map[string]model.Garbage, w, h int) (model.Matrix, map[string]model.Garbage) {
	polyominos := make([]*polyomino, 0, len(piecesOfGarbage))
	pid2polyomino := map[int]*polyomino{}

	pid := 1
	for gid, garbage := range piecesOfGarbage {
		g, gw, gh := garbage.Normalize()
		mat := model.EmptyMatrix(gw, gh)
		for _, cell := range g {
			mat[cell[1]][cell[0]] = 1
		}

		p := &polyomino{
			matrix:    mat,
			id:        pid,
			garbageId: gid,
		}

		polyominos = append(polyominos, p)
		pid2polyomino[p.id] = p

		pid++
	}

	// algorithm

	grid := model.EmptyMatrix(w, h)
	costMap := model.EmptyMatrix(w, h)
	for _, y := range []int{0, h - 1} {
		for x := 0; x < w; x++ {
			costMap[y][x]++
		}
	}
	for _, x := range []int{0, w - 1} {
		for y := 0; y < w; y++ {
			costMap[y][x]++
		}
	}

	slices.SortFunc(polyominos, func(a, b *polyomino) int {
		return -cmp.Compare(area(a.matrix), area(b.matrix))
	})

	for _, p := range polyominos {
		found := false
		maxCost := 0
		x := 0
		y := 0
		rot := 0
		sp := p.matrix
		for rotationId := 0; rotationId < 4; rotationId++ {
			traverse(grid, func(gx, gy, gv int) bool {
				valid := true
				localCost := 0
				traverse(sp, func(dx, dy, dv int) bool {
					rx, ry := gx+dx, gy+dy
					if rx >= w || ry >= h {
						valid = false
						return false
					}
					if grid[ry][rx] > 0 && dv > 0 {
						valid = false
						return false
					}
					if dv > 0 {
						localCost += costMap[ry][rx]
					}
					return true
				})

				if !valid {
					return true
				}

				if maxCost < localCost {
					x = gx
					y = gy
					rot = rotationId
					maxCost = localCost
					found = true
				}

				return true
			})

			sp = sp.RotateCW()
		}

		if !found {
			continue
		}
		used := rotateN(rot, p.matrix)
		traverse(used, func(dx, dy, dv int) bool {
			if dv == 0 {
				return true
			}

			rx, ry := x+dx, y+dy
			grid[ry][rx] = p.id
			flowerVicinity(rx, ry, func(fx, fy int) bool {
				if fx < 0 || fy < 0 {
					return true
				}
				if fx >= w || fy >= h {
					return true
				}
				costMap[fy][fx]++
				return true
			})
			return true
		})
	}

	newGarbage := map[string]model.Garbage{}

	traverse(grid, func(x, y, val int) bool {
		if val == 0 {
			return true
		}

		p := pid2polyomino[val]
		newGarbage[p.garbageId] = append(newGarbage[p.garbageId], model.Cell{x, y})

		return true
	})

	return grid, newGarbage
}

type polyomino struct {
	matrix    model.Matrix
	id        int
	garbageId string
}

func rotateN(n int, m model.Matrix) model.Matrix {
	for i := 0; i < n; i++ {
		m = m.RotateCW()
	}

	return m
}

func flowerVicinity(x, y int, fn func(x, y int) bool) {
	if !fn(x-1, y) {
		return
	}
	if !fn(x+1, y) {
		return
	}
	if !fn(x, y-1) {
		return
	}
	if !fn(x, y+1) {
		return
	}
}

func area(mat model.Matrix) int {
	if len(mat) == 0 {
		return 0
	}
	return len(mat) * len(mat[0])
}

func countCells(mat model.Matrix) int {
	c := 0
	for _, row := range mat {
		for _, val := range row {
			if val <= 0 {
				continue
			}
			c += 1
		}
	}

	return c
}

func traverse(mat model.Matrix, fn func(x, y, val int) bool) {
	for y, row := range mat {
		for x, val := range row {
			if !fn(x, y, val) {
				return
			}
		}
	}
}
