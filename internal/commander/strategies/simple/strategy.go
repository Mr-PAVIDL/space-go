package simple

import (
	"context"
	"github.com/samber/lo"
	"log"
	"maps"
	"math"
	"os"
	"space-go/internal/commander"
	"space-go/internal/commander/strategies"
	"space-go/internal/model"
)

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func MakeBigGarbage(size int) map[string]model.Garbage {
	g := make(model.Garbage, 0)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			g = append(g, model.Cell{i, j})
		}
	}
	return map[string]model.Garbage{"big_garbage": g}
}

func (strategy *Strategy) Next(ctx context.Context, state *commander.State) commander.Command {
	if state.Planet.Name == strategies.EdenName || len(state.Garbage) == 0 {
		// picking the planet to go to next
		var candidates []string

		minTiles := int(math.Ceil(float64(state.CapacityY*state.CapacityX) * 0.30))
		for _, planet := range state.Universe.Planets {
			// if garbage is nil, the planet is expected to be in unknown state
			// if garbage is an empty map, then the planet is emptied
			if planet.Garbage == nil || len(planet.Garbage) > 0 {
				if planet.Garbage != nil && commander.CountTiles(planet.Garbage) < minTiles {
					continue
				}
				candidates = append(candidates, planet.Name)
			}
		}

		if len(candidates) == 0 {
			os.Exit(0)
		}

		nearest := lo.Sample(candidates) //state.Universe.Nearest(state.Planet.Name, candidates)

		return commander.Sequential(
			commander.GoTo(nearest),
			commander.Collect(),
		)
	} else {
		packer := commander.DumboPacker{}
		myGarbage := maps.Clone(state.Garbage)
		haveALotOfSpace := packer.Pack(state.CapacityX, state.CapacityY, model.PileGarbage(myGarbage, MakeBigGarbage(4)))
		if len(haveALotOfSpace) > len(myGarbage) {
			// find uncovered planet because we are guaranteed to find something fitting there
			var candidates []string

			for _, planet := range state.Universe.Planets {
				if planet.Garbage == nil {
					candidates = append(candidates, planet.Name)
				}
			}
			if len(candidates) != 0 {
				log.Println("OPTIMIZATION: have a lot of empty space, will uncover nearest planet")
				nearest := state.Universe.Nearest(state.Planet.Name, candidates)
				return commander.Sequential(
					commander.GoTo(nearest),
					commander.Collect(),
				)
			} else {
				candidates = []string{}
				for _, planet := range state.Universe.Planets {
					if planet.Garbage != nil && len(planet.Garbage) > 0 {
						candidates = append(candidates, planet.Name)
					}
				}
				if len(candidates) != 0 {
					log.Println("OPTIMIZATION: have a lot of empty space, will grab something from nearby planet")
					nearest := state.Universe.Nearest(state.Planet.Name, candidates)
					return commander.Sequential(
						commander.GoTo(nearest),
						commander.Collect(),
					)
				}
			}
		}

		path := state.Universe.ShortestPath(state.Planet.Name, strategies.EdenName)
		for _, planet := range path[1 : len(path)-1] {
			garbage := state.Universe.Planets[planet].Garbage
			if garbage != nil && len(garbage) > 0 {
				garbage := maps.Clone(garbage)
				for name, val := range state.Garbage {
					garbage[name] = val
				}
				for name, val := range garbage {
					garbage[name] = val.Normalize()
				}

				packed := packer.Pack(state.CapacityX, state.CapacityY, garbage)
				was, now := commander.CountTiles(state.Garbage), commander.CountTiles(packed)
				minimalCells := int(math.Ceil(float64(state.CapacityY*state.CapacityX) * 0.05))

				// TODO: verify...
				if now == state.CapacityX*state.CapacityY ||
					now >= minimalCells+was {
					log.Println("OPTIMIZATION: ", "on route from ", state.Planet.Name, " to eden we can visit ",
						planet)
					return commander.Sequential(
						commander.GoTo(planet),
						commander.Collect(),
					)
				}

			}
		}

		return commander.Sequential(
			commander.GoTo(strategies.EdenName),
		)
	}
}
