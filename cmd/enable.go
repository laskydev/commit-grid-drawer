package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"commit-grid/internal/config"
	"commit-grid/internal/schedule"

	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Activa el programador diario (cron/launchd)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, _, err := config.Load()
		if err != nil { return fmt.Errorf("config no encontrada, corre 'commit-grid init'") }
		bin, err := os.Executable()
		if err != nil { return err }
		// resolver symlinks
		if l, err := filepath.EvalSymlinks(bin); err == nil { bin = l }
		if err := schedule.Enable(bin, c.Hour24, c.Minute); err != nil { return err }
		fmt.Println("✅ Programador activado")
		return nil
	},
}
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Desactiva el programador diario",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := schedule.Disable(); err != nil { return err }
		fmt.Println("🛑 Programador desactivado")
		return nil
	},
}
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Muestra estado del programador",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := schedule.Status()
		if err != nil { return err }
		fmt.Println("Estado:", s)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(enableCmd, disableCmd, statusCmd)
}
