package commander

import (
	"reflect"
	"space-go/internal/model"
	"testing"
)

func TestUniverse_Nearest(t *testing.T) {
	cases := []struct {
		name           string
		universe       Universe
		from           string
		to             []string
		expectedPlanet string
	}{
		{
			name: "Direct path to target",
			universe: Universe{
				Planets: map[string]model.Planet{
					"Earth": {},
					"Mars":  {},
				},
				Graph: map[string]map[string]int{
					"Earth": {"Mars": 10},
				},
			},
			from:           "Earth",
			to:             []string{"Mars"},
			expectedPlanet: "Mars",
		},
		{
			name: "Multiple targets, choose nearest",
			universe: Universe{
				Planets: map[string]model.Planet{
					"Earth": {},
					"Mars":  {},
					"Venus": {},
				},
				Graph: map[string]map[string]int{
					"Earth": {"Mars": 10, "Venus": 5},
				},
			},
			from:           "Earth",
			to:             []string{"Mars", "Venus"},
			expectedPlanet: "Venus",
		},
		{
			name: "No path to target",
			universe: Universe{
				Planets: map[string]model.Planet{
					"Earth": {},
					"Mars":  {},
				},
				Graph: map[string]map[string]int{
					"Earth": {},
				},
			},
			from:           "Earth",
			to:             []string{"Mars"},
			expectedPlanet: "",
		},
		{
			name: "Difficult",
			universe: Universe{
				Planets: map[string]model.Planet{
					"Earth": {}, "Alpha": {}, "Beta": {}, "Gamma": {}, "Delta": {}, "Eden": {},
				},
				Graph: map[string]map[string]int{
					"Earth": {"Alpha": 5, "Beta": 10},
					"Alpha": {"Gamma": 15},
					"Beta":  {"Gamma": 20, "Delta": 2},
					"Gamma": {"Eden": 30},
					"Delta": {"Eden": 25},
				},
			},
			from:           "Earth",
			to:             []string{"Gamma", "Eden"},
			expectedPlanet: "Gamma",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nearest := tc.universe.Nearest(tc.from, tc.to)
			if nearest != tc.expectedPlanet {
				t.Errorf("Expected %s, got %s", tc.expectedPlanet, nearest)
			}
		})
	}
}

//func (universe *Universe) Nearest(from string, to []string) string {
//	// Initialize distances map
//	distances := make(map[string]int)
//	for planet := range universe.Planets {
//		if planet == from {
//			distances[planet] = 0
//		} else {
//			distances[planet] = math.MaxInt32
//		}
//	}
//
//	// Dijkstra's algorithm
//	visited := make(map[string]bool)
//	for len(visited) < len(universe.Planets) {
//		// Find the unvisited planet with the smallest distance
//		minDistance := math.MaxInt32
//		minPlanet := ""
//		for planet, distance := range distances {
//			if !visited[planet] && distance <= minDistance {
//				minDistance = distance
//				minPlanet = planet
//			}
//		}
//
//		// Mark the planet as visited
//		visited[minPlanet] = true
//
//		// Update the distances to the neighboring planets
//		for neighbor, distance := range universe.Graph[minPlanet] {
//			if newDistance := distances[minPlanet] + distance; newDistance < distances[neighbor] {
//				distances[neighbor] = newDistance
//			}
//		}
//	}
//
//	// Find the nearest planet from the 'to' list
//	minDistance := math.MaxInt32
//	nearestPlanet := ""
//	for _, planet := range to {
//		if distances[planet] < minDistance {
//			minDistance = distances[planet]
//			nearestPlanet = planet
//		}
//	}
//
//	return nearestPlanet
//}

func TestUniverse_ShortestPath(t *testing.T) {
	tests := []struct {
		name     string
		universe Universe
		from     string
		to       string
		wantPath []string
	}{
		{
			name: "Direct Path",
			universe: *NewUniverse([]model.Transition{
				{FromPlanet: "Earth", ToPlanet: "Mars", FuelCost: 10},
			}),
			from:     "Earth",
			to:       "Mars",
			wantPath: []string{"Earth", "Mars"},
		},
		{
			name: "Indirect Path",
			universe: *NewUniverse([]model.Transition{
				{FromPlanet: "Earth", ToPlanet: "Alpha", FuelCost: 5},
				{FromPlanet: "Alpha", ToPlanet: "Beta", FuelCost: 10},
				{FromPlanet: "Beta", ToPlanet: "Mars", FuelCost: 15},
			}),
			from:     "Earth",
			to:       "Mars",
			wantPath: []string{"Earth", "Alpha", "Beta", "Mars"},
		},
		{
			name: "No Path Exists",
			universe: *NewUniverse([]model.Transition{
				{FromPlanet: "Earth", ToPlanet: "Alpha", FuelCost: 5},
			}),
			from:     "Earth",
			to:       "Mars",
			wantPath: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath := tt.universe.ShortestPath(tt.from, tt.to)
			if !reflect.DeepEqual(gotPath, tt.wantPath) {
				t.Errorf("ShortestPath() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
