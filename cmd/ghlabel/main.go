package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type label struct {
	ID     int
	URL    string
	Name   string
	Color  string
	Action string
}

func main() {
	ctx, cli := getClient()

	owner, repo, parent, action := parseFlags()
	parentLabels := getLabels(ctx, cli, owner, parent)
	if repo != "" {
		currentLabels := getLabels(ctx, cli, owner, repo)
		targetLabels := processLabels(parentLabels, currentLabels)
		commit(ctx, cli, owner, repo, targetLabels)
		return
	}

	switch action {
	case "preview":
		previewAllRepos(ctx, cli, owner, parentLabels)
	case "":
		updateAllRepos(ctx, cli, owner, parentLabels)
	}
}

func parseFlags() (owner string, repo string, parent string, action string) {
	flag.StringVar(&owner, "owner", "", "The organization or user that owns the repositories.")
	flag.StringVar(&repo, "repo", "", "A specific repository to focus on.")
	flag.StringVar(&parent, "parent", "", "The repository to replicate labels from.")
	flag.Parse()
	action = flag.Arg(0)

	if owner == "" {
		log.Fatal("The owner flag is required. Use -h for help.")
	}
	if parent == "" {
		log.Fatal("The parent flag is required. Use -h for help.")
	}
	return owner, repo, parent, action
}

func getClient() (ctx context.Context, cli *github.Client) {
	ctx = context.Background()
	githubToken := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli = github.NewClient(tc)
	return ctx, cli
}

func commit(ctx context.Context, client *github.Client, owner string, repo string, labels map[string]label) {
	for _, v := range labels {
		label := new(github.Label)

		color := string(v.Color)
		name := string(v.Name)
		url := string(v.URL)
		id := int(v.ID)

		label.ID = &id
		label.Color = &color
		label.URL = &url
		label.Name = &name

		if v.Action == "edit" {
			client.Issues.EditLabel(ctx, owner, repo, v.Name, label)
		}
		if v.Action == "create" {
			client.Issues.CreateLabel(ctx, owner, repo, label)
		}
		if v.Action == "delete" {
			client.Issues.DeleteLabel(ctx, owner, repo, v.Name)
		}
	}
}

func processLabels(parent map[string]label, current map[string]label) map[string]label {
	labelsMap := make(map[string]label)
	// Move all parent items into labelsMap with action create
	for k, v := range parent {
		v.ID = 0
		v.URL = ""
		v.Action = "create"
		labelsMap[k] = v
	}

	// Move all current items into labelsMap with updated action
	for k, v := range current {
		if targetLabel, ok := labelsMap[v.Name]; ok {
			// update color if it is different
			if v.Color != targetLabel.Color {
				v.Action = "edit"
				v.Color = targetLabel.Color
			}
		} else {
			v.Action = "delete"
		}
		labelsMap[k] = v
	}

	// Remove anything that has a nil action.
	for _, v := range labelsMap {
		if v.Action == "" {
			delete(labelsMap, v.Name)
		}
	}
	return labelsMap
}

// Get the currently available label set for a repository.
func getLabels(ctx context.Context, client *github.Client, owner string, repo string) map[string]label {
	labelsMap := make(map[string]label)
	opt := &github.ListOptions{
		PerPage: 10,
	}
	for {
		labels, resp, err := client.Issues.ListLabels(ctx, owner, repo, opt)
		if err != nil {
			log.Fatal(err)
		}
		for _, labelDetail := range labels {
			labelsMap[labelDetail.GetName()] = label{ID: labelDetail.GetID(), URL: labelDetail.GetURL(), Name: labelDetail.GetName(), Color: labelDetail.GetColor(), Action: ""}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
		break
	}
	return labelsMap
}

func previewAllRepos(ctx context.Context, client *github.Client, owner string, parentLabels map[string]label) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Type:        "all",
	}

	fmt.Printf("\nOwner: %s\n\n", owner)
	fmt.Println("             CHANGES STAGED FOR COMMIT             ")
	fmt.Println("===================================================")
	fmt.Printf("| %-28s| %-18s\n", "REPOSITORY", "ACTIONABLE")
	fmt.Print("===================================================")

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			currentLabels := getLabels(ctx, client, owner, repo.GetName())
			targetLabels := processLabels(parentLabels, currentLabels)
			fmt.Printf("\n| %-28s", repo.GetName())
			r, _ := json.MarshalIndent(targetLabels, "|                             |", "  ")
			fmt.Printf("| %-18s\n", string(r))
			fmt.Print("---------------------------------------------------")
		}

		if resp.NextPage == 0 {
			fmt.Printf("\n\nRun `ghlabel --owner --parent run` to proceed with these changes.\n\n")
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
}

func updateAllRepos(ctx context.Context, client *github.Client, owner string, parentLabels map[string]label) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Type:        "all",
	}

	fmt.Printf("\nOwner: %s\n\n", owner)
	fmt.Println("               UPDATING REPOSITORIES              ")
	fmt.Println("===================================================")
	fmt.Printf("| %-28s| %-18s\n", "REPOSITORY", "ACTIONABLE")
	fmt.Print("===================================================")

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			currentLabels := getLabels(ctx, client, owner, repo.GetName())
			targetLabels := processLabels(parentLabels, currentLabels)
			commit(ctx, client, owner, repo.GetName(), targetLabels)
			fmt.Printf("\n| %-28s", repo.GetName())
			r, _ := json.MarshalIndent(targetLabels, "|                             |", "  ")
			fmt.Printf("| %-18s\n", string(r))
			fmt.Print("---------------------------------------------------")
		}
		if resp.NextPage == 0 {
			fmt.Printf("\nSuccessfully updated labels.\n")
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
}
