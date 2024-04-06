package model

type Matrix [][]bool

func EmptyMatrix(width, height int) Matrix {
	mat := make(Matrix, height)
	for i := 0; i < height; i++ {
		mat[i] = make([]bool, width)
	}
	return mat
}

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

func (g Garbage) Add(x, y int) Garbage {
	g2 := make(Garbage, len(g))
	for i := 0; i < len(g2); i++ {
		g2[i] = Cell{g[i][0] + x, g[i][1] + y}
	}
	return g2
}

func PileGarbage(gs ...map[string]Garbage) map[string]Garbage {
	result := make(map[string]Garbage)
	for _, g := range gs {
		for s, garbage := range g {
			result[s] = garbage
		}
	}
	return result
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
