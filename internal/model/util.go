package model

import "sync"

type Matrix [][]int

var pools = map[int]map[int]*sync.Pool{}

func allocateMatrix(width, height int) Matrix {
	subPools, ok := pools[width]
	if !ok {
		subPools = make(map[int]*sync.Pool)
		pools[width] = subPools
	}

	pool, ok := subPools[height]
	if !ok {
		pool = &sync.Pool{
			New: func() any {
				mat := make(Matrix, height)
				for i := 0; i < height; i++ {
					mat[i] = make([]int, width)
				}
				return mat
			},
		}
		subPools[height] = pool
	}

	return pool.Get().(Matrix)
}

func deallocateMatrix(mat Matrix) {
	if mat == nil {
		return
	}
	h := len(mat)
	if h == 0 {
		return
	}

	w := len(mat[0])
	subPools, ok := pools[w]
	if !ok {
		return
	}

	pool, ok := subPools[h]
	if !ok {
		return
	}

	pool.Put(mat)
}

func EmptyMatrix(width, height int) Matrix {
	//mat := make(Matrix, height)
	//for i := 0; i < height; i++ {
	//	mat[i] = make([]int, width)
	//}
	//return mat

	return allocateMatrix(width, height)
}

func (m Matrix) Close() {
	deallocateMatrix(m)
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

func (g Garbage) Normalize() (Garbage, int, int) {
	minX, minY := g[0][0], g[0][1]
	for _, cell := range g {
		minX, minY = min(minX, cell[0]), min(minY, cell[1])
	}
	w, h := 0, 0
	g2 := g.Copy()
	for i, _ := range g {
		g2[i] = Cell{g2[i][0] - minX, g2[i][1] - minY}

		w = max(w, g2[i][0])
		h = max(h, g2[i][1])
	}
	return g2, w + 1, h + 1
}

func (mat Matrix) RotateCW() Matrix {
	m, n := len(mat), len(mat[0])
	rotated := make([][]int, n) // New matrix with flipped dimensions

	for i := range rotated {
		rotated[i] = make([]int, m)
		for j := range rotated[i] {
			rotated[i][j] = mat[m-j-1][i]
		}
	}
	return rotated
}

func (mat Matrix) RotateCCW() Matrix {
	m, n := len(mat), len(mat[0])
	rotated := make([][]int, n) // New matrix with flipped dimensions

	for i := range rotated {
		rotated[i] = make([]int, m)
		for j := range rotated[i] {
			rotated[i][j] = mat[m-j-1][i]
		}
	}
	return rotated
}
