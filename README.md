# Readme

## Running the Server

1. `cd cadenceservice`
2. `docker-compose up`
3. Register the `simple-domain` with `docker run --network=host --rm ubercadence/cli:master --do simple-domain domain register --rd 10`

## Worker and API Server

Navigate back to the project root folder. Make sure go is installed in system and install dependencies with `go mod vendor`

* HttpServer
    1. `make httpserver`
    2. `./bins/httpserver`

* Worker
    1. `make worker`
    2. `./bins/worker`

## Endpoints

1. To start workflow
   * POST request to `http://localhost:3030/api/start-hello-world` 

2. To signal workflow
    * Check Cadence UI for the WorkflowID of the child workflow
    * POST Request to `http://localhost:3030/api/signal-hello-world?workflowId=<workflowId>&age=25`
    * Make sure to replace <workflowId> with the id retrieved from cadence web ui