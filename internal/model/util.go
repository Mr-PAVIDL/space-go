package model

type Matrix [][]bool

//func (g Garbage) Matrix() Matrix {
//
//}
//
//func (g Garbage) Copy() Garbage {
//	g2 := make(Garbage, len(g))
//	for i := range g {
//		g2[i] = make([]int, len(g[i]))
//		copy(g2[i], g[i])
//	}
//	return g2
//}
//
//func (g Garbage) RotateCW() Garbage {
//	m, n := len(g), len(g[0])
//	rotated := make([][]int, n) // New matrix with flipped dimensions
//
//	for i := range rotated {
//		rotated[i] = make([]int, m)
//		for j := range rotated[i] {
//			rotated[i][j] = g[m-j-1][i]
//		}
//	}
//	return rotated
//}
//
//func (g Garbage) RotateCW() Garbage {
//	m, n := len(g), len(g[0])
//	rotated := make([][]int, n) // New matrix with flipped dimensions
//
//	for i := range rotated {
//		rotated[i] = make([]int, m)
//		for j := range rotated[i] {
//			rotated[i][j] = g[m-j-1][i]
//		}
//	}
//	return rotated
//}
