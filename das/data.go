package das

import (
	"errors"
)

type Dinosaur struct {
	ID      uint   `json:"id"`
	Species string `json:"species"`
	Name    string `json:"name"`
	Diet    string `json:"diet"`
	Cage    uint   `json:"cage"`
}

var (
	ErrSpeciesBadName = errors.New("species name must not be empty or missing")
	ErrSpeciesBadDiet = errors.New("species diet must be herbivore or carnivore")
)

const (
	Herbivore     = "herbivore"
	HerbivoreCode = "H"
	Carnivore     = "carnivore"
	CarnivoreCode = "C"
)

type Cage struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Capacity int    `json:"capaciy"`
	Count    int    `json:"count"`
	Kind     string `json:"kind"`
}
