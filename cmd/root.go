package cmd

import (
	"fmt"
	"os"

	"github.com/pft/internal/database"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "pft",
		Short: "Personal Finance Tracker is a cli app for tracking your expenses",
		Long: `
	With Personal Finance Tracker you can see how much you spend on categories 
	such as Food, Rent and Public Transport over time and plan your budget better.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			fmt.Fprintf(cmd.OutOrStdout(), "Welcome to PFT!\n")
		},
	}
	addCmd := NewAddCmd(database.DB)
	rootCmd.AddCommand(addCmd)

	summaryCmd := NewSummaryCmd(database.DB, "summaries")
	rootCmd.AddCommand(summaryCmd)

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
