package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	Error string `json:"error"`
}

type CollectRequest struct {
	Garbage map[string][][]int `json:"garbage"`
}

type TravelRequest struct {
	Planets []string `json:"planets"`
}

type CollectResponse struct {
	Garbage map[string][][]int `json:"garbage"`
	Leaved  []string           `json:"leaved"`
}

type PlanetDiff struct {
	From string `json:"from"`
	Fuel int    `json:"fuel"`
	To   string `json:"to"`
}

type Planet struct {
	Garbage map[string][][]int `json:"garbage"`
	Name    string             `json:"name"`
}

type Ship struct {
	CapacityX int                `json:"capacityX"`
	CapacityY int                `json:"capacityY"`
	FuelUsed  int                `json:"fuelUsed"`
	Garbage   map[string][][]int `json:"garbage"`
	Planet    Planet             `json:"planet"`
}

type Transition struct {
	FromPlanet string `json:"fromPlanet"`
	ToPlanet   string `json:"toPlanet"`
	FuelCost   int    `json:"fuelCost"`
}

type Universe []Transition

func (u *Universe) UnmarshalJSON(data []byte) error {
	var raw [][]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var transitions []Transition
	for _, item := range raw {
		if len(item) != 3 {
			return errors.New("universe transition must have exactly 3 elements")
		}

		var fromPlanet, toPlanet string
		var fuelCost int
		if err := json.Unmarshal(item[0], &fromPlanet); err != nil {
			return fmt.Errorf("unmarshalling fromPlanet: %w", err)
		}
		if err := json.Unmarshal(item[1], &toPlanet); err != nil {
			return fmt.Errorf("unmarshalling toPlanet: %w", err)
		}
		if err := json.Unmarshal(item[2], &fuelCost); err != nil {
			return fmt.Errorf("unmarshalling fuelCost: %w", err)
		}

		transitions = append(transitions, Transition{
			FromPlanet: fromPlanet,
			ToPlanet:   toPlanet,
			FuelCost:   fuelCost,
		})
	}

	*u = transitions
	return nil
}

type Player struct {
	Name       string   `json:"name"`
	RoundName  string   `json:"roundName"`
	RoundEndIn int      `json:"roundEndIn"`
	Ship       Ship     `json:"ship"`
	Universe   Universe `json:"universe"`
	Attempt    int      `json:"attempt"`
}

type Round struct {
	StartAt     string `json:"startAt"`
	EndAt       string `json:"endAt"`
	IsCurrent   bool   `json:"isCurrent"`
	Name        string `json:"name"`
	PlanetCount int    `json:"planetCount"`
}

type RoundList struct {
	Rounds []Round `json:"rounds"`
}

type TravelResponse struct {
	FuelDiff      int                `json:"fuelDiff"`
	PlanetDiffs   []PlanetDiff       `json:"planetDiffs"`
	PlanetGarbage map[string][][]int `json:"planetGarbage"`
	ShipGarbage   map[string][][]int `json:"shipGarbage"`
}

type AcceptedResponse struct {
	Success bool `json:"success"`
}
