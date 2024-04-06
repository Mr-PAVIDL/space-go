package simple

import (
	"context"
	"github.com/samber/lo"
	"log"
	"maps"
	"math"
	"os"
	"space-go/internal/commander"
	packer2 "space-go/internal/commander/packer"
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

const ScoutUntilUncoveredFraction = 1.0

func (strategy *Strategy) Next(ctx context.Context, state *commander.State) commander.Command {
	//uncoveredCount := 0
	//for _, planet := range state.Universe.Planets {
	//	if planet.Garbage == nil {
	//		uncoveredCount++
	//	}
	//}
	//if float64(uncoveredCount) < float64(len(state.Universe.Planets))*ScoutUntilUncoveredFraction {
	if scout := strategy.scout(state); scout != nil {
		return scout
	}
	//}

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
		packer := packer2.DuploPacker{}
		if strategy.hasGuaranteedSpace(state) {
			{
				candidates := []string{}
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
					garbage[name], _, _ = val.Normalize()
				}

				packed := packer.Pack(state.CapacityX, state.CapacityY, garbage, true, 0)
				was, now := commander.CountTiles(state.Garbage), commander.CountTiles(packed)
				minimalCells := int(math.Ceil(float64(state.CapacityY*state.CapacityX) * 0.05))

				// TODO: verify...
				if now == state.CapacityX*state.CapacityY ||
					len(garbage) == len(packed) ||
					now >= minimalCells+was {
					log.Println("OPTIMIZATION: ", "on route from ", state.Planet.Name, " to eden we can visit ",
						planet)
					return commander.Sequential(
						commander.GoTo(planet),
						commander.CollectWithProposal(packed),
					)
				}

			}
		}

		return commander.Sequential(
			commander.GoTo(strategies.EdenName),
		)
	}
}

func (strategy *Strategy) hasGuaranteedSpace(state *commander.State) bool {
	packer := packer2.DuploPacker{}
	myGarbage := maps.Clone(state.Garbage)
	haveALotOfSpace := packer.Pack(state.CapacityX, state.CapacityY, model.PileGarbage(myGarbage, MakeBigGarbage(4)), true, 0)
	return len(haveALotOfSpace) > len(myGarbage)
}

func (strategy *Strategy) scout(state *commander.State) commander.Command {
	if !strategy.hasGuaranteedSpace(state) {
		return nil
	}

	// find not-uncovered planet because we are guaranteed to find something fitting there
	var candidates []string
	for _, planet := range state.Universe.Planets {
		if planet.Garbage == nil {
			candidates = append(candidates, planet.Name)
		}
	}
	if len(candidates) != 0 {
		log.Println("OPTIMIZATION: have a lot of empty space, will scout")
		nearest := state.Universe.Nearest(state.Planet.Name, candidates)
		return commander.Sequential(
			commander.GoTo(nearest),
			commander.Collect(),
		)
	}
	return nil
}
