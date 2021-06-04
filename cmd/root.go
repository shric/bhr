package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bhr",
		Short: "bhr is a command line interface for BambooHR",
		Long:  `A command line interface for BambooHR`,
	}
)

func init() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
