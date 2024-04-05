package simple

import (
	"context"
	"space-go/internal/commander"
	"time"
)

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func (strategy *Strategy) Next(ctx context.Context, state *commander.State) commander.Command {
	if state.Planet.Name == EdenName {
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
			return commander.Idle(time.Second)
		}

		nearest := state.Universe.Nearest(state.Planet.Name, candidates)

		return commander.GoTo(nearest)
	} else {
		return commander.Sequential(
			commander.Collect(),
			commander.GoTo(EdenName),
		)
	}
}
