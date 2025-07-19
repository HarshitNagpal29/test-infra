package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	Type    string          `json:"type"`
	Repo    RepoInfo        `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

type RepoInfo struct {
	Name string `json:"name"`
}

type PushEventPayload struct {
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Message string `json:"message"`
}

func main() {
	username := os.Args[1]
	if username == "" {
		println("No username provided")
		return
	}
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := http.Get(url)
	if err != nil {
		println("Error fetching events:", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		println("Error fetching events:", resp.Status)
		return
	}

	var events []Event
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Error reading response body:", err.Error())
		return
	}
	err = json.Unmarshal(body, &events)
	if err != nil {
		println("Error parsing JSON:", err.Error())
		return
	}
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			var payload PushEventPayload
			err := json.Unmarshal(event.Payload, &payload)
			if err != nil {
				println("Error parsing PushEvent payload:", err.Error())
				continue
			}
			fmt.Printf("Repository: %s\n", event.Repo.Name)
			for _, commit := range payload.Commits {
				fmt.Printf("Commit Message: %s\n", commit.Message)
			}
		case "IssuesEvent":
			fmt.Printf("Repository: %s\n", event.Repo.Name)
			fmt.Println("Issue event detected, but no specific details provided in this example.")
		case "WatchEvent":
			fmt.Printf("Repository: %s\n", event.Repo.Name)
			fmt.Println("Watch event detected, but no specific details provided in this example.")
		default:
			fmt.Printf("Repository: %s\n", event.Repo.Name)
			fmt.Printf("Event Type: %s\n", event.Type)
			fmt.Println("No specific details provided for this event type in this example.")
		}

	}
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Rate Limit Remaining: %s\n", resp.Header.Get("X-RateLimit-Remaining"))
	fmt.Println("Raw response body:")
	fmt.Println(string(body))
	fmt.Printf("Fetched %d events\n", len(events))

}
