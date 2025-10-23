package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Request represents a Git LFS transfer request
type Request struct {
	Event   string                   `json:"event"`
	Objects []map[string]interface{} `json:"objects,omitempty"`
}

// Response represents a Git LFS transfer response
type Response struct {
	Event   string                   `json:"event"`
	Success bool                     `json:"success"`
	Error   string                   `json:"error,omitempty"`
	Objects []map[string]interface{} `json:"objects,omitempty"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		var request Request
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			continue // Skip invalid JSON
		}

		logRequest(request)

		response := handleRequest(request)
		logResponse(response)

		// Write response to stdout
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func logRequest(request Request) {
	fmt.Fprintln(os.Stderr, "\n== Request ==")
	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	fmt.Fprintln(os.Stderr, string(requestJSON))
	fmt.Fprintln(os.Stderr, "================")
}

func logResponse(response Response) {
	fmt.Fprintln(os.Stderr, "\n== Response ==")
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Fprintln(os.Stderr, string(responseJSON))
	fmt.Fprintln(os.Stderr, "================")
}

func handleRequest(request Request) Response {
	switch request.Event {
	case "init":
		return Response{Event: "init", Success: true}
	case "terminate":
		return Response{Event: "terminate", Success: true}
	case "upload":
		return handleUpload(request)
	case "download":
		return handleDownload(request)
	default:
		return Response{
			Event:   request.Event,
			Success: false,
			Error:   "Unsupported event",
		}
	}
}

func handleUpload(request Request) Response {
	if len(request.Objects) == 0 {
		return Response{
			Event:   "upload",
			Success: false,
			Error:   "No object specified",
		}
	}

	object := request.Objects[0]
	oid, _ := object["oid"].(string)
	size, _ := object["size"].(float64)

	return Response{
		Event:   "upload",
		Success: true,
		Objects: []map[string]interface{}{
			{
				"oid":  oid,
				"size": size,
				"actions": map[string]interface{}{
					"upload": map[string]interface{}{
						"href": fmt.Sprintf("https://example.com/upload/%s", oid),
					},
				},
			},
		},
	}
}

func handleDownload(request Request) Response {
	if len(request.Objects) == 0 {
		return Response{
			Event:   "download",
			Success: false,
			Error:   "No object specified",
		}
	}

	object := request.Objects[0]
	oid, _ := object["oid"].(string)
	size, _ := object["size"].(float64)

	return Response{
		Event:   "download",
		Success: true,
		Objects: []map[string]interface{}{
			{
				"oid":  oid,
				"size": size,
				"actions": map[string]interface{}{
					"download": map[string]interface{}{
						"href": fmt.Sprintf("https://example.com/download/%s", oid),
					},
				},
			},
		},
	}
}
