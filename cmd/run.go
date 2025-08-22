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
		fmt.Println("🔄 Iniciando commit-grid...")

		// Cargar configuración
		fmt.Println("📁 Cargando configuración...")
		c, _, err := config.Load()
		if err != nil {
			return fmt.Errorf("❌ Error al cargar configuración: %w\n💡 Ejecuta 'commit-grid init' para configurar la aplicación", err)
		}
		fmt.Printf("✅ Configuración cargada desde: %s\n", c.RepoPath)

		// Validar repositorio Git
		fmt.Println("🔍 Validando repositorio Git...")
		if err := gitutil.ValidateRepo(c.RepoPath); err != nil {
			return fmt.Errorf("❌ Error en validación del repositorio: %w", err)
		}
		fmt.Println("✅ Repositorio Git válido")

		// Crear archivo del día
		fmt.Println("📝 Creando archivo del día...")
		if err := pattern.TouchToday(c.RepoPath); err != nil {
			return fmt.Errorf("❌ Error al crear archivo del día: %w", err)
		}
		fmt.Println("✅ Archivo del día creado")

		// Calcular intensidad
		intensity := pattern.IntensityToday(pattern.Strategy(c.IntensityStrategy), c.IntensityValue)
		fmt.Printf("🎯 Intensidad del día: %d commits\n", intensity)

		// Ejecutar commits
		for i := 1; i <= intensity; i++ {
			msg := fmt.Sprintf("grid: %s (%d/%d)", time.Now().UTC().Format("2006-01-02:15:04:05"), i, intensity)
			fmt.Printf("🔄 Haciendo commit %d/%d: %s\n", i, intensity, msg)

			if err := gitutil.CommitPush(c.RepoPath, c.GitUser, c.GitEmail, msg); err != nil {
				return fmt.Errorf("❌ Error en commit %d/%d: %w", i, intensity, err)
			}
			fmt.Printf("✅ Commit %d/%d completado\n", i, intensity)
		}

		fmt.Printf("🎉 ¡Completado! Se realizaron %d commits exitosamente\n", intensity)
		return nil
	},
}

func init() { rootCmd.AddCommand(runCmd) }
