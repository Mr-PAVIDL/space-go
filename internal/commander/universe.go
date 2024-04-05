package commander

import "space-go/internal/model"

type Universe struct {
	Planets map[string]model.Planet
	Graph   map[string]map[string]int
}

func (universe *Universe) Nearest(from string, to []string) string {
	// TODO
	return ""
}
