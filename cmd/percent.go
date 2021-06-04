package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/shric/bhr/pkg/bhr"
)

func init() {
	percentCmd.PersistentFlags().String(flagName, "", "Name of employee (API key owner if unspecified)")
	percentCmd.PersistentFlags().Int(flagID, -1, "ID of employee")
	rootCmd.AddCommand(percentCmd)
}

var percentCmd = &cobra.Command{
	Use:   "percent",
	Short: "Show percent of an employee",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		c := bhr.PercentCmd{}

		var apiKey string
		var ok bool

		if apiKey, ok = os.LookupEnv("BAMBOOHR_API_KEY"); !ok {
			return errors.New("BAMBOOHR_API_KEY not set")
		}
		result, err := flags.GetString(flagName)
		if err != nil {
			return err
		}
		c.Name, err = nameFlag(result)
		if err != nil {
			return err
		}
		c.ID, err = flags.GetInt(flagID)
		if err != nil {
			return err
		}
		c.Client = bhr.NewClient(apiKey)
		return c.Run()
	},
}
