package main

import (
	"github.com/spf13/cobra"
	"os"
)

// RootCmd is the Cobra root for ghlabel command.
var RootCmd = &cobra.Command{
	Use:   "ghlabel",
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
		if User != "" {
			if Repository != "" {
				client.ListByUserRepository()
				return
			}
			client.ListByUser()
			return
		}
		if Organization != "" {
			if Repository != "" {
				client.ListByOrgRepository()
				return
			}
			client.ListByOrg()
			return
		}
	},
}

// Globally accessible flags
var (
	User         string
	Organization string
	Repository   string
	Reference    string
	Run          bool
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&User, "user", "u", "", "The user that owns the repositories")
	RootCmd.PersistentFlags().StringVarP(&Organization, "org", "o", "", "The organization that owns the repositories.")
	RootCmd.PersistentFlags().StringVarP(&Repository, "repo", "", "", "A specific repository to sync.")
	RootCmd.PersistentFlags().StringVarP(&Reference, "ref", "", "", "Required: the repository to replicate labels from.")
	RootCmd.PersistentFlags().BoolVarP(&Run, "run", "r", false, "Run currently staged label updates.")
}

// Execute runs Cobra
func Execute() {
	RootCmd.Execute()
}
