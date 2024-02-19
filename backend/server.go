package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type WorkflowRun struct {
	Status       string `json:"Status"`
	ReleaseCount int    `json:"Release Count"`
	Velocity     string `json:"Velocity"`
	Volatility   string `json:"Volatility"`
	WorkflowId   int64  `json:"Workflow ID"`
}

var WorkflowRuns []WorkflowRun

var owner = "kyledeanreinford"
var repo = "doraproject"
var branch = "main"

func getLatestWorkflow(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GH_AUTH"))

	opts := &github.ListWorkflowRunsOptions{Branch: branch}

	workflowRuns, _, workflowErr := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)

	for _, workflowRun := range workflowRuns.WorkflowRuns {
		var newWorkflow = WorkflowRun{
			Status:       workflowRun.GetStatus(),
			ReleaseCount: workflowRun.GetRunNumber(),
			Velocity:     "mildly volatile",
			Volatility:   "super fast",
			WorkflowId:   workflowRun.GetWorkflowID(),
		}

		WorkflowRuns = append(WorkflowRuns, newWorkflow)
	}

	if workflowErr != nil {
		fmt.Printf("\nerror: %v\n", workflowErr)
		return
	}

	log.Print("Getting latest workflow")

	var workflowRun = WorkflowRun{
		Status:       workflowRuns.WorkflowRuns[0].GetStatus(),
		ReleaseCount: workflowRuns.WorkflowRuns[0].GetRunNumber(),
		Velocity:     "mildly volatile",
		Volatility:   "super fast",
		WorkflowId:   workflowRuns.WorkflowRuns[0].GetWorkflowID(),
	}

	err := json.NewEncoder(w).Encode(workflowRun)
	if err != nil {
		log.Print("Error encoding workflow")
		return
	}
}

func getAllWorkflows(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getAllWorkflows")
	err := json.NewEncoder(w).Encode(WorkflowRuns)
	if err != nil {
		return
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", getLatestWorkflow)
	myRouter.HandleFunc("/workflows", getAllWorkflows)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		log.Print(".env file loaded")
	}

	handleRequests()
}
