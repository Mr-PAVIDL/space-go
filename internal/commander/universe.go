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

func (universe *Universe) Nearest(from string, to []string) string {
	dist := make(map[string]int)
	for planet := range universe.Planets {
		dist[planet] = math.MaxInt32
	}
	dist[from] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{name: from, fuel: 0})

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		if visited[current.name] {
			continue
		}
		visited[current.name] = true

		for _, target := range to {
			if current.name == target {
				return current.name
			}
		}

		for neighbor, fuelCost := range universe.Graph[current.name] {
			if !visited[neighbor] {
				newFuel := current.fuel + fuelCost
				if newFuel < dist[neighbor] {
					dist[neighbor] = newFuel
					heap.Push(&pq, &Item{name: neighbor, fuel: newFuel})
				}
			}
		}
	}

	return ""
}
