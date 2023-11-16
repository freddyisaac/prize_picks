package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"dinocage/das"
)

const (
	EnvDBHost       = "ENV_DB_HOST"
	EnvDBPort       = "ENV_DB_PORT"
	EnvDBName       = "ENV_DB_NAME"
	EnvDBUsr        = "ENV_DB_USR"
	EnvDBPass       = "ENV_DB_PWD"
	EnvSvrEndpoint  = "ENV_SVR_ENDPOINT"
	DefaultEndpoint = ":8000"
)

type EnvParams struct {
	ServerEndpoint string
	DbHost         string
	DbPort         string
	DbName         string
	DbUser         string
	DbPass         string
}

// extract params from env - should implement defaults
func InitConfigFromEnv() EnvParams {
	var ep EnvParams
	ep.DbHost = os.Getenv(EnvDBHost)
	ep.DbPort = os.Getenv(EnvDBPort)
	ep.DbName = os.Getenv(EnvDBName)
	ep.DbUser = os.Getenv(EnvDBUsr)
	ep.DbPass = os.Getenv(EnvDBPass)
	ep.ServerEndpoint = os.Getenv(EnvSvrEndpoint)
	if len(ep.ServerEndpoint) == 0 {
		ep.ServerEndpoint = DefaultEndpoint
	}
	return ep
}

var speciesFilename *string = flag.String("sf", "species.json", "species reference file")

func main() {
	flag.Parse()

	envCfg := InitConfigFromEnv()
	log.Printf("cfg : %+v\n", envCfg)

	// read static species list
	// this should be a persistent mutable repository (database or other)
	speciesMap, err := ReadSpecies(*speciesFilename)

	if err != nil {
		// TODO provide an api to update species in a repo
		log.Fatalf("unable to load a species map : %v", err)
		return
	}

	// connect to database
	dap, err := das.Connect(envCfg.DbHost, envCfg.DbPort, envCfg.DbUser, envCfg.DbPass, envCfg.DbName)
	if err != nil {
		log.Printf("unable to connect to database : %v", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// start server in background
		err := StartServer(ctx, envCfg.ServerEndpoint, &AppHandlers{dap: dap, speciesMap: speciesMap})
		log.Printf("server returned %v - shutting down", err)
		wg.Done()
	}()

	// keep server alive until
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	// try to shutdown gracefully

	// cancel context to initiate server shutdown
	cancel()
	// wait for http server to shudown
	wg.Wait()
	// close database
	dap.Close()
	log.Printf("done...")
}
