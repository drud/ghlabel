package main

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Globally accessible flags
var (
	User         string
	Organization string
	Repository   string
	Reference    string
	ApplyLabels  bool
)

// RootCmd is the Cobra root for ghlabel command.
var RootCmd = &cobra.Command{
	Use:   "ghlabel --owner --reference -flags",
	Short: "ghlabel automatically manages issue labels.",
	Long: `GitHub Label (ghlabel) automatically updates
			a user or organization's GitHub issue labels.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validation requires the owner and parent flags.
		if !validateFlags() {
			os.Exit(1)
		}
		// Create a new ghlabel Client.
		client := NewClient()
		if User != "" && Repository != "" {
			// If we have both a user and repository then we run the tool on just that repository.
			client.ListByUserRepository()
		} else if User != "" {
			// If we Just have a user but no repository, then run the tool for all repositories for that user.
			client.ListByUser()
		} else if Organization != "" && Repository != "" {
			// If we have an org and a repository, run the tool for just that repository.
			client.ListByOrgRepository()
		} else if Organization != "" {
			// If we just have an org and no repository, run the tool on all repos in that org.
			client.ListByOrg()
		} else {
			log.Fatal("You must specify either an organization or user. Use -h for help.")
		}
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&User, "user", "u", "", "The user that owns the repositories")
	RootCmd.PersistentFlags().StringVarP(&Organization, "org", "o", "", "The organization that owns the repositories.")
	RootCmd.PersistentFlags().StringVarP(&Repository, "repo", "", "", "A specific repository to sync.")
	RootCmd.PersistentFlags().StringVarP(&Reference, "ref", "", "", "Required: the repository from which to replicate labels.")
	RootCmd.PersistentFlags().BoolVarP(&ApplyLabels, "apply", "a", false, "Apply currently staged label changes.")
}

// Execute runs Cobra
func Execute() {
	RootCmd.Execute()
}
