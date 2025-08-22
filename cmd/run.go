package cmd

import (
	"fmt"
	"time"

	"commit-grid/internal/config"
	"commit-grid/internal/gitutil"
	"commit-grid/internal/pattern"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Ejecuta la tarea del día (commits)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, _, err := config.Load()
		if err != nil { return fmt.Errorf("config no encontrada, ejecuta 'commit-grid init': %w", err) }
		// Aquí podrías setear TZ si necesitas cálculos locales
		if err := pattern.TouchToday(c.RepoPath); err != nil { return err }
		intensity := pattern.IntensityToday(pattern.Strategy(c.IntensityStrategy), c.IntensityValue)
		for i := 1; i <= intensity; i++ {
			msg := fmt.Sprintf("grid: %s (%d/%d)", time.Now().UTC().Format("2006-01-02"), i, intensity)
			if err := gitutil.CommitPush(c.RepoPath, c.GitUser, c.GitEmail, msg); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(runCmd) }
