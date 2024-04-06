package packer

import (
	"cmp"
	"github.com/samber/lo"
	"slices"
	"space-go/internal/model"
	"time"
)

type Sorter func(a, b *Polyomino) int

var (
	ReverseAreaSorter = func(a, b *Polyomino) int {
		return -cmp.Compare(area(a.Matrix), area(b.Matrix))
	}

	AreaSorter = func(a, b *Polyomino) int {
		return cmp.Compare(area(a.Matrix), area(b.Matrix))
	}

	ReverseCellsAmountSorter = func(a, b *Polyomino) int {
		return -cmp.Compare(countCells(a.Matrix), countCells(b.Matrix))
	}

	CellsAmountSorter = func(a, b *Polyomino) int {
		return -cmp.Compare(countCells(a.Matrix), countCells(b.Matrix))
	}

	RandomSorter = func(a, b *Polyomino) int {
		return lo.Sample([]int{-1, 0, 1})
	}
)

type GridStat struct {
	TotalGarbagePieces int
	TotalGarbageCells  int
	TotalCells         int
	AverageCost        float64
}

func Stat(grid model.Matrix) (stat GridStat) {
	index := map[int]struct{}{}
	traverse(grid, func(_, _, val int) bool {
		stat.TotalCells++
		if val == 0 {
			return true
		}

		stat.TotalGarbageCells++
		if _, ok := index[val]; !ok {
			stat.TotalGarbagePieces++
			index[val] = struct{}{}
		}

		return true
	})

	return stat
}

type DuploPacker struct{}

func (p DuploPacker) Pack(w, h int, piecesOfGarbage map[string]model.Garbage, scouting bool, minTiles int) map[string]model.Garbage {
	polyominos := make([]*Polyomino, 0, len(piecesOfGarbage))
	pid2polyomino := map[int]*Polyomino{}
	pid := 1
	for gid, garbage := range piecesOfGarbage {
		if len(garbage) == 0 {
			continue
		}
		g, gw, gh := garbage.Normalize()
		mat := model.EmptyMatrix(gw, gh)
		for _, cell := range g {
			mat[cell[1]][cell[0]] = 1
		}

		p := &Polyomino{
			Matrix:    mat,
			ID:        pid,
			GarbageID: gid,
		}

		polyominos = append(polyominos, p)
		pid2polyomino[p.ID] = p

		pid++
	}

	timeout := time.Millisecond * 50
	if scouting {
		timeout = time.Millisecond * 10
	}

	grid := BoostedRawPack(polyominos, timeout, 10_000, w, h)
	newGarbage := map[string]model.Garbage{}

	traverse(grid, func(x, y, val int) bool {
		if val == 0 {
			return true
		}

		p := pid2polyomino[val]
		newGarbage[p.GarbageID] = append(newGarbage[p.GarbageID], model.Cell{x, y})

		return true
	})

	return newGarbage
}

func BoostedRawPack(polyominos []*Polyomino, timeout time.Duration, countLimit int, w, h int) model.Matrix {
	now := time.Now()
	limit := now.Add(timeout)

	var optimal model.Matrix
	var stat GridStat

	for _, sorter := range []Sorter{
		ReverseAreaSorter,
		AreaSorter,
		ReverseCellsAmountSorter,
		CellsAmountSorter,
	} {
		avgCost, result := RawPack(polyominos, sorter, w, h)
		lstat := Stat(result)
		lstat.AverageCost = avgCost
		if isStatBetter(lstat, stat) {
			stat = lstat
			optimal = result
		}

		if time.Now().After(limit) {
			return optimal
		}
	}

	for i := 0; i < countLimit; i++ {
		avgCost, result := RawPack(polyominos, RandomSorter, w, h)
		lstat := Stat(result)
		lstat.AverageCost = avgCost
		if isStatBetter(lstat, stat) {
			stat = lstat
			optimal = result
		}

		if time.Now().After(limit) {
			return optimal
		}
	}

	return optimal
}

func isStatBetter(stat1, stat2 GridStat) bool {
	if stat1.AverageCost > stat2.AverageCost {
		return true
	}
	if stat1.TotalGarbagePieces > stat2.TotalGarbagePieces {
		return true
	}

	return false
}

func RawPack(polyominos []*Polyomino, sorter func(a, b *Polyomino) int, w, h int) (averageCost float64, mat model.Matrix) {
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

	slices.SortFunc(polyominos, sorter)

	for _, p := range polyominos {
		found := false
		maxCost := 0
		x := 0
		y := 0
		rot := 0
		sp := p.Matrix
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
		used := rotateN(rot, p.Matrix)
		traverse(used, func(dx, dy, dv int) bool {
			if dv == 0 {
				return true
			}

			rx, ry := x+dx, y+dy
			grid[ry][rx] = p.ID
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

	averageCost = 0.0
	traverse(costMap, func(x, y, val int) bool {
		averageCost += float64(val) / float64(w*h)
		return true
	})

	return averageCost, grid
}

type Polyomino struct {
	Matrix    model.Matrix
	ID        int
	GarbageID string
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
