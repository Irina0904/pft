package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-echarts/go-echarts/v2/components"

	"github.com/pft/internal/summary"
	"github.com/spf13/cobra"
)

func NewSummaryCmd(db *sql.DB, directory string) *cobra.Command {
	summaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Show report on spending over time",
		Long: `
	Shows how much money was spent on each category for the selected month`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do Stuff Here
			month, _ := cmd.Flags().GetString("month")

			if month != "" {
				entries, err := summary.GetSummary(month, db)
				if err != nil {
					return fmt.Errorf("could not retrieve summary. %s", err)
				}

				if len(entries) == 0 {
					fmt.Printf("There was no data for the month of %s\n", strings.ToUpper(month[:1]) + strings.ToLower(month[1:]))
					return nil
				}

				page := components.NewPage()
				page.AddCharts(summary.PieSummary(entries, month))

				filename := fmt.Sprintf("summary_%s.html", strings.ToLower(month))
				filePath := filepath.Join(directory, filename)

				// Create the directory if it doesn't exist
				err = os.MkdirAll(directory, 0755) // 0755 grants read/write/execute for owner, read/execute for group/others
				if err != nil {
					fmt.Printf("Error creating directory: %v\n", err)
				}

				f, err := os.Create(filePath)
				if err != nil {
					return fmt.Errorf("could not create summary file: %s", err)
				}
				page.Render(io.MultiWriter(f))
			}
			return nil
		},
	}
	summaryCmd.Flags().String("month", "", "Specify the month to report on")

	return summaryCmd
}
