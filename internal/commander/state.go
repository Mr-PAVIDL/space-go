package commander

type State struct {
	Universe         *Universe
	Planet           string
	GarbageCollected int
	FuelUsed         int
}
