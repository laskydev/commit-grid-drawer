package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "commit-grid",
	Short: "Dibuja patrones en tu contribution graph con commits diarios",
	Long:  "CLI para programar commits diarios y 'dibujar' patrones en GitHub, con UX moderna.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Aquí podrías agregar flags globales si los necesitas
}
