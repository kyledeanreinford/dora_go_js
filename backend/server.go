package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

type WorkflowRun struct {
	Status       string
	ReleaseCount int
	Velocity     string
	Volatility   string
	WorkflowId   int64
}

var WorkflowRuns []WorkflowRun

type Repo struct {
	Name string
	Owner  string
	Branch string
	Auth string
}

var repo Repo

var help bool

func getLatestWorkflow(w http.ResponseWriter, r *http.Request) {
	log.Print("Getting latest workflow")

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GH_AUTH"))

	opts := &github.ListWorkflowRunsOptions{Branch: repo.Branch}

	workflowRuns, _, workflowErr := client.Actions.ListRepositoryWorkflowRuns(ctx, repo.Owner, repo.Name, opts)

	if workflowErr != nil {
		fmt.Printf("\nerror: %v\n", workflowErr)
		return
	}

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
	reqBody, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &repo)
	if err != nil {
		log.Print(err)
		return
	}

	WorkflowRuns = nil

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(repo.Auth)

	opts := &github.ListWorkflowRunsOptions{Branch: repo.Branch}

	workflowRuns, _, workflowErr := client.Actions.ListRepositoryWorkflowRuns(ctx, repo.Owner, repo.Name, opts)

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

	encodeErr := json.NewEncoder(w).Encode(WorkflowRuns)
	if encodeErr != nil {
		return
	}

	fmt.Printf("Here's a repo: %v", repo)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("This is the homepage")
	if err != nil {
		return
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/workflow", getLatestWorkflow)
	myRouter.HandleFunc("/workflows", getAllWorkflows).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	flag.StringVar(&repo.Owner, "o", "kyledeanreinford", "Repository owner")
	flag.StringVar(&repo.Name, "r", "doraproject", "Repository name")
	flag.StringVar(&repo.Branch, "b", "main", "Branch name")
	flag.BoolVar(&help, "help", false, "Help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		log.Print(".env file loaded")
	}

	handleRequests()
}
