package commander

import (
	"github.com/samber/lo"
	"log"
	"sort"
	"space-go/internal/model"
	algorithmX "space-go/pack-test/Algorithm-X"
	"strings"
)

type Packer interface {
	Pack(width, height int, garbage map[string]model.Garbage) map[string]model.Garbage
}

type DumboPacker struct{}

type pair struct {
	Name string
	G    model.Garbage
}

func (p DumboPacker) Pack(width, height int, garbage map[string]model.Garbage) map[string]model.Garbage {
	pairs := lo.MapToSlice(garbage, func(key string, val model.Garbage) pair {
		return pair{Name: key, G: val}
	})
	sort.Slice(pairs, func(i, j int) bool {
		return len(pairs[i].G) < len(pairs[j].G)
	})
	best := pack(width, height, pairs)
	for i := 0; i < 500; i++ {
		lo.Shuffle(pairs)
		result := pack(width, height, pairs)
		if len(result) > len(best) {
			log.Println("updated best: ", len(best), "->", len(result))
			best = result
		}
	}
	if len(garbage) > 0 {
		print(best, width, height)
	}
	return best
}

func pack(width, height int, pairs []pair) map[string]model.Garbage {
	mat := model.EmptyMatrix(width, height)
	added := map[string]model.Garbage{}
	cellCount := 0
	for _, p := range pairs {
		if width*height-cellCount < len(p.G) {
			continue
		}
		if ok, g := tryFit(width, height, mat, p.G); ok {
			added[p.Name] = g // save garbage with offset
			cellCount += len(g)
		}
	}
	return added
}

func CountTiles(garbage map[string]model.Garbage) int {
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

type PackX struct{}

// Calculate the minimum and maximum x and y coordinates of garbage to determine its bounds.
func calculateBounds(garbage model.Garbage) (minX, minY, maxX, maxY int) {
	minX, minY = garbage[0][0], garbage[0][1]
	maxX, maxY = minX, minY
	for _, cell := range garbage {
		if cell[0] < minX {
			minX = cell[0]
		}
		if cell[0] > maxX {
			maxX = cell[0]
		}
		if cell[1] < minY {
			minY = cell[1]
		}
		if cell[1] > maxY {
			maxY = cell[1]
		}
	}
	return
}

//// Convert garbage to the Polyomino format expected by algorithmX.
//func garbageToPolyomino(garbage model.Garbage) algorithmX.Polyomino {
//	minX, minY, maxX, maxY := calculateBounds(garbage)
//	width, height := maxX-minX+1, maxY-minY+1
//	tiles := make([][]bool, height)
//	for i := range tiles {
//		tiles[i] = make([]bool, width)
//	}
//
//	for _, cell := range garbage {
//		x, y := cell[0]-minX, cell[1]-minY
//		tiles[y][x] = true
//	}
//
//	return algorithmX.Polyomino{Tiles: tiles}
//}

func garbageToPolyomino(g model.Garbage) algorithmX.Polyomino {
	minX, minY, maxX, maxY := calculateBounds(g)
	width, height := maxX-minX+1, maxY-minY+1
	tiles := make([][]bool, height)
	for i := range tiles {
		tiles[i] = make([]bool, width)
	}

	for _, cell := range g {
		x, y := cell[0]-minX, cell[1]-minY
		tiles[y][x] = true
	}

	return algorithmX.Polyomino{Tiles: tiles, Name: ""} // Name will be assigned later.
}

func (p PackX) Pack(width, height int, garbage map[string]model.Garbage) map[string]model.Garbage {
	var polys []algorithmX.Polyomino

	// Convert each garbage piece into a named Polyomino.
	for name, g := range garbage {
		poly := garbageToPolyomino(g)
		poly.Name = name // Assign the garbage piece's name to the Polyomino.
		polys = append(polys, poly)
	}

	packedGarbage := algorithmX.SolvePacking(polys[:10], width, height)

	return packedGarbage
}
