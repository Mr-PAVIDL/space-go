package simple

import (
	"context"
	"github.com/samber/lo"
	"os"
	"space-go/internal/commander"
	"space-go/internal/commander/strategies"
)

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func (strategy *Strategy) Next(ctx context.Context, state *commander.State) commander.Command {
	if state.Planet.Name == strategies.EdenName || len(state.Garbage) == 0 {
		// picking the planet to go to next
		var candidates []string

		for _, planet := range state.Universe.Planets {
			// if garbage is nil, the planet is expected to be in unknown state
			// if garbage is an empty map, then the planet is emptied
			if planet.Garbage == nil || len(planet.Garbage) > 0 {
				candidates = append(candidates, planet.Name)
			}
		}

		if len(candidates) == 0 {
			os.Exit(0)
		}

		//nearest := state.Universe.Farthest(state.Planet.Name, candidates)

		return commander.Sequential(
			commander.GoTo(lo.Sample(candidates)),
			commander.Collect(),
		)
	} else {
		//packer := commander.DumboPacker{}
		//path := state.Universe.ShortestPath(state.Planet.Name, strategies.EdenName)
		//for _, planet := range path[1 : len(path)-1] {
		//	garbage := state.Universe.Planets[planet].Garbage
		//	if garbage != nil && len(garbage) > 0 {
		//		garbage := maps.Clone(garbage)
		//		for name, val := range state.Garbage {
		//			garbage[name] = val
		//		}
		//		for name, val := range garbage {
		//			garbage[name] = val.Normalize()
		//		}
		//
		//		packed := packer.Pack(state.CapacityX, state.CapacityY, garbage)
		//		was, now := commander.CalcTiles(state.Garbage), commander.CalcTiles(packed)
		//		wasPercent := 100 * was / state.CapacityX * state.CapacityY
		//		nowPercent := 100 * now / state.CapacityX * state.CapacityY
		//
		//		// TODO: verify...
		//		if now == state.CapacityX*state.CapacityY ||
		//			nowPercent >= wasPercent+5 {
		//			log.Println("OPTIMIZATION: ", "on route from ", state.Planet.Name, " to eden we can visit ",
		//				planet)
		//			return commander.Sequential(
		//				commander.GoTo(planet),
		//				commander.Collect(),
		//			)
		//		}
		//
		//	}
		//}

		return commander.Sequential(
			commander.GoTo(strategies.EdenName),
		)
	}
}
