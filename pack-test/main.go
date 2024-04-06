package main

import (
	"fmt"
	"sort"
)

type Cell [2]int
type Garbage []Cell

type model struct{}
type PiecePlacement struct {
	Origin Cell    // The top-left corner where the piece is placed
	Piece  Garbage // The relative cells of the piece from the origin
}

// A modified structure to hold the packing state and results.
type PackingState struct {
	grid    [][]bool
	width   int
	height  int
	pieces  map[string]Garbage
	order   []string           // Keep an order of pieces keys
	results map[string]Garbage // Store the final positions
}

// Checks if a piece can be placed at a given position.
func (s *PackingState) canPlace(piece Garbage, x, y int) bool {
	for _, cell := range piece {
		cx := x + cell[0]
		cy := y + cell[1]
		if cx >= s.width || cy >= s.height || cx < 0 || cy < 0 || s.grid[cy][cx] {
			return false
		}
	}
	return true
}

// Place or remove a piece on the grid and record its placement.
func (s *PackingState) place(pieceKey string, piece Garbage, x, y int, put bool) {
	for _, cell := range piece {
		s.grid[y+cell[1]][x+cell[0]] = put
	}
	if put {
		// Convert the piece to absolute positions and store it.
		var placement Garbage
		for _, cell := range piece {
			placement = append(placement, Cell{x + cell[0], y + cell[1]})
		}
		s.results[pieceKey] = placement
	}
}

// The recursive backtracking function.
func (s *PackingState) backtrack(index int) bool {
	if index == len(s.order) {
		return true // All pieces placed
	}
	key := s.order[index]
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			if s.canPlace(s.pieces[key], x, y) {
				s.place(key, s.pieces[key], x, y, true)
				if s.backtrack(index + 1) {
					return true
				}
				s.place(key, s.pieces[key], x, y, false) // Backtrack
			}
		}
	}
	return false
}

// Initiates the packing process and returns the placement results.
func (s *PackingState) Pack(width, height int, garbage map[string]Garbage) map[string]Garbage {
	s.width = width
	s.height = height
	s.grid = make([][]bool, height)
	for i := range s.grid {
		s.grid[i] = make([]bool, width)
	}
	s.pieces = garbage
	s.results = make(map[string]Garbage)

	// Sort pieces by their size in descending order for heuristic.
	for key := range s.pieces {
		s.order = append(s.order, key)
	}
	sort.Slice(s.order, func(i, j int) bool {
		return len(s.pieces[s.order[i]]) > len(s.pieces[s.order[j]])
	})

	if !s.backtrack(0) {
		fmt.Println("Couldn't fit all pieces.")
	}

	return s.results
}

func printResultsTable(width, height int, results map[string]Garbage) {
	// Initialize the grid
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = '.' // Use '.' to represent empty cells
		}
	}

	// Assign a unique symbol to each piece
	symbols := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	pieceSymbols := make(map[string]rune)
	i := 0
	for key := range results {
		pieceSymbols[key] = rune(symbols[i])
		i++
	}

	// Fill the grid with the pieces
	for key, cells := range results {
		symbol := pieceSymbols[key]
		for _, cell := range cells {
			grid[cell[1]][cell[0]] = symbol
		}
	}

	// Print the grid
	fmt.Println("Packing Result Table:")
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Println()
	}
}

func main() {
	garbage := map[string]Garbage{
		"2FwbGL": {{0, 0}, {0, 1}, {1, 1}, {0, 2}, {1, 2}, {1, 3}},
		"2FycEL": {{3, 0}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {2, 2}, {2, 3}},
		"2HQaLA": {{6, 0}, {6, 1}, {7, 1}},
		"2HSaqf": {{4, 2}, {3, 3}, {4, 3}, {5, 3}, {3, 4}, {2, 5}, {3, 5}},
	}
	packer := PackingState{}
	results := packer.Pack(8, 11, garbage)
	fmt.Println("Packing Results:", results)
}
