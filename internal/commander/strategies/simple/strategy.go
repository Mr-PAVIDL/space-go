package simple

import (
	"context"
	"space-go/internal/commander"
)

type Strategy struct {
	visited map[string]struct{}
}

func New() *Strategy {
	return &Strategy{}
}

func (strategy *Strategy) Next(ctx context.Context, state *commander.State) commander.Command {
	strategy.visited[state.Planet] = struct{}{}

	if state.Planet == EdenName {
		// picking the planet to go to next
		candidates := []string{}

		for planet := range state.Universe.Planets {
			if strategy.isVisited(planet) {
				candidates = append(candidates, planet)
			}
		}

		nearest := state.Universe.Nearest()
	} else {

	}
}

func (strategy *Strategy) isVisited(planet string) bool {
	if planet == EarthName {
		return true
	}

}
