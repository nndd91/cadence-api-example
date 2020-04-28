package main

import (
	"context"
	"encoding/json"
	"github.com/nndd91/cadence-api-example/app/adapters/cadenceAdapter"
	"github.com/nndd91/cadence-api-example/app/config"
	"github.com/nndd91/cadence-api-example/app/worker/workflows"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	cadenceAdapter *cadenceAdapter.CadenceAdapter
	logger         *zap.Logger
}

func (h *Service) triggerHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountId := r.URL.Query().Get("accountId")

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.Workflow, accountId)
		if err != nil {
			http.Error(w, "Error starting workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Started work flow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) signalHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		workflowId := r.URL.Query().Get("workflowId")
		age, err := strconv.Atoi(r.URL.Query().Get("age"))
		if err != nil {
			h.logger.Error("Failed to parse age from request!")
		}

		err = h.cadenceAdapter.CadenceClient.SignalWorkflow(context.Background(), workflowId, "", workflows.SignalName, age)
		if err != nil {
			http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Signaled work flow with the following params!", zap.String("WorkflowId", workflowId), zap.Int("Age", age))

		js, _ := json.Marshal("Success")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func main() {
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	service := Service{&cadenceClient, appConfig.Logger}
	http.HandleFunc("/api/start-hello-world", service.triggerHelloWorld)
	http.HandleFunc("/api/signal-hello-world", service.signalHelloWorld)

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
