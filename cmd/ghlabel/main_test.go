package main

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

var (
	repo1   = "test-ghlabel"
	orgName = "ghlabel"
	client = NewClient()
)

func TestClient_ListByUser(t *testing.T) {
	Reference = "community"
	User = orgName

	assert.NoError(t, client.ListByOrgRepository(), "ListByOrgRepository() returned an error.")
}

func TestClient_ListByUserRepository(t *testing.T) {
	Reference = "community"
	Repository = "junkrepo"
	User = orgName

	assert.NoError(t, client.ListByOrgRepository(), "ListByOrgRepository() returned an error.")
}

func TestClient_ListByOrg(t *testing.T) {
	Reference = "community"
	Organization = orgName

	assert.NoError(t, client.ListByOrg(), "ListByOrg() returned an error.")
}

func TestClient_ListByOrgRepository(t *testing.T) {
	Reference = "community"
	Repository = "junkrepo"
	Organization = orgName

	assert.NoError(t, client.ListByOrgRepository(), "ListByOrgRepository() returned an error.")
}

// TestClient_GetLabels makes sure values returned by GetLabels are contained
func TestClient_GetLabels(t *testing.T) {
	expectedLabels := []string{"actionable", "hibernate", "showstopper", "incubate",
	"work in progress", "security", "needs decision", "needs tests", "needs docs"}
	actualLabels := client.GetLabels("community", orgName)

	for actual, _ := range actualLabels {
		assert.Contains(t, expectedLabels, actual, "GetLabels() Test failed.")
	}
}

func createRepo(org string, name string) (*github.Repository, *github.Response, error) {
	repo := &github.Repository{
		Name:      github.String(name),
		Private:   github.Bool(true),
		HasIssues: github.Bool(true),
	}

	return client.GitHub.Repositories.Create(client.Context, org, repo)
}

func deleteRepo(org string, name string) (*github.Response, error) {
	return client.GitHub.Repositories.Delete(client.Context, org, name)
}
