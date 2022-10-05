package cmd

import (
	"log"

	"github.com/charlieegan3/tool-subpub/internal/pkg/config"
	"github.com/charlieegan3/tool-subpub/internal/pkg/jobs"
	"github.com/spf13/cobra"
)

var cfgFile string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the main task",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			log.Fatalf("config file path must be set with --config")
		}

		cfg, err := config.Load(cfgFile)
		if err != nil {
			log.Fatalf("failed to load config file: %s", err)
		}

		for _, j := range cfg.Jobs {
			err = jobs.Run(&j)
			if err != nil {
				log.Fatalf("failed to run job: %s", err)
			}
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&cfgFile, "config", "", "config file path (required)")
	rootCmd.AddCommand(runCmd)
}
