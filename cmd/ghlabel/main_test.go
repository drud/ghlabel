package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	orgName = "ghlabel"
	usrName = "drud-ghlabel-test"
	client  = NewClient()
)

func TestClient_ListByUser(t *testing.T) {
	Reference = "community"
	User = usrName

	assert.NoError(t, client.ListByUser(), "ListByOrg() returned an error.")
}

func TestClient_ListByUserRepository(t *testing.T) {
	Reference = "community"
	Repository = "junkrepo"
	User = usrName

	assert.NoError(t, client.ListByUserRepository(), "ListByOrgRepository() returned an error.")
}

func TestClient_ListByOrg(t *testing.T) {
	Reference = "community"
	Organization = orgName

	assert.NoError(t, client.ListByOrg(), "ListByOrg() returned an error.")
}

// TestClient_ListByOrgRepository
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

	for actual := range actualLabels {
		assert.Contains(t, expectedLabels, actual, "GetLabels() Test failed.")
	}
}
