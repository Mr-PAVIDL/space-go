package simple

import (
	"context"
	"space-go/internal/commander"
	"time"
)

type Strategy struct {
	visited map[string]struct{}
}

func New() *Strategy {
	return &Strategy{
		visited: map[string]struct{}{},
	}
}

func (strategy *Strategy) Next(ctx context.Context, state *commander.State) commander.Command {
	strategy.markVisited(state.Planet)

	if state.Planet == EdenName {
		// picking the planet to go to next
		var candidates []string

		for planet := range state.Universe.Planets {
			if !strategy.isVisited(planet) {
				candidates = append(candidates, planet)
			}
		}

		if len(candidates) == 0 {
			return commander.Idle(time.Second)
		}

		nearest := state.Universe.Nearest(state.Planet, candidates)

		return commander.GoTo(nearest)
	} else {
		return commander.Sequential(
			commander.Collect(),
			commander.GoTo(EdenName),
		)
	}
}

func (strategy *Strategy) markVisited(planet string) {
	strategy.visited[planet] = struct{}{}
}

func (strategy *Strategy) isVisited(planet string) bool {
	if planet == EarthName {
		return true
	}

	_, ok := strategy.visited[planet]
	return ok
}
