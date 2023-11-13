package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	. "dinocage/das"

	"github.com/gorilla/mux"
)

var (
	MesgOK = []byte(`{"msg":"ok"}`)
)

// Core application data structure
type AppHandlers struct {
	dap        DataAccessProvider
	speciesMap *GenMap[string, string]
}

// check species against in memory species list
func (ah AppHandlers) CheckSpecies(speciesName string) bool {
	_, ok := ah.speciesMap.Load(speciesName)
	return ok
}

// add species to in memory list
func (ah AppHandlers) NewSpecies(name, diet string) bool {
	diet = strings.ToUpper(diet)
	if diet != HerbivoreCode && diet != CarnivoreCode {
		return false
	}
	ah.speciesMap.Store(strings.ToLower(name), diet)
	return true
}

// write message back to client utility
func WriteMsg(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Printf("unable to write back to client")
	}
}

// write http success utiliry
func WriteOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(MesgOK)
	if err != nil {
		log.Printf("unable to write back to client")
	}
}

// app server health responder
func (ah AppHandlers) healthcheck(w http.ResponseWriter, r *http.Request) {
	WriteOk(w)
}

// persist a new dinosaur to database - assigning to an open or new cage
func (ah AppHandlers) AddDinosaur(w http.ResponseWriter, r *http.Request) {
	// read payload
	dino := Dinosaur{}
	err := json.NewDecoder(r.Body).Decode(&dino)
	if err != nil {
		WriteMsg(w, http.StatusBadRequest, "bad payload "+err.Error())
		return
	}
	if !ah.CheckSpecies(strings.ToLower(dino.Species)) {
		WriteMsg(w, http.StatusBadRequest, "unknown species "+dino.Species)
		return
	}
	defer r.Body.Close()
	err = ah.dap.AddDinosaur(r.Context(), dino)
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	WriteOk(w)
}

// create a new cage of the given dietary type
func (ah AppHandlers) AddCage(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	kind := vars["diet"]
	paramCap := r.URL.Query().Get("cap")
	cap := CageCapacity
	if len(paramCap) != 0 {
		cap, err = strconv.Atoi(paramCap)
		if err != nil {
			WriteMsg(w, http.StatusBadRequest, "bad capacity parameter must be an integer")
			return
		}
		if cap < 1 {
			WriteMsg(w, http.StatusBadRequest, "bad capacity parameter value must be > 0")
			return
		}
	}
	var id int
	switch kind {
	case "H":
		id, err = ah.dap.NewCage(r.Context(), cap, HerbivoreCode)
	case "C":
		id, err = ah.dap.NewCage(r.Context(), cap, CarnivoreCode)
	default:
		WriteMsg(w, http.StatusBadRequest, "diet must be H (herbivore) V (carnivore)")
		return
	}
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// write back using anonymous id struct
	v := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}
	b, _ := json.Marshal(v)
	WriteMsg(w, http.StatusOK, string(b))
}

// list cages handler
func (ah AppHandlers) GetCages(w http.ResponseWriter, r *http.Request) {
	var err error
	var cages []Cage
	activeOpt := r.URL.Query().Get("status")

	// only apply filter if it exists and is valid
	if len(activeOpt) != 0 && ValidStatus(activeOpt) {
		cages, err = ah.dap.GetCages(r.Context(), activeOpt)
	} else {
		cages, err = ah.dap.GetCages(r.Context())
	}
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, "database error : "+err.Error())
		return
	}
	b, err := json.Marshal(cages)
	if err != nil {
		WriteMsg(w, http.StatusInternalServerError, "unable to parse "+err.Error())
		return
	}
	WriteMsg(w, http.StatusOK, string(b))
}

// list handler for dinosaurs in a given cage
func (ah AppHandlers) GetCageDinosaurs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paramCageID := vars["cageid"]
	if len(paramCageID) == 0 {
		WriteMsg(w, http.StatusBadRequest, "no valid cage id given")
		return
	}
	cageID, err := strconv.Atoi(paramCageID)
	if err != nil {
		WriteMsg(w, http.StatusBadRequest, fmt.Sprintf("invalid cage id given %s is not an integer", paramCageID))
		return
	}
	dinoList, err := ah.dap.GetDinosaursForCage(r.Context(), cageID)
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, "database error : "+err.Error())
		return
	}
	b, _ := json.Marshal(dinoList)
	WriteMsg(w, http.StatusOK, string(b))
}

// add new species handler
func (ah AppHandlers) AddSpecies(w http.ResponseWriter, r *http.Request) {
	var species Species
	err := json.NewDecoder(r.Body).Decode(&species)
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, "bad request for add species")
		return
	}
	log.Printf("adding %s %s", species.Name, species.Diet)
	ah.NewSpecies(species.Name, species.Diet)
	WriteOk(w)
}

// list species handler
func (ah AppHandlers) ListSpecies(w http.ResponseWriter, r *http.Request) {
	var species []Species
	ah.speciesMap.Range(func(name, diet string) bool {
		species = append(species, Species{Name: name, Diet: diet})
		return true
	})
	b, err := json.Marshal(species)
	if err != nil {
		WriteMsg(w, http.StatusInternalServerError, "marshal error : "+err.Error())
		return
	}
	WriteMsg(w, http.StatusOK, string(b))
}

// set the status of a specified cage handler
func (ah AppHandlers) SetCageStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]
	if len(status) == 0 || !ValidStatus(status) {
		WriteMsg(w, http.StatusBadRequest, "invalid status")
		return
	}
	paramCageID := vars["cageid"]
	if len(paramCageID) == 0 {
		WriteMsg(w, http.StatusBadRequest, "invalid cage status")
		return
	}
	cageID, err := strconv.Atoi(paramCageID)
	if err != nil {
		WriteMsg(w, http.StatusBadRequest, "cageid is not a valid value")
		return
	}
	err = ah.dap.SetCageStatus(r.Context(), cageID, status)
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, fmt.Sprintf("unable to set cage %d status to %s check that the cage is empty or powered down", cageID, status))
		return
	}
	WriteOk(w)
}

// list dinosaur handler
func (ah AppHandlers) GetDinosaurs(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")
	var err error
	var dinos []Dinosaur
	if len(species) != 0 && ah.CheckSpecies(species) {
		dinos, err = ah.dap.GetDinosaurs(r.Context(), species)
	} else {
		dinos, err = ah.dap.GetDinosaurs(r.Context())
	}
	if err != nil {
		WriteMsg(w, http.StatusUnprocessableEntity, "database error accessing dinosaurs : %v "+err.Error())
		return
	}
	buf, err := json.Marshal(dinos)
	if err != nil {
		WriteMsg(w, http.StatusInternalServerError, "unable to marshal response "+err.Error())
		return
	}
	WriteMsg(w, http.StatusOK, string(buf))
}

// place dinosaur in cage handler
func (ah AppHandlers) AddDinoToCage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paramCageID := vars["cageid"]
	if len(paramCageID) == 0 {
		WriteMsg(w, http.StatusBadRequest, "cage id not present")
		return
	}
	cageID, err := strconv.Atoi(paramCageID)
	if err != nil {
		WriteMsg(w, http.StatusBadRequest, "cage id not a valid value")
		return
	}
	var dino Dinosaur
	err = json.NewDecoder(r.Body).Decode(&dino)
	if err != nil {
		WriteMsg(w, http.StatusBadRequest, "bad dino data "+err.Error())
		return
	}

	defer r.Body.Close()
	err = ah.dap.PlaceDinosaurInCage(r.Context(), cageID, dino)
	if err != nil {
		WriteMsg(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteOk(w)
}

// create mux and start server
func StartServer(ctx context.Context, listenAddr string, appHandlers *AppHandlers) error {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", appHandlers.healthcheck).Methods("GET")
	r.HandleFunc("/dino/add", appHandlers.AddDinosaur).Methods("POST")
	r.HandleFunc("/dino/list", appHandlers.GetDinosaurs).Methods("GET")
	r.HandleFunc("/cages", appHandlers.GetCages).Methods("GET")
	r.HandleFunc("/cage/{diet}/add", appHandlers.AddCage).Methods("POST")
	r.HandleFunc("/cage/{cageid}/list_dinosaurs", appHandlers.GetCageDinosaurs).Methods("GET")
	r.HandleFunc("/cage/{cageid}/status/{status}", appHandlers.SetCageStatus).Methods("POST")
	r.HandleFunc("/cage/{cageid}/add_dino", appHandlers.AddDinoToCage).Methods("POST")
	r.HandleFunc("/species/add", appHandlers.AddSpecies).Methods("POST")
	r.HandleFunc("/species/list", appHandlers.ListSpecies).Methods("GET")

	server := http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	// prepare for shutdown initiated from context
	go func() {
		<-ctx.Done()
		// received context done
		// shutdown server
		server.Shutdown(context.Background())
	}()
	log.Printf("Starting serv on %s", listenAddr)
	err := server.ListenAndServe()
	return err
}
