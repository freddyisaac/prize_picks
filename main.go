package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"dinocage/das"
)

const (
	ENV_DB_HOST      = "ENV_DB_HOST"
	ENV_DB_PORT      = "ENV_DB_PORT"
	ENV_DB_NAME      = "ENV_DB_NAME"
	ENV_DB_USR       = "ENV_DB_USR"
	ENV_DB_PASS      = "ENV_DB_PWD"
	ENV_SVR_ENDPOINT = "ENV_SVR_ENDPOINT"
	DefaultEndpoint  = ":8000"
)

type EnvParams struct {
	ServerEndpoint string
	DbHost         string
	DbPort         string
	DbName         string
	DbUser         string
	DbPass         string
}

func InitFromEnv() EnvParams {
	var ep EnvParams
	ep.DbHost = os.Getenv(ENV_DB_HOST)
	ep.DbPort = os.Getenv(ENV_DB_PORT)
	ep.DbName = os.Getenv(ENV_DB_NAME)
	ep.DbUser = os.Getenv(ENV_DB_USR)
	ep.DbPass = os.Getenv(ENV_DB_PASS)
	ep.ServerEndpoint = os.Getenv(ENV_SVR_ENDPOINT)
	if len(ep.ServerEndpoint) == 0 {
		ep.ServerEndpoint = DefaultEndpoint
	}
	return ep
}

var speciesFilename *string = flag.String("sf", "species.json", "species reference file")

func main() {
	flag.Parse()

	envCfg := InitFromEnv()
	log.Printf("cfg : %+v\n", envCfg)

	// read static species list
	// this should be a persistent mutable repository (database or other)
	speciesMap, err := ReadSpecies(*speciesFilename)

	if err != nil {
		// TODO provide an api to update species in a repo
		log.Fatalf("unable to load a species map : %v", err)
		return
	}

	dap, err := das.Connect(envCfg.DbHost, envCfg.DbPort, envCfg.DbUser, envCfg.DbPass, envCfg.DbName)
	if err != nil {
		log.Printf("unable to connect to database : %v", err)
		return
	}
	go func() {
		err := StartServer(envCfg.ServerEndpoint, &AppHandlers{dap: dap, speciesMap: speciesMap})
		log.Printf("server returned %v - shutting down", err)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Printf("done...")
}
