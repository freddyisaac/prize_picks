package das

const (
	Herbivore     = "herbivore"
	HerbivoreCode = "H"
	Carnivore     = "carnivore"
	CarnivoreCode = "C"
	StatusDown    = "DOWN"
	StatusActive  = "ACTIVE"

	CageCapacity = 20
)

type Dinosaur struct {
	ID      uint   `json:"id"`
	Species string `json:"species"`
	Name    string `json:"name"`
	Diet    string `json:"diet"`
	Cage    uint   `json:"cage"`
}

type Cage struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Capacity int    `json:"capacity"`
	Count    int    `json:"count"`
	Kind     string `json:"kind"`
}

func ValidStatus(status string) bool {
	switch status {
	case StatusDown, StatusActive:
		return true
	default:
		return false
	}
	return false
}
