package commander

import "space-go/internal/model"

type State struct {
	Universe               *Universe
	Planet                 model.Planet
	Garbage                map[string]model.Garbage
	GarbageCellsCollected  int
	GarbagePiecesCollected int
	FuelUsed               int
}
