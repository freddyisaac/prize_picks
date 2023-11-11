package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"dinocage/das"
)

type Species struct {
	Name string `json:"name"`
	Diet string `json:"diet"`
}

func ReadSpecies(name string) (*GenMap[string, string], error) {
	speciesMap := &GenMap[string, string]{}
	var species []Species
	b, err := os.ReadFile(name)
	if err != nil {
		return speciesMap, err
	}
	err = json.Unmarshal(b, &species)
	for _, s := range species {
		diet := strings.ToUpper(s.Diet)
		if diet != das.HerbivoreCode && diet != das.CarnivoreCode {
			// invalid code given skip
			continue
		}
		speciesMap.Store(strings.ToLower(s.Name), s.Diet)
	}
	speciesMap.Range(func(name, diet string) bool {
		log.Printf("known species: %s diet: %s", name, diet)
		return true
	})
	return speciesMap, err
}
