package algorithmX

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

func divisors(n int) []int {
	var div []int

	for i := 1; i <= n; i++ {
		if n%i == 0 {
			div = append(div, i)
		}
	}

	return div
}

// return rounded int of p/q
func truncDiv(p int, q int) int {
	return int(math.Round(float64(p) / float64(q)))
}

func ceilDiv(p int, q int) int {
	return int(math.Ceil(float64(p) / float64(q)))
}

/*
+-----------------------------------------------+--------------------+---------------------+
|                       -                       | Tile covered (n^2) |  Rectangle used (m) |
+-----------------------------------------------+--------------------+---------------------+
| Rectangle 1 (1, 1)                            |                    |                     |
| ...                                           |                    |                     |
| Rectangle 1 (n - sizeX(1), n - sizeY(1))      |                    |                     |
| Rectangle 1 rot. (1, 1)                       |                    |                     |
| ...                                           |                    |                     |
| Rectangle 1 rot. (n - sizeY(1), n - sizeX(1)) |                    |                     |
| ....                                          |                    |                     |
| Rectangle m rot. (n - sizeY(m), n - sizeX(m)) |                    |                     |
+-----------------------------------------------+--------------------+---------------------+
*/

func rectangleMatrix(rectsX []int, rectsY []int, n int) *Matrix {
	m := Initialize(n*n, len(rectsX))

	for i, _ := range rectsX {
		for _, rot := range [2]bool{true, false} {
			sx := 0
			sy := 0
			if rot {
				sx = rectsX[i]
				sy = rectsY[i]
			} else {
				sx = rectsY[i]
				sy = rectsX[i]
			}

			for px := 0; px <= n-sx; px++ {
				for py := 0; py <= n-sy; py++ {
					var row []int

					for x := px + 1; x <= px+sx; x++ {
						for y := py + 1; y <= py+sy; y++ {
							row = append(row, x+(y-1)*n)
						}
					}

					row = append(row, n*n+i+1) // i-th rectangle used
					sort.Ints(row)
					AddRow(m, row)
				}
			}
		}
	}

	return m
}

func printRectangles(rows [][]bool, n int) {
	var output [][]int32
	for i := 0; i < n; i++ {
		output = append(output, make([]int32, n))
	}

	for i := 0; i < len(rows); i++ {
		for j := 0; j < n*n; j++ {
			if rows[i][j] {
				output[j/n][j%n] = rune('A' + i)
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

func SolveMondrian(n int) {
	minPieces := 4 //TODO 9 // proof in paper
	d := divisors(n * n)

	for _, r := range d {
		if r >= minPieces {
			area := truncDiv(n*n, r)
			dA := divisors(area)

			var s []int
			for _, i := range dA {
				if i <= n && area/i <= n {
					s = append(s, i)
				}
			}

			if ceilDiv(len(s), 2) < r {
				continue
			}

			var rectsX []int
			var rectsY []int
			for i := 1; i <= ceilDiv(len(s), 2); i++ {
				rectsX = append(rectsX, s[i-1])
				rectsY = append(rectsY, truncDiv(area, s[i-1]))
			}

			m := rectangleMatrix(rectsX, rectsY, n)

			fmt.Println("---------------------------------------------")
			fmt.Println("n = " + strconv.Itoa(n) + ", r = " + strconv.Itoa(r))
			fmt.Println(rectsX)
			fmt.Println(rectsY)

			start := time.Now()
			FindFirst(m)
			elapsed := time.Since(start)

			fmt.Printf("Time: %s\n", elapsed)

			fmt.Println("---------------------------------------------")
		}
	}
}
