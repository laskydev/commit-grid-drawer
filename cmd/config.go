package cmd

import (
	"fmt"

	"commit-grid/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{ Use: "config", Short: "Lee/escribe configuración" }

var configGetCmd = &cobra.Command{
	Use: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, _, err := config.Load()
		if err != nil { return err }
		fmt.Println(config.Debug(c))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
}
