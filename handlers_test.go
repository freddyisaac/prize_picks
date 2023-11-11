package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "dinocage/das"
	"dinocage/mocks"

	gomock "github.com/golang/mock/gomock"
)

func Closer(da DataAccessProvider) {
	da.Close()
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockDataAccessProvider(ctrl)

	m.EXPECT().Close()

	Closer(m)

}

func TestHealthCheckHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDap := mocks.NewMockDataAccessProvider(ctrl)

	r, err := http.NewRequest("GET", "http://localhost:8000/healthcheck", nil)
	if err != nil {
		t.Errorf("NewRequest failed with %v", err)
	}
	w := httptest.NewRecorder()

	ah := &AppHandlers{dap: mockDap}

	ah.healthcheck(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("healthcheck failed with status %v", resp.StatusCode)
	}
}

func TestAddDino(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDap := mocks.NewMockDataAccessProvider(ctrl)

	newDino := &Dinosaur{
		Species: "Tyrannosaurus",
		Name:    "Barnie",
		Diet:    "C",
		Cage:    1,
	}

	serialDino, err := json.Marshal(newDino)

	if err != nil {
		t.Errorf("TestAddDino marshal failed with %v", err)
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8000/dino/add", bytes.NewBuffer(serialDino))
	if err != nil {
		t.Errorf("NewRequest failed with %v", err)
	}
	w := httptest.NewRecorder()

	speciesMap := &GenMap[string, string]{}
	speciesMap.Store("tyrannosaurus", "C")
	speciesMap.Range(func(species, diet string) bool {
		t.Logf("%s : %s", species, diet)
		return true
	})
	ah := &AppHandlers{dap: mockDap, speciesMap: speciesMap}

	// prepare the data access call

	var dino Dinosaur
	mockDap.EXPECT().AddDinosaur(gomock.Any(), gomock.AssignableToTypeOf(dino)).DoAndReturn(
		func(v interface{}, arg Dinosaur) error {
			dino = arg
			t.Logf("TestAddDino::.AddDinosaur received Dino : %+v", dino)
			return nil
		},
	)

	ah.AddDinosaur(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("TestAddDino did not return success but gave %v", resp.StatusCode)
	}

}

func TestErrorAddDino(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDap := mocks.NewMockDataAccessProvider(ctrl)

	newDino := &Dinosaur{
		Species: "Tyrannosaurus",
		Name:    "Barnie",
		Diet:    "C",
		Cage:    1,
	}

	serialDino, err := json.Marshal(newDino)

	if err != nil {
		t.Errorf("TestAddDino marshal failed with %v", err)
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8000/dino/add", bytes.NewBuffer(serialDino))
	if err != nil {
		t.Errorf("NewRequest failed with %v", err)
	}
	w := httptest.NewRecorder()

	speciesMap := &GenMap[string, string]{}
	speciesMap.Store("tyrannosaurus", "C")
	ah := &AppHandlers{dap: mockDap, speciesMap: speciesMap}

	// prepare the data access call

	var dino Dinosaur
	mockDap.EXPECT().AddDinosaur(gomock.Any(), gomock.AssignableToTypeOf(dino)).DoAndReturn(
		func(v interface{}, arg Dinosaur) error {
			dino = arg
			t.Logf("TestAddDino:AddDinosaur received Dino : %+v", dino)
			return fmt.Errorf("Test Error")
		},
	)

	ah.AddDinosaur(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("TestAddDino:AddDinosaur did not return %v but gave %v", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}
