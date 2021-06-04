package cmd

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/shric/bhr/pkg/bhr"
	"github.com/spf13/cobra"
)

const (
	flagName  = "name"
	flagID    = "id"
	flagImage = "image"
)

func init() {
	employeeCmd.PersistentFlags().String(flagName, "", "Name of employee (API key owner if unspecified)")
	employeeCmd.PersistentFlags().Int(flagID, -1, "ID of employee")
	employeeCmd.PersistentFlags().Bool(flagImage, false, "Display profile image using sixel")
	rootCmd.AddCommand(employeeCmd)
}

var employeeCmd = &cobra.Command{
	Use:     "employee",
	Aliases: []string{"emp", "get"},
	Short:   "Show information about an employee",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		c := bhr.EmployeeCmd{}

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

		c.Image, err = flags.GetBool(flagImage)
		if err != nil {
			return err
		}

		c.Client = bhr.NewClient(apiKey)
		return c.Run()
	},
}

func nameFlag(flag string) (string, error) {
	if flag != "" {
		flag = strings.Replace(flag, " ", ".*", -1)
		flag = "(?i)" + flag
		_, err := regexp.Compile(flag)
		if err != nil {
			return "", err
		}
	}
	return flag, nil
}
