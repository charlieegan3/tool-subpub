package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "subpub",
	Short: "Simple CLI to fetch data, make substitutions and upload the data",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
