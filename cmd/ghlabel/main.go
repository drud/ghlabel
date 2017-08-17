package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
)

// Label represents information about a GitHub issue label.
type Label struct {
	ID     int
	URL    string
	Name   string
	Color  string
	Action string
}

// Client represents information about an instance of ghlabel.
type Client struct {
	// Command line context
	Context context.Context
	// Authenticated GitHub API client.
	GitHub *github.Client
}

// NewClient is the preferred method for making a new authenticated Client.
func NewClient() *Client {
	ctx := context.Background()
	githubToken := os.Getenv("GHLABEL_GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("NewClient error: Please provide a GitHub API token.")
		os.Exit(1)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli := github.NewClient(tc)

	return &Client{Context: ctx, GitHub: cli}
}

// ListByUser lists repositories for a user.
func (c *Client) ListByUser() error {
	if ApplyLabels {
		printCommitHeader()
	}else{
		printPreviewHeader()
	}
	referenceLabels := c.GetLabels(Reference, User)
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		repos, resp, err := c.GitHub.Repositories.List(c.Context, User, opt)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("client_ListByUser: Failed %s request. Status %d",
				resp.Request.Method, resp.StatusCode)
		}

		for _, repo := range repos {
			currentLabels := c.GetLabels(repo.GetName(), User)
			targetLabels := processLabels(referenceLabels, currentLabels)
			// If the run flag was used, execute the staged label changes
			if ApplyLabels {
				err := commit(c.Context, c.GitHub, User, repo.GetName(), targetLabels)
				if err != nil {
					return err
				}
			}else{
				printPreviewData(User, repo.GetName(), targetLabels)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	fmt.Print("\nDone.\n")
	return nil
}

// ListByUserRepository lists a single repository for a user.
func (c *Client) ListByUserRepository() error {
	if ApplyLabels {
		printCommitHeader()
	}else{
		printPreviewHeader()
	}
	referenceLabels := c.GetLabels(Reference, User)
	repo, resp, err := c.GitHub.Repositories.Get(c.Context, User, Repository)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("client_ListByUserRepository: Failed %s request. Status %d",
			resp.Request.Method, resp.StatusCode)
	}

	currentLabels := c.GetLabels(repo.GetName(), User)
	targetLabels := processLabels(referenceLabels, currentLabels)

	// If the run flag was used, execute the staged label changes
	if ApplyLabels {
		err := commit(c.Context, c.GitHub, User, repo.GetName(), targetLabels)
		if err != nil {
			return err
		}
	} else {
		printPreviewData(User, repo.GetName(), targetLabels)
	}
	fmt.Print("\nDone.\n")
	return nil
}

// ListByOrg lists repositories for an organization.
func (c *Client) ListByOrg() error {
	if ApplyLabels {
		printCommitHeader()
	}else{
		printPreviewHeader()
	}
	referenceLabels := c.GetLabels(Reference, Organization)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
		Type:        "all",
	}
	for {
		repos, resp, err := c.GitHub.Repositories.ListByOrg(c.Context, Organization, opt)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("client_ListByOrg: Failed %s request. Status %d",
				resp.Request.Method, resp.StatusCode)
		}

		for _, repo := range repos {
			currentLabels := c.GetLabels(repo.GetName(), Organization)
			targetLabels := processLabels(referenceLabels, currentLabels)
			// If the run flag was used, execute the staged label changes

			if ApplyLabels {
				err := commit(c.Context, c.GitHub, Organization, repo.GetName(), targetLabels)
				if err != nil {
					return err
				}
			}else {
				printPreviewData(Organization, repo.GetName(), targetLabels)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	fmt.Print("\nDone.\n")
	return nil
}

// ListByOrgRepository lists a single repository for an organization.
func (c *Client) ListByOrgRepository() error {
	if ApplyLabels {
		printCommitHeader()
	}else{
		printPreviewHeader()
	}
	referenceLabels := c.GetLabels(Reference, Organization)
	repo, resp, err := c.GitHub.Repositories.Get(c.Context, Organization, Repository)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("client_ListByOrgRepository: Failed %s request. Status %d",
			resp.Request.Method, resp.StatusCode)
	}

	currentLabels := c.GetLabels(repo.GetName(), Organization)
	targetLabels := processLabels(referenceLabels, currentLabels)

	// If the run flag was used, execute the staged label changes
	if ApplyLabels {
		err := commit(c.Context, c.GitHub, Organization, repo.GetName(), targetLabels)
		if err != nil {
			return err
		}
	} else {
		printPreviewData(Organization, repo.GetName(), targetLabels)
	}
	fmt.Print("\nDone.\n")
	return nil
}

// GetLabels returns the currently available label set for a repository.
func (c *Client) GetLabels(repo string, owner string) map[string]Label {
	labelsMap := make(map[string]Label)
	opt := &github.ListOptions{
		PerPage: 10,
	}
	for {
		labels, resp, err := c.GitHub.Issues.ListLabels(c.Context, owner, repo, opt)
		if err != nil {
			log.Fatal(err)
		}
		for _, labelDetail := range labels {
			labelsMap[labelDetail.GetName()] = Label{ID: labelDetail.GetID(), URL: labelDetail.GetURL(), Name: labelDetail.GetName(), Color: labelDetail.GetColor(), Action: ""}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return labelsMap
}

func validateFlags() bool {
	if User == "" && Organization == "" {
		log.Print("You must specify an owner using either the --user or --org flag. Use -h for help.")
		return false
	}
	if Reference == "" {
		log.Print("You must specify a reference repository using the --ref flag. User -h for help.")
		return false
	}
	return true
}

func printCommitHeader() {
	color.Cyan("Applying label changes...\n\n")
}

func printPreviewHeader() {
	color.Yellow("Running in preview mode...\n\n")
	fmt.Println("Description:")
	fmt.Println("  View currently staged label updates.\n")
	fmt.Println("Instructions:")
	fmt.Println("  To apply label changes, use the -a flag.\n")
}

func printPreviewData(owner, repo string, targetLabels map[string]Label) {
	if len(targetLabels) > 0 {
		fmt.Println("\r+-------------------------------------------------+")
		color.Green("%s/%s\n\n", owner, repo)
		r, _ := json.MarshalIndent(targetLabels, "", "    ")
		color.Green("%s\n", string(r))
		return
	}
	fmt.Print("Scanning repo...")
	fmt.Printf("%s\n", repo)
}

func commit(ctx context.Context, client *github.Client, owner string, repo string, labels map[string]Label) error {
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
			_, _, err := client.Issues.EditLabel(ctx, owner, repo, v.Name, label)
			if err != nil {
				log.Printf("%s. Failed to apply label changes.", err.Error())
				return err
			}
		}
		if v.Action == "create" {
			_, _, err := client.Issues.CreateLabel(ctx, owner, repo, label)
			if err != nil {
				log.Printf("%s. Failed to apply label changes.", err.Error())
				return err
			}
		}
		if v.Action == "delete" {
			_, err := client.Issues.DeleteLabel(ctx, owner, repo, v.Name)
			if err != nil {
				log.Printf("%s. Failed to apply label changes.", err.Error())
				return err
			}
		}
	}
	return nil
}

func processLabels(parent map[string]Label, current map[string]Label) map[string]Label {
	labelsMap := make(map[string]Label)
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

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
