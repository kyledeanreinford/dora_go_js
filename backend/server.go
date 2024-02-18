package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type WorkflowRun struct {
	status       string
	releaseCount int
	velocity     string
	volatility   string
	workflowId   int64
}

func getLatestWorkflow(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GH_AUTH"))

	opts := &github.ListWorkflowRunsOptions{Branch: "main"}

	workflowRuns, _, workflowErr := client.Actions.ListRepositoryWorkflowRuns(ctx, "kyledeanreinford", "doraproject", opts)

	if workflowErr != nil {
		fmt.Printf("\nerror: %v\n", workflowErr)
		return
	}

	log.Print("Getting latest workflow")

	err := json.NewEncoder(w).Encode(WorkflowRun{
		status: workflowRuns.WorkflowRuns[0].GetStatus(),
		releaseCount: workflowRuns.WorkflowRuns[0].GetRunNumber(),
		velocity: "mildly volatile",
		volatility: "super fast",
		workflowId: workflowRuns.WorkflowRuns[0].GetWorkflowID(),
	})
	if err != nil {
		log.Print("Error encoding workflow")
		return
	}
}

func handleRequests() {
	http.HandleFunc("/", getLatestWorkflow)
	log.Fatal(http.ListenAndServe(":8081", nil))
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