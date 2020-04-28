package main

import (
	"fmt"
	"github.com/nndd91/cadence-api-example/app/adapters/cadenceAdapter"
	"github.com/nndd91/cadence-api-example/app/config"
	"github.com/nndd91/cadence-api-example/app/worker/workflows"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func startWorkers(h *cadenceAdapter.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	startWorkers(&cadenceClient, workflows.TaskListName)
	// The workers are supposed to be long running process that should not exit.
	select {}
}
