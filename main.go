package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Request struct represents the expected JSON input
type Request struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

// Response struct represents the JSON output
type Response struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

// executeCommand runs the given command and returns the result
func executeCommand(cmd string, args []string) Response {
	// Create command
	command := exec.Command(cmd, args...)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	// Run command
	err := command.Run()

	// Get exit code
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}
}

// handler processes incoming requests
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Execute command
	result := executeCommand(req.Cmd, req.Args)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Get PORT from environment, default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/execute", handler)
	log.Printf("Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
