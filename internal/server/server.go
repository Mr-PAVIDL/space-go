package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"maps"
	"os"
	"sort"
	"space-go/internal/model"
)

const PathCostIncrement = 10

type Server struct {
	graph      map[string]map[string]int
	planets    map[string]model.Planet
	ship       model.Ship
	allGarbage map[string]model.Garbage
	//trunk     [][]string // keep in sync with ship.garbage!
	totalFuel int
}

func FromDump(planetsPath, graphPath string) Server {
	graph := map[string]map[string]int{}
	planets := map[string]model.Planet{}
	allGarbage := map[string]model.Garbage{}

	graphJson, err := os.ReadFile(graphPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(graphJson, &graph)
	if err != nil {
		panic(err)
	}
	planetsData := map[string]map[string]model.Garbage{}
	planetsJson, err := os.ReadFile(planetsPath)
	err = json.Unmarshal(planetsJson, &planetsData)
	for planet, garbages := range planetsData {
		planets[planet] = model.Planet{
			Garbage: garbages,
			Name:    planet,
		}
		for name, garbage := range garbages {
			allGarbage[name] = garbage
		}
	}

	return Server{
		graph:   graph,
		planets: planets,
		ship: model.Ship{
			CapacityX: 8,
			CapacityY: 11,
			FuelUsed:  0,
			Garbage:   map[string]model.Garbage{},
			Planet:    planets["Earth"],
		},
		totalFuel:  0,
		allGarbage: allGarbage,
	}
}

func (s *Server) Universe() model.UniverseResponse {
	uni := make(model.Universe, 0)
	for from, toMap := range s.graph {
		for to, cost := range toMap {
			uni = append(uni, model.Transition{
				FromPlanet: from,
				ToPlanet:   to,
				FuelCost:   cost,
			})
		}
	}

	return model.UniverseResponse{
		Name:       "test",
		RoundName:  "test",
		RoundEndIn: 3600,
		Ship:       s.ship,
		Universe:   uni,
		Attempt:    0,
	}
}

func (s *Server) Travel(request model.TravelRequest) (model.TravelResponse, error) {
	if len(request.Planets) == 0 {
		return model.TravelResponse{}, errors.New("at least one planet is required")
	}

	totalCost := 0
	prevPlanet := s.ship.Planet.Name
	var planetsDiffs []model.PlanetDiff
	for _, planet := range request.Planets {
		if _, ok := s.graph[prevPlanet]; !ok {
			panic("you are on non existing planet")
			//return model.TravelResponse{}, errors.New("unknown planet in request")
		}

		if cost, ok := s.graph[prevPlanet][planet]; !ok {
			return model.TravelResponse{}, fmt.Errorf("no path from %s to %s", prevPlanet, planet)
		} else {
			totalCost += cost
			planetsDiffs = append(planetsDiffs, model.PlanetDiff{
				From: prevPlanet,
				To:   planet,
				Fuel: cost + PathCostIncrement,
			})
			prevPlanet = planet
		}
	}

	return model.TravelResponse{
		FuelDiff:      totalCost,
		PlanetDiffs:   planetsDiffs,
		PlanetGarbage: s.planets[request.Planets[len(request.Planets)-1]].Garbage,
		ShipGarbage:   s.ship.Garbage,
	}, nil

}

func checkOverlaps(newTrunk [][]int, garbage model.Garbage) bool {
	for _, cell := range garbage {
		if newTrunk[cell[0]][cell[1]] != 0 {
			return true
		}
	}

	return false
}

func checkOutOfBounds(garbage model.Garbage, x, y int) bool {
	for _, cell := range garbage {
		if cell[0] < 0 || cell[0] >= y || cell[1] < 0 || cell[1] >= x {
			return true
		}
	}

	return false
}

func (s *Server) Collect(request model.CollectRequest) (model.CollectResponse, error) {
	if len(request.Garbage) == 0 {
		return model.CollectResponse{}, errors.New("at least one garbage is required")
	}

	newTrunk := make([][]int, s.ship.CapacityY)
	for i := range newTrunk {
		newTrunk[i] = make([]int, s.ship.CapacityX)
	}

	for name, garbage := range request.Garbage {
		onPlanet := s.ship.Planet.Garbage[name]
		onShip := s.ship.Garbage[name]

		var expected model.Garbage
		if onPlanet != nil {
			expected = onPlanet
		}
		if onShip != nil {
			expected = onShip
		}
		if expected == nil {
			return model.CollectResponse{}, fmt.Errorf("no garbage %s on planet %s or on ship", name, s.ship.Planet.Name)
		}
		if len(expected) != len(garbage) {
			return model.CollectResponse{}, fmt.Errorf("garbage %s has incorrect form", name)
		}

		if checkOverlaps(newTrunk, garbage) {
			return model.CollectResponse{}, fmt.Errorf("garbage %s overlaps with other garbage", name)
		}

		if checkOutOfBounds(garbage, s.ship.CapacityX, s.ship.CapacityY) {
			return model.CollectResponse{}, fmt.Errorf("garbage %s is out of bounds", name)
		}

		for _, cell := range garbage {
			newTrunk[cell[0]][cell[1]] = 1
		}

		// how to write fucking comparator?
		sort.Slice(expected, func(i, j int) bool {
			return expected[i][0] < expected[j][0] || (expected[i][1] == expected[j][1] && expected[i][1] < expected[j][1])
		})
		sort.Slice(garbage, func(i, j int) bool {
			return garbage[i][0] < garbage[j][0] || (garbage[i][1] == garbage[j][1] && garbage[i][1] < garbage[j][1])
		})
		for i := 0; i < len(garbage); i++ {
			if expected[i] != garbage[i] {
				return model.CollectResponse{}, fmt.Errorf("garbage %s has incorrect form", name)
			}
		}
	}

	leaved := maps.Clone(s.ship.Planet.Garbage)
	for name, garbage := range s.ship.Garbage {
		leaved[name] = garbage
	}
	for name, _ := range request.Garbage {
		delete(leaved, name)
	}
	s.ship.Garbage = request.Garbage
	s.ship.Planet.Garbage = leaved
	s.planets[s.ship.Planet.Name] = s.ship.Planet

	return model.CollectResponse{
		Garbage: request.Garbage,
		Leaved:  lo.Keys(leaved),
	}, nil
}
