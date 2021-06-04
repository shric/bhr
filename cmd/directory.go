package cmd

import (
	"errors"
	"os"
	"regexp"

	"github.com/shric/bhr/pkg/bhr"
	"github.com/spf13/cobra"
)

const (
	flagDepartment = "department"
	flagTitle      = "title"
)

func init() {
	directoryCmd.PersistentFlags().String(flagDepartment, "", "Filter by department (case insensitive regex)")
	directoryCmd.PersistentFlags().String(flagTitle, "", "Filter by title (case insensitive regex)")
	rootCmd.AddCommand(directoryCmd)
}

var directoryCmd = &cobra.Command{
	Use:     "directory",
	Aliases: []string{"dir"},
	Short:   "List directory of employees",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		c := bhr.DirectoryCmd{}

		var apiKey string
		var ok bool
		var err error

		if apiKey, ok = os.LookupEnv("BAMBOOHR_API_KEY"); !ok {
			return errors.New("BAMBOOHR_API_KEY not set")
		}

		c.Department, err = flags.GetString(flagDepartment)
		if err != nil {
			return err
		}
		c.Department = "(?i)" + c.Department
		_, err = regexp.Compile(c.Department)
		if err != nil {
			return err
		}

		c.Title, err = flags.GetString(flagTitle)
		if err != nil {
			return err
		}
		c.Title = "(?i)" + c.Title
		_, err = regexp.Compile(c.Title)
		if err != nil {
			return err
		}

		c.Client = bhr.NewClient(apiKey)
		return c.Run()
	},
}
