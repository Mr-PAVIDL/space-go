package commander

import (
	"github.com/samber/lo"
	"log"
	"sort"
	"space-go/internal/model"
	"strings"
)

type Packer interface {
	Pack(width, height int, garbage map[string]model.Garbage) map[string]model.Garbage
}

type DumboPacker struct{}

func (p DumboPacker) Pack(width, height int, garbage map[string]model.Garbage) map[string]model.Garbage {
	type Pair struct {
		Name string
		G    model.Garbage
	}
	pairs := lo.MapToSlice(garbage, func(key string, val model.Garbage) Pair {
		return Pair{Name: key, G: val}
	})
	sort.Slice(pairs, func(i, j int) bool {
		return len(pairs[i].G) < len(pairs[j].G)
	})
	mat := model.EmptyMatrix(width, height)
	added := map[string]model.Garbage{}
	for _, p := range pairs {
		if ok, g := tryFit(width, height, mat, p.G); ok {
			added[p.Name] = g // save garbage with offset
		}
	}
	if len(garbage) > 0 {
		//print(added, width, height)
	}
	return added
}

func CalcTiles(garbage map[string]model.Garbage) int {
	s := 0
	for _, g := range garbage {
		s += len(g)
	}
	return s
}

func print(added map[string]model.Garbage, width, height int) {
	table := make([][]byte, height)
	for i := 0; i < height; i++ {
		table[i] = make([]byte, width)
		for j := 0; j < width; j++ {
			table[i][j] = '.'
		}
	}
	id2sym := map[string]byte{}
	for id, _ := range added {
		id2sym[id] = byte('a' + len(id2sym))
	}
	for id, garbage := range added {
		for _, cell := range garbage {
			table[cell[1]][cell[0]] = id2sym[id]
		}
	}
	log.Println("┌" + strings.Repeat("-", width*2+1) + "┐")
	for y := 0; y < height; y++ {
		line := ""
		line += "| "
		for x := 0; x < width; x++ {
			line += string(table[y][x])
			line += " "
		}
		log.Println(line + "|")
	}
	log.Print("└" + strings.Repeat("-", width*2+1) + "┘")
}

func tryFit(width int, height int, mat model.Matrix, g model.Garbage) (bool, model.Garbage) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if fits(mat, g, x, y) {
				for _, cell := range g {
					mat[y+cell[1]][x+cell[0]] = true
				}
				return true, g.Add(x, y)
			}
		}
	}
	return false, model.Garbage{}
}

func fits(mat model.Matrix, g model.Garbage, x int, y int) bool {
	for _, cell := range g {
		y := y + cell[1]
		x := x + cell[0]
		if y < 0 || x < 0 || y >= len(mat) || x >= len(mat[y]) {
			return false
		}
		if mat[y][x] {
			return false
		}
	}
	return true
}
