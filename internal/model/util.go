package model

type Matrix [][]bool

// func (g Garbage) Matrix() Matrix {
//
// }
func (g Garbage) Copy() Garbage {
	g2 := make(Garbage, len(g))
	for i := range g {
		g2[i] = Cell{g[i][0], g[i][1]}
	}
	return g2
}

//

func (g Garbage) Normalize() Garbage {
	minX, minY := g[0][0], g[0][1]
	for _, cell := range g {
		minX, minY = min(minX, cell[0]), min(minY, cell[1])
	}
	g2 := g.Copy()
	for i, _ := range g {
		g2[i] = Cell{g2[i][0] - minX, g2[i][1] - minY}
	}
	return g2
}

func (mat Matrix) RotateCW() Matrix {
	m, n := len(mat), len(mat[0])
	rotated := make([][]bool, n) // New matrix with flipped dimensions

	for i := range rotated {
		rotated[i] = make([]bool, m)
		for j := range rotated[i] {
			rotated[i][j] = mat[m-j-1][i]
		}
	}
	return rotated
}

func (mat Matrix) RotateCCW() Matrix {
	m, n := len(mat), len(mat[0])
	rotated := make([][]bool, n) // New matrix with flipped dimensions

	for i := range rotated {
		rotated[i] = make([]bool, m)
		for j := range rotated[i] {
			rotated[i][j] = mat[m-j-1][i]
		}
	}
	return rotated
}
