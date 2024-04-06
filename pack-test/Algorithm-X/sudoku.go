package algorithmX

import (
	"fmt"
	"time"
)

// Sudoku represents a sudoku board with a two-dimensional array. Zero entries encode empty cells.
type Sudoku struct {
	Board [9][9]int
}

// ShowSudoku prints out the sudoku s.
func ShowSudoku(s Sudoku) {
	for i := 0; i <= 8; i++ {
		if i%3 == 0 {
			fmt.Println("+-------+-------+-------+")
		}
		for j := 0; j <= 8; j++ {
			if j%3 == 0 {
				fmt.Print("| ")
			}
			if s.Board[i][j] != 0 {
				fmt.Printf("%d ", s.Board[i][j])
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println("|")
	}

	fmt.Println("+-------+-------+-------+")
}

/*
maxQueensMatrix constructs the corresponding exact cover matrix for the sudoku problem (in each row, column and
3x3 group every number 1:9 is contained exactly once). To encode a specific sudoku board s we use ForceOption to
preselect the options already chosen by the entries in s.

+-------------+----------------------------------------------------------------------------------------------------------------+
|      -      |    Field filled (n^2)   Row has all numbers (n^2)   Column has all numbers (n^2)   Group has all numbers (n^2) |
+-------------+----------------------------------------------------------------------------------------------------------------+
| 1 at (1, 1) |                                                                                                                |
| ...         |                                                                                                                |
| n at (1, 1) |                                                                                                                |
| 1 at (1, 2) |                                                                                                                |
| ...         |                                                                                                                |
| n at (1, 2) |                                                                                                                |
| ...         |                                                                                                                |
| 1 at (n, n) |                                                                                                                |
| ...         |                                                                                                                |
| n at (n, n) |                                                                                                                |
+-------------+----------------------------------------------------------------------------------------------------------------+
*/
func sudokuMatrix(s *Sudoku) *Matrix {
	m := Initialize(4*81, 0)

	for x := 1; x <= 9; x++ {
		for y := 1; y <= 9; y++ {
			for n := 1; n <= 9; n++ {
				g := (x-1)/3*3 + (y-1)/3 + 1 // group number

				// [field (x, y) filled, row x contains n, column y contains n, group g contains n]
				AddRow(m, []int{x + (y-1)*9, 81 + n + (x-1)*9, 2*81 + n + (y-1)*9, 3*81 + n + (g-1)*9})
			}
		}
	}

	for x := 1; x <= 9; x++ {
		for y := 1; y <= 9; y++ {
			if s.Board[x-1][y-1] != 0 {
				// (x, y) already filled in sudoku s
				row := s.Board[x-1][y-1] + (y-1)*9 + (x-1)*81
				ForceOption(m, x+(y-1)*9, row)
			}
		}
	}

	return m
}

// SolveSudoku solves the sudoku s in-situ and displays the search time.
func SolveSudoku(s *Sudoku) bool {
	m := sudokuMatrix(s)
	start := time.Now()
	solution := FindFirst(m)
	elapsed := time.Since(start)

	fmt.Printf("Time: %s\n\n", elapsed)

	if &solution != nil {
		for _, e := range solution {
			n := (e-1)%9 + 1
			y := (((e - 1) / 9) % 9) + 1
			x := (e-1)/81 + 1

			s.Board[x-1][y-1] = n
		}

		return true
	}

	return false
}
