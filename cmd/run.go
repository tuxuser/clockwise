package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/syncfast/clockwise/internal/scrape"
	"github.com/syncfast/clockwise/internal/tui"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run clockwise",
	Long:   ``,
	PreRun: toggleDebug,
	RunE: func(cmd *cobra.Command, args []string) error {
		manual, err := cmd.Flags().GetBool("manual")
		if err != nil {
			return err
		}

		// We declare data here because it's consumed by both the `tui` and
		// `scrape` packages.
		var data tui.Data

		if !manual {
			log.Println("Initializing playwright to scrape participant count.")
			pw, err := scrape.InitializePlaywright()
			if err != nil {
				return err
			}

			log.Info("Initializing TUI.")
			url, err := cmd.Flags().GetString("url")
			go func() {
				err = scrape.GetParticipants(url, 1, &data, pw)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		tui.Start(manual, &data)
		log.Info("Clockwise has been stopped.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("manual", "m", false, "Set participant count manually via the TUI.")
	runCmd.Flags().StringP("url", "u", "", "The Zoom ")
}