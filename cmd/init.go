package cmd

import (
	"fmt"

	"commit-grid/internal/config"
	"commit-grid/internal/tui"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Asistente interactivo de configuración",
	RunE: func(cmd *cobra.Command, args []string) error {
		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Render("Commit Grid Draw")
		fmt.Println(title, "— configuración inicial")
		cfg, err := tui.Onboard()
		if err != nil { return err }
		path, err := config.Save(cfg)
		if err != nil { return err }
		fmt.Println("✅ Config guardada en:", path)
		fmt.Println("Sugerencia: ejecuta `commit-grid enable` para programarlo.")
		return nil
	},
}

func init() { rootCmd.AddCommand(initCmd) }
