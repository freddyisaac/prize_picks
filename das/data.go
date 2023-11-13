package das

type Dinosaur struct {
	ID      uint   `json:"id"`
	Species string `json:"species"`
	Name    string `json:"name"`
	Diet    string `json:"diet"`
	Cage    uint   `json:"cage"`
}

const (
	Herbivore     = "herbivore"
	HerbivoreCode = "H"
	Carnivore     = "carnivore"
	CarnivoreCode = "C"
)

type Cage struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Capacity int    `json:"capacity"`
	Count    int    `json:"count"`
	Kind     string `json:"kind"`
}
