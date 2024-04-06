package algorithmX

import (
	"fmt"
	"time"
)

// Polyomino represents a polyomino with a two-dimensional array Tiles and an id field for assignment.
type Polyomino struct {
	Tiles [][]bool
	id    int
}

// PrintPolyomino prints a given polyomino.
func PrintPolyomino(p Polyomino) {
	for _, vec := range p.Tiles {
		for _, e := range vec {
			if e {
				fmt.Print("X")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
}

/*
packingMatrix constructs the corresponding exact cover matrix for the polyomino packing problem (fit all given
polyominoes perfectly in a wxh rectangle with rotation and flipping being allowed). In the table the options
where the polyomino is out of bounds need to be filtered out.

+-------------------------------+----------------------+--------------------+
|               -               | Tile covered (w * h) | Polyomino used (n) |
+-------------------------------+----------------------+--------------------+
| Polyomino 1 state 1 on (1, 1) |                      |                    |
| Polyomino 1 state 2 on (1, 1) |                      |                    |
| ...                           |                      |                    |
| Polyomino 1 state 1 on (1, 2) |                      |                    |
| ...                           |                      |                    |
| Polyomino 1 state 8 on (h, w) |                      |                    |
| Polyomino 2 state 1 on (1, 1) |                      |                    |
| ...                           |                      |                    |
| Polyomino n state 1 on (h, w) |                      |                    |
+-------------------------------+----------------------+--------------------+
*/
func packingMatrix(polys []Polyomino, w int, h int) *Matrix {
	m := Initialize(w*h+len(polys), 0)

	// enumerate polyominoes from 1 to n
	for i := range polys {
		polys[i].id = i + 1
	}

	// find all orientations of polyominoes
	var options []Polyomino
	for _, e := range polys {
		uniqRotations(&options, &e)
	}

	for _, p := range options {
		for pCol := 0; pCol <= w-len(p.Tiles[0]); pCol++ {
			for pRow := 0; pRow <= h-len(p.Tiles); pRow++ {
				var row []int

				// tile (j+pCol, i+pRow) covered
				for i := 0; i < len(p.Tiles); i++ {
					for j := 0; j < len(p.Tiles[0]); j++ {
						if p.Tiles[i][j] {
							row = append(row, (j+pCol)+(i+pRow)*w+1)
						}
					}
				}

				row = append(row, w*h+p.id) // id-th polyomino used
				AddRow(m, row)
			}
		}
	}

	return m
}

// printPacking prints out the packing solution rows with "A" for the first polyomino, "B" for the second, ...
func printPacking(rows [][]bool, w int, h int) {
	var output [][]int32
	for i := 0; i < h; i++ {
		output = append(output, make([]int32, w))
	}

	for i := 0; i < len(rows); i++ {
		// find id of polyomino
		var id int
		for j := len(rows[i]) - 1; j >= 0; j-- {
			if rows[i][j] {
				id = j - w*h
				break
			}
		}

		for j := 0; j < len(rows[0]); j++ {
			if rows[i][j] && j != id+w*h {
				// reverse (x, y) -> x + y * w
				output[j/w][j%w] = rune('A' + id)
			}
		}
	}

	for i := 0; i < len(output); i++ {
		for j := 0; j < len(output[0]); j++ {
			fmt.Printf("%c", output[i][j])
			fmt.Print(" ")
		}
		fmt.Println("")
	}
}

// SolvePacking solves the packing problem and displays the search time.
func SolvePacking(polys []Polyomino, w int, h int) {
	m := packingMatrix(polys, w, h)

	start := time.Now()
	rows := FindRows(m)
	elapsed := time.Since(start)

	printPacking(rows, w, h)

	fmt.Printf("Time: %s\n", elapsed)
}

// CountPacking counts the number of solutions to the packing problem and displays the search time.
func CountPacking(polys []Polyomino, w int, h int) {
	m := packingMatrix(polys, w, h)

	start := time.Now()
	coll := FindAll(m)
	elapsed := time.Since(start)

	fmt.Printf("Solutions: %d\n", len(coll))
	fmt.Printf("Time: %s\n", elapsed)
}

// uniqRotations returns the up to 8 possibilities from rotation and flipping filtered for uniqueness
func uniqRotations(coll *[]Polyomino, p *Polyomino) {
	var rotations []Polyomino
	for i := 0; i < 8; i++ {
		q := rotatePolyomino(*p, i)
		if !contains(&rotations, &q) {
			rotations = append(rotations, rotatePolyomino(*p, i))
		}
	}

	for _, e := range rotations {
		*coll = append(*coll, e)
	}
}

// rotatePolyomino returns for 0 <= d <= 3 the polyomino p with d 90 degree right turns and for 4 <= d <= 7 the
// left-right mirrored polyomino p with (d - 4) 90 degree turns
func rotatePolyomino(p Polyomino, d int) Polyomino {
	d = d % 8

	pNew := Polyomino{id: p.id, Tiles: p.Tiles}

	if d >= 4 {
		// mirror matrix

		var tiles [][]bool
		for i := 0; i < len(pNew.Tiles); i++ {
			tiles = append(tiles, make([]bool, len(pNew.Tiles[0])))
		}

		for i := 0; i < len(pNew.Tiles); i++ {
			for j := 0; j < len(pNew.Tiles[0]); j++ {
				tiles[i][len(pNew.Tiles[0])-j-1] = pNew.Tiles[i][j]
			}
		}

		pNew.Tiles = tiles
		d = d - 3
	}

	for i := 0; i < d; i++ {
		rotate90(&pNew)
	}

	return pNew
}

// rotate90 rotates a polyomino 90 degrees clockwise
func rotate90(p *Polyomino) {
	h := len(p.Tiles)
	w := len(p.Tiles[0])

	var tiles [][]bool
	for i := 0; i < w; i++ {
		tiles = append(tiles, make([]bool, h))
	}

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			tiles[j][h-i-1] = p.Tiles[i][j]
		}
	}

	p.Tiles = tiles
}

// equals checks if two polyominoes are equal
func equals(p *Polyomino, q *Polyomino) bool {
	if len(p.Tiles) != len(q.Tiles) || len(p.Tiles[0]) != len(q.Tiles[0]) {
		return false
	}

	for i := 0; i < len(p.Tiles); i++ {
		for j := 0; j < len(p.Tiles[0]); j++ {
			if p.Tiles[i][j] != q.Tiles[i][j] {
				return false
			}
		}
	}
	return true
}

// contains checks if polyomino vector contains the polyomino p by the trivial approach.
func contains(polys *[]Polyomino, p *Polyomino) bool {
	for _, e := range *polys {
		if equals(&e, p) {
			return true
		}
	}
	return false
}
