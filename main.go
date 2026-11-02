package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Event : GitHub API event struct
type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
	} `json:"payload"`
}

// fetchUserEvents ecth last user events calling GitHub APIs
func fetchUserEvents(username string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user %s not found", username)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(response.Body).Decode(&events); err != nil {
		return nil, err
	}

	return events, nil
}

// translateEvent turns an event into a readable sentence.
func translateEvent(e Event) string {
	switch e.Type {
	case "PushEvent":
		branch := strings.TrimPrefix(e.Payload.Ref, "refs/heads/")
		return fmt.Sprintf("Pushed to %s in %s", branch, e.Repo.Name)
	case "IssuesEvent":
		return fmt.Sprintf("%s an issue in %s", e.Payload.Action, e.Repo.Name)
	case "WatchEvent":
		return fmt.Sprintf("Starred %s", e.Repo.Name)
	case "ForkEvent":
		return fmt.Sprintf("Forked %s", e.Repo.Name)
	case "CreateEvent":
		return fmt.Sprintf("Created a %s in %s", e.Payload.RefType, e.Repo.Name)
	case "PullRequestEvent":
		return fmt.Sprintf("%s a pull request in %s", e.Payload.Action, e.Repo.Name)
	case "IssueCommentEvent":
		return fmt.Sprintf("Commented on an issue in %s", e.Repo.Name)
	default:
		return fmt.Sprintf("%s in %s", e.Type, e.Repo.Name)
	}
}

func main() {
	var username string

	// Username as arg
	if len(os.Args) >= 2 {
		username = os.Args[1]
	} else {
		fmt.Print("Enter GitHub username: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			_, err := fmt.Fprintln(os.Stderr, "Errore di lettura:", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		username = strings.TrimSpace(input)
	}

	if username == "" {
		_, err := fmt.Fprintln(os.Stderr, "Error: empty username")
		if err != nil {
			return
		}
		os.Exit(1)
	}

	events, err := fetchUserEvents(username)
	if err != nil {
		fmt.Println(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("No events found.")
		return
	}

	for _, event := range events {
		fmt.Println("-", translateEvent(event))
	}
}
