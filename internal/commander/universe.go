package commander

import (
	"container/heap"
	"math"
	"space-go/internal/model"
)

type Item struct {
	name  string
	fuel  int
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].fuel < pq[j].fuel
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type Universe struct {
	Planets map[string]model.Planet
	Graph   map[string]map[string]int
}

func NewUniverse(transitions []model.Transition) *Universe {
	universe := &Universe{
		Planets: map[string]model.Planet{},
		Graph:   map[string]map[string]int{},
	}

	for _, transition := range transitions {
		universe.Planets[transition.ToPlanet] = model.Planet{
			Garbage: nil,
			Name:    transition.ToPlanet,
		}

		universe.Planets[transition.FromPlanet] = model.Planet{
			Garbage: nil,
			Name:    transition.FromPlanet,
		}

		if _, ok := universe.Graph[transition.FromPlanet]; !ok {
			universe.Graph[transition.FromPlanet] = map[string]int{}
		}

		if _, ok := universe.Graph[transition.ToPlanet]; !ok {
			universe.Graph[transition.ToPlanet] = map[string]int{}
		}

		universe.Graph[transition.FromPlanet][transition.ToPlanet] = transition.FuelCost
	}

	return universe
}

func (universe *Universe) Nearest(from string, to []string) string {
	dist, _ := universe.DijkstraWithPaths(from)

	nearest := ""
	minFuelCost := math.MaxInt32

	for _, target := range to {
		if fuelCost, exists := dist[target]; exists && fuelCost < minFuelCost {
			nearest = target
			minFuelCost = fuelCost
		}
	}

	return nearest
}

func (universe *Universe) DijkstraWithPaths(from string) (map[string]int, map[string]string) {
	dist := make(map[string]int)
	prev := make(map[string]string)
	for planet := range universe.Planets {
		dist[planet] = math.MaxInt32
		prev[planet] = ""
	}
	dist[from] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{name: from, fuel: 0})

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		visited[current.name] = true

		for neighbor, fuelCost := range universe.Graph[current.name] {
			if !visited[neighbor] {
				newFuel := current.fuel + fuelCost
				if newFuel < dist[neighbor] {
					dist[neighbor] = newFuel
					prev[neighbor] = current.name
					heap.Push(&pq, &Item{name: neighbor, fuel: newFuel})
				}
			}
		}
	}

	return dist, prev
}

func (universe *Universe) ShortestPath(from, to string) []string {
	_, prev := universe.DijkstraWithPaths(from)
	path := make([]string, 0)

	for at := to; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
	}

	if len(path) == 0 || path[0] != from {
		return nil
	}

	return path
}
