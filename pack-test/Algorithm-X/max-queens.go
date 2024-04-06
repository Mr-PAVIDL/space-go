package algorithmX

import (
	"fmt"
	"time"
)

/*
maxQueensMatrix constructs the corresponding exact cover matrix for the n queens problem (placing n non-attacking
chess queens on an nxn board). For each row and column we use primary constraints to ensure a queen is placed, in
the diagonals we use secondary constraints to ensure that at most one queen is placed. The diagonals pointing
downwards are characterised by c = x - y with c in -(n-1):n-1 constant for all field (x, y) in the diagonal. The
diagonals pointing upwards are characterised by c = x + y with c in 2:2n constant for all field (x, y) in the
diagonal.

+-----------------+------------------+---------------------+---------------------------------+------------------------------+
|        -        | Queen on row (n) | Queen on column (n) | Queen on down diagonal (2n - 1) | Queen on up diagonal (2n -1) |
+-----------------+------------------+---------------------+---------------------------------+------------------------------+
| Queen on (1, 1) |                  |                     |                                 |                              |
| Queen on (1, 2) |                  |                     |                                 |                              |
| ...             |                  |                     |                                 |                              |
| Queen on (2, 1) |                  |                     |                                 |                              |
| ...             |                  |                     |                                 |                              |
| Queen on (n, n) |                  |                     |                                 |                              |
+-----------------+------------------+---------------------+---------------------------------+------------------------------+
*/
func maxQueensMatrix(n int) *Matrix {
	m := Initialize(2*n, 4*n-2)
	for x := 1; x <= n; x++ {
		for y := 1; y <= n; y++ {
			AddRow(m, []int{x, y + n, ((x - y) + n) + 2*n, ((x + y) - 1) + (2*n + 2*n - 1)})
		}
	}

	return m
}

// PrintMaxQueens print a solution to the n queens problem and displays the search time.
func PrintMaxQueens(n int) {
	// solve exact cover problem
	m := maxQueensMatrix(n)

	start := time.Now()
	solution := FindFirst(m)
	elapsed := time.Since(start)

	for x := 1; x <= n; x++ {
		fmt.Print("|")
		for y := 1; y <= n; y++ {
			if vecContains(solution, x+(y-1)*n) {
				fmt.Print("X|")
			} else {
				fmt.Print(" |")
			}
		}
		fmt.Println("")
	}

	fmt.Printf("Time: %s\n\n", elapsed)
}

// CountMaxQueens counts the number of solutions to the n queens problem and displays the search time. The sequence
// can be found in "The On-Line Encyclopedia of Integer Sequences" under A000170.
func CountMaxQueens(n int) {
	m := maxQueensMatrix(n)
	start := time.Now()
	solution := FindAll(m)
	elapsed := time.Since(start)

	fmt.Printf("Solutions: %d", len(solution))
	fmt.Println("")

	fmt.Printf("Time: %s\n\n", elapsed)
}

// vecContains checks if vec contains the element el by the trivial approach.
func vecContains(vec []int, el int) bool {
	for _, e := range vec {
		if e == el {
			return true
		}
	}
	return false
}

// (c) Mia Muessig
