package algorithmX

import (
	"fmt"
)

// Matrix represents a sparse matrix as a toroidal double linked list.
// Each column represents a constraint, which either has to covered exactly once (primary constraint)
// or at most once (secondary constraint). Each row represents an option, which when included in the solution
// covers all columns of its ones entries.
type Matrix struct {
	columns []Column // columns[0] is root node
	numCols int
	numRows int
	sol     []*Node
}

// Column build a double linked list of headers to access to columns of the matrix
type Column struct {
	left, right *Column
	head        Node
	length      int
	id          int
}

// Node links to its neighbours and its column. It saves its row in the field value.
type Node struct {
	up, down, left, right *Node
	col                   *Column
	value                 int
}

// Initialize constructs an empty matrix with the specified number of primary (exactly once)
// and secondary constraints (at most once).
func Initialize(primary int, secondary int) *Matrix {
	numCols := primary + secondary

	// setup first header row of matrix
	m := Matrix{columns: make([]Column, numCols+1), numCols: numCols, numRows: 0}

	// root node has an empty column
	m.columns[0].head.down = &m.columns[0].head
	m.columns[0].head.up = &m.columns[0].head

	m.columns[0].right = &m.columns[1]
	m.columns[0].left = &m.columns[primary]

	// initialize primary columns
	for i := 1; i <= primary; i++ {
		m.columns[i] = Column{left: &m.columns[i-1], right: &m.columns[(i+1)%(primary+1)], id: i}
		m.columns[i].head.down = &m.columns[i].head
		m.columns[i].head.up = &m.columns[i].head
	}

	// initialize secondary columns (not linked into header row)
	for i := primary + 1; i <= numCols; i++ {
		m.columns[i] = Column{left: &m.columns[i], right: &m.columns[i], id: i}
		m.columns[i].head.down = &m.columns[i].head
		m.columns[i].head.up = &m.columns[i].head
	}

	return &m
}

// AddRow appends the sparse row indices to the matrix m. So the input {1, 4, 5} is appended
// as row (1, 0, 0, 1, 1) to a matrix of width 5. The function validate that the indices of the
// ones entries are in range and sorted in ascending order.
func AddRow(m *Matrix, indices []int) {
	// input verification
	last := -1
	for _, element := range indices {
		if element <= 0 || element > m.numCols {
			fmt.Println("Error: Index out of range.")
		}
		if element <= last {
			fmt.Println("Error: Indices not ordered")
		}
		last = element
	}

	m.numRows++

	for _, e := range indices {
		// insert new node in last row of column
		current := Node{value: m.numRows, col: &m.columns[e], down: &m.columns[e].head, up: m.columns[e].head.up}

		m.columns[e].head.up.down = &current
		m.columns[e].head.up = &current

		m.columns[e].length++
	}

	// link new nodes to left/right neighbour (head.up still points to new node)
	for i, e := range indices {
		m.columns[e].head.up.right = m.columns[indices[(i+1)%len(indices)]].head.up
		m.columns[indices[(i+1)%len(indices)]].head.up.left = m.columns[e].head.up
	}
}

// cover removes c from the selection and removes all colliding rows. If the column c has a one entry at row i,
// then remove this row from the matrix.
func cover(c *Column) {
	// hide colliding rows
	for p := c.head.down; p != &(c.head); p = p.down {
		hide(p)
	}

	// remove from header list
	c.left.right = c.right
	c.right.left = c.left
}

func hide(p *Node) {
	for q := p.right; q != p; q = q.right {
		q.up.down = q.down
		q.down.up = q.up

		q.col.length--
	}
}

// uncover undoes the deletion done by cover(c).
func uncover(c *Column) {
	// add to header list
	c.left.right = c
	c.right.left = c

	// unhide colliding rows
	for p := c.head.up; p != &(c.head); p = p.up {
		unhide(p)
	}
}

func unhide(p *Node) {
	for q := p.left; q != p; q = q.left {
		q.up.down = q
		q.down.up = q

		q.col.length++
	}
}

// mrv uses Knuths MRV heuristic, which always chooses the column that can be covered by the least number
// of rows.
func mrv(m *Matrix) *Column {
	minLen := -1
	var minCol *Column

	// check if smaller column exists
	for c := m.columns[0].right; c != &m.columns[0]; c = c.right {
		if c.length < minLen || minLen == -1 {
			minLen = c.length
			minCol = c
		}
	}

	return minCol
}

// solve does a backtracking search using the matrix data structure. The next column is selected with the MRV
// heuristic. Then the search goes through the options (rows) to cover the primary constraint (column). If first
// is true, the search is stopped after the first solution is found.
func solve(m *Matrix, collector *[][]int, first bool) []*Node {
	// problem is solved
	if m.columns[0].left == &m.columns[0] {
		var solution []int
		for _, e := range m.sol {
			solution = append(solution, e.value)
		}

		*collector = append(*collector, solution)
		return m.sol
	}

	// MRV heuristic
	col := mrv(m)
	if col.length == 0 {
		return nil
	}

	cover(col)

	// classic backtracking algorithm
	for r := col.head.down; r != &col.head; r = r.down {
		// add r to solution
		m.sol = append(m.sol, r)

		// cover r
		for n := r.right; n != r; n = n.right {
			cover(n.col)
		}

		result := solve(m, collector, first)

		// remove r from solution
		undo := m.sol[len(m.sol)-1]
		col = undo.col
		m.sol = m.sol[:len(m.sol)-1]

		// uncover r
		for n := undo.left; n != undo; n = n.left {
			uncover(n.col)
		}

		if first && result != nil {
			return result
		}
	}

	uncover(col)

	return nil
}

// FindFirst uses solve to find one solution to the exact cover problem for m. It returns the options of the solution
// represented by their row id.
func FindFirst(m *Matrix) []int {
	var coll [][]int

	solve(m, &coll, true)

	if len(coll) == 0 {
		fmt.Println("Error: No solution found.")
		return nil
	}
	return coll[0]
}

// FindRows uses solve to find one solution to the exact cover problem for m. It returns the options of the solution
// represented by their row id and content.
func FindRows(m *Matrix) [][]bool {
	// Return row ids and their content

	sol := solve(m, &[][]int{}, true)

	var rows [][]bool
	for i := 0; i < len(sol); i++ {
		rows = append(rows, make([]bool, m.numCols))
	}

	for i, e := range sol {
		rows[i][e.col.id-1] = true
		for n := e.right; n != e; n = n.right {
			rows[i][n.col.id-1] = true
		}
	}

	return rows
}

// FindAll uses solve to find all solution to the exact cover problem for m. It returns the options of the solution
// represented by their row id.
func FindAll(m *Matrix) [][]int {
	var coll [][]int

	solve(m, &coll, false)

	return coll
}

// PrettyPrint prints a non-sparse table representation of the matrix m for debug purposes.
func PrettyPrint(m *Matrix) {
	var elements [][]bool
	for i := 0; i < m.numRows; i++ {
		elements = append(elements, make([]bool, m.numCols))
	}

	for i := 1; i <= m.numCols; i++ {
		h := m.columns[i]
		for e := h.head.down; *e != h.head; e = e.down {
			elements[e.value-1][i-1] = true
		}
	}

	numPrimary := 0
	for c := m.columns[0].right; c != &m.columns[0]; c = c.right {
		numPrimary++
	}

	for x := 1; x <= m.numRows; x++ {
		fmt.Print("|")
		for y := 1; y <= m.numCols; y++ {
			if elements[x-1][y-1] {
				fmt.Print("X|")
			} else {
				fmt.Print(" |")
			}
			if y == numPrimary {
				fmt.Print("|")
			}
		}
		fmt.Println("")
	}
}

// ForceOption ensures that a specific option value has to be included in all solutions. To find the row specified
// by value in the sparse matrix a column containing a one entry in that row is given.
func ForceOption(m *Matrix, id int, value int) {
	// get column from id
	col := &m.columns[id]

	// collect columns covered by forced option to cover together later, otherwise not all are found
	var columns []*Column
	columns = append(columns, col)

	for r := col.head.down; r != &col.head; r = r.down {
		if r.value == value {
			// add r to solution
			m.sol = append(m.sol, r)

			for n := r.right; n != r; n = n.right {
				columns = append(columns, n.col)
			}

			// cover r
			for _, e := range columns {
				cover(e)
			}

			break
		}
	}
}

// (c) Mia Muessig
