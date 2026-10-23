package main

import (	
	"fmt"
	"os"
	"io"
	"net/http"
	"encoding/json"
)

type Repository struct {
	Name 			string 			`json:"name"`
}

type Event struct {
	Type 			string 			`json:"type"`
	Repo 			Repository 		`json:"repo"`
	CreatedAt 		string 			`json:"created_at"`
}

func describeEvent(event Event)string {
	switch event.Type {
    case "PushEvent":
        return "Pushed commits to " + event.Repo.Name
    case "WatchEvent":
        return "Starred " + event.Repo.Name
    case "IssuesEvent":
        return "Updated an issue in " + event.Repo.Name
    case "CreateEvent":
        return "Created a branch or repository in " + event.Repo.Name
    case "ForkEvent":
        return "Forked " + event.Repo.Name
    case "PullRequestEvent":
        return "Updated a pull request in " + event.Repo.Name
    default:
        return event.Type + " in " + event.Repo.Name
    }
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: github-activity <username>")
		return
	}

	username := os.Args[1];
	apiUrl := "https://api.github.com/users/" + username + "/events"

	response, err := http.Get(apiUrl)

	if err != nil {
		fmt.Println("error fetching activity: ", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("github api returned: ", response.Status)
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("error reading response", err)
		return
	}

	var events []Event

	err = json.Unmarshal(body, &events)

	if err != nil {
		fmt.Println("error parsing json:", err)
		return
	}

	if len(events) == 0 {
		fmt.Println("no recent github activity found for", username)
		return
	}

	fmt.Println("-----recent github activity for:", username, "-----")

	for _,event := range events {
		fmt.Printf("- %s\n", describeEvent(event))
	}

}
