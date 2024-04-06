package commander

import (
	"context"
	"fmt"
	"log"
	"maps"
	"math"
	"space-go/internal/model"
	"strings"
	"time"
)

type Command interface {
	Execute(ctx context.Context, commander *Commander) error
}

type TravelCommand struct {
	Path []string
}

func Travel(planets ...string) TravelCommand {
	return TravelCommand{planets}
}

func (cmd TravelCommand) Execute(ctx context.Context, commander *Commander) error {
	if len(cmd.Path) == 0 {
		// empty path means we stay?
		return nil
	}

	response, err := commander.API.Travel(ctx, model.TravelRequest{
		Planets: cmd.Path,
	})

	if err != nil {
		return fmt.Errorf("failed to travel: %w", err)
	}

	for _, diff := range response.PlanetDiffs {
		commander.State.Universe.Graph[diff.From][diff.To] += diff.Fuel
	}
	commander.State.FuelUsed += response.FuelDiff
	commander.State.Planet = model.Planet{
		Garbage: response.PlanetGarbage,
		Name:    cmd.Path[len(cmd.Path)-1],
	}
	commander.State.Universe.Planets[commander.State.Planet.Name] = commander.State.Planet
	commander.State.Garbage = response.ShipGarbage
	if commander.State.Garbage == nil {
		commander.State.Garbage = map[string]model.Garbage{}
	}

	return nil
}

func (cmd TravelCommand) String() string {
	return formatCommand("travel", strings.Join(cmd.Path, " -> "))
}

type CollectCommand struct {
	proposal map[string]model.Garbage
	scouting bool
}

func CollectScouting() CollectCommand {
	return CollectCommand{proposal: nil, scouting: true}
}

func CollectWithProposal(proposal map[string]model.Garbage) CollectCommand {
	return CollectCommand{proposal: proposal}
}

func Collect() CollectCommand {
	return CollectCommand{proposal: nil}
}

var cntDeposits int = 0
var sumDeposits int = 0

func (cmd CollectCommand) Execute(ctx context.Context, commander *Commander) error {
	if len(commander.State.Planet.Garbage) == 0 {
		return nil
	}

	garbage := maps.Clone(commander.State.Garbage)
	for name, val := range commander.State.Planet.Garbage {
		garbage[name] = val
	}
	for name, val := range garbage {
		garbage[name], _, _ = val.Normalize()
	}
	minTiles := 0
	if len(commander.State.Garbage) == 0 {
		minTiles = int(math.Ceil(float64(commander.State.CapacityY*commander.State.CapacityX) * 0.30))
	} else {
		minTiles = int(math.Ceil(float64(commander.State.CapacityY*commander.State.CapacityX)*0.05)) + len(commander.State.Garbage)
	}

	newGarbage := commander.Packer.Pack(commander.State.CapacityX, commander.State.CapacityY, garbage, cmd.scouting, minTiles)
	if CountTiles(cmd.proposal) > CountTiles(newGarbage) {
		newGarbage = cmd.proposal
	}

	if len(commander.State.Garbage) > len(newGarbage) {
		log.Println("didn't improve: ", len(commander.State.Garbage),
			"->", len(newGarbage), commander.State.Garbage, newGarbage)
		//return nil
	}
	//commander.State.Garbage = newGarbage
	if len(newGarbage) != 0 {
		sumDeposits += CountTiles(newGarbage)
		cntDeposits++
		log.Println("deposit stats: ", float64(sumDeposits)/float64(cntDeposits), "avg deposit size / ",
			commander.State.CapacityY*commander.State.CapacityX)

		response, err := commander.API.CollectGarbage(ctx, model.CollectRequest{Garbage: newGarbage})

		if err != nil {
			if strings.Contains(err.Error(), "no garbage on this planet") {
				commander.State.Planet.Garbage = map[string]model.Garbage{}
				commander.State.Universe.Planets[commander.State.Planet.Name] = commander.State.Planet
			}

			return err
		}
		planetGarbage := map[string]model.Garbage{}
		for _, id := range response.Leaved {
			planetGarbage[id] = garbage[id]
		}
		commander.State.Planet.Garbage = planetGarbage
		commander.State.Universe.Planets[commander.State.Planet.Name] = commander.State.Planet
		commander.State.Garbage = response.Garbage
		if commander.State.Garbage == nil {
			commander.State.Garbage = make(map[string]model.Garbage)
		}
	}

	return nil
}

func (cmd CollectCommand) String() string {
	return formatCommand("collect", "collecting garbage...")
}

type GotoCommand struct {
	Destination string
}

func GoTo(planet string) GotoCommand {
	return GotoCommand{
		Destination: planet,
	}
}

func (cmd GotoCommand) Execute(ctx context.Context, commander *Commander) error {
	from := commander.State.Planet.Name

	path := commander.State.Universe.ShortestPath(from, cmd.Destination)

	if len(path) == 0 {
		return fmt.Errorf("no path from %s to %s", from, cmd.Destination)
	}
	if len(path) == 1 {
		return nil
	}

	return Travel(path[1:]...).Execute(ctx, commander)
}

func (cmd GotoCommand) String() string {
	return formatCommand("goto", cmd.Destination)
}

type IdleCommand struct {
	Duration time.Duration
}

func Idle(duration time.Duration) IdleCommand {
	return IdleCommand{duration}
}

func (cmd IdleCommand) Execute(ctx context.Context, commander *Commander) error {
	time.Sleep(cmd.Duration)

	return nil
}

func (cmd IdleCommand) String() string {
	return formatCommand("idle", cmd.Duration.String())
}

type SequentialCommand []Command

func Sequential(commands ...Command) SequentialCommand {
	return commands
}

func (cmd SequentialCommand) Execute(ctx context.Context, commander *Commander) error {
	for _, subCmd := range cmd {
		if err := commander.Execute(ctx, subCmd); err != nil {
			return err
		}
	}

	return nil
}

func (cmd SequentialCommand) String() string {
	return formatCommand("sequential", fmt.Sprintf("%d", len(cmd)))
}

func formatCommand(name string, message string) string {
	return fmt.Sprintf("[%s]: %s", name, message)
}
