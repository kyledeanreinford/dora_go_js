package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type WorkflowRun struct {
	Status       string `json:"Status"`
	ReleaseCount int `json:"Release Count"`
	Velocity     string `json:"Velocity"`
	Volatility   string `json:"Volatility"`
	WorkflowId   int64 `json:"Workflow ID"`
}

var WorkflowRuns []WorkflowRun

func getLatestWorkflow(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	//client := github.NewClient(nil).WithAuthToken(os.Getenv("GH_AUTH"))
	//
	//opts := &github.ListWorkflowRunsOptions{Branch: "main"}
	//
	//workflowRuns, _, workflowErr := client.Actions.ListRepositoryWorkflowRuns(ctx, "kyledeanreinford", "doraproject", opts)
	//
	//if workflowErr != nil {
	//	fmt.Printf("\nerror: %v\n", workflowErr)
	//	return
	//}

	log.Print("Getting latest workflow")

	var mockWorkflowRun = WorkflowRun{
		Status:       "completed",
		ReleaseCount: 1000,
		Velocity:     "super fast",
		Volatility:   "mildly volatile",
		WorkflowId:   1234,
	}

	//var workflowRun = WorkflowRun{
	//	status: workflowRuns.WorkflowRuns[0].GetStatus(),
	//	releaseCount: workflowRuns.WorkflowRuns[0].GetRunNumber(),
	//	velocity: "mildly volatile",
	//	volatility: "super fast",
	//	workflowId: workflowRuns.WorkflowRuns[0].GetWorkflowID(),
	//}

	err := json.NewEncoder(w).Encode(mockWorkflowRun)
	if err != nil {
		log.Print("Error encoding workflow")
		return
	}
}

func getAllWorkflows(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getAllWorkflows")
	json.NewEncoder(w).Encode(WorkflowRuns)
}

func handleRequests() {
	http.HandleFunc("/", getLatestWorkflow)
	http.HandleFunc("/workflows", getAllWorkflows)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		log.Print(".env file loaded")
	}

	WorkflowRuns = []WorkflowRun{
		{"in progress", 123, "fast!", "volatile", 12345},
		{"completed", 4, "faster!", "more volatile", 555},
	}
	handleRequests()
}
