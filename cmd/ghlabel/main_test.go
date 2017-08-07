package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

var (
	repo1   = "test-github-labels-1"
	repo2   = "test-github-labels-2"
	orgName = "drud-admin"
)

func setup() {
	ctx, cli := getClient()
	repo, response, err := createRepo(ctx, cli, orgName, repo1)
	if err != nil {
		fmt.Println("fooooooooo")
		fmt.Println(repo)
		fmt.Println(response)

	}

	// create two labels in our parent repo
	label := new(github.Label)
	color := "0e8a16"
	name := "foo"
	label.Color = &color
	label.Name = &name
	cli.Issues.CreateLabel(ctx, orgName, repo1, label)

	color = "c2e0c6"
	name = "bar"
	label.Color = &color
	label.Name = &name
	cli.Issues.CreateLabel(ctx, orgName, repo1, label)

	// Delete the default labels
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "bug")
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "duplicate")
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "enhancement")
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "help wanted")
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "invalid")
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "question")
	cli.Issues.DeleteLabel(ctx, orgName, repo1, "wontfix")

	createRepo(ctx, cli, orgName, repo2)

	// Create a label in our child repo that has the same label with a different color.
	color = "1d76db"
	label.Color = &color
	label.Name = &name
	cli.Issues.CreateLabel(ctx, orgName, repo2, label)

}

func TestDevLogs(t *testing.T) {
	assert.Equal(t, 123, 123, "they should be equal")
}

func teardown() {
	ctx, cli := getClient()
	//cli.Repositories.Delete(ctx, orgName, repo1)
	cli.Repositories.Delete(ctx, orgName, repo2)
}

func TestCredentials(t *testing.T) {
	setup()
	teardown()
}

func createRepo(ctx context.Context, cli *github.Client, org string, name string) (*github.Repository, *github.Response, error) {

	repo := &github.Repository{
		Name:    github.String(name),
		Private: github.Bool(true),
		HasIssues: github.Bool(true),
	}

	return cli.Repositories.Create(ctx, org, repo)
}

func deleteRepo(ctx context.Context, cli *github.Client, org string, name string) (*github.Response, error) {

	return cli.Repositories.Delete(ctx, org, name)
}
