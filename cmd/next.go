package cmd

import (
	"fmt"
	"time"

	"commit-grid/internal/config"
	"commit-grid/internal/schedule"

	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Muestra cuándo será la próxima sincronización",
	Long:  "Muestra información sobre la próxima ejecución programada y el estado actual del scheduler",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Consultando próxima sincronización...")

		// Cargar configuración
		c, _, err := config.Load()
		if err != nil {
			return fmt.Errorf("❌ Error al cargar configuración: %w\n💡 Ejecuta 'commit-grid init' para configurar la aplicación", err)
		}

		// Obtener estado del scheduler
		status, err := schedule.Status()
		if err != nil {
			return fmt.Errorf("❌ Error al consultar estado del scheduler: %w", err)
		}

		fmt.Printf("📅 Estado del scheduler: %s\n", status)

		if status == "deshabilitado" || status == "deshabilitado (cron)" || status == "deshabilitado (launchd)" {
			fmt.Println("\n⚠️  El scheduler está deshabilitado")
			fmt.Println("💡 Para habilitarlo, ejecuta: commit-grid enable")
			return nil
		}

		// Calcular próxima ejecución basada en la hora configurada
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), c.Hour24, c.Minute, 0, 0, now.Location())

		// Si ya pasó la hora de hoy, la próxima será mañana
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		// Calcular tiempo restante
		timeUntil := nextRun.Sub(now)

		fmt.Printf("⏰ Hora configurada: %02d:%02d\n", c.Hour24, c.Minute)
		fmt.Printf("🔄 Próxima ejecución: %s\n", nextRun.Format("2006-01-02 15:04:05"))
		fmt.Printf("⏳ Tiempo restante: %s\n", formatDuration(timeUntil))

		// Mostrar información adicional sobre la estrategia
		fmt.Printf("\n📊 Estrategia de intensidad: %s\n", c.IntensityStrategy)
		if c.IntensityStrategy == "fixed" {
			fmt.Printf("🎯 Commits por día: %d\n", c.IntensityValue)
		} else if c.IntensityStrategy == "random" {
			fmt.Printf("🎯 Commits por día: 1-4 (aleatorio)\n")
		}

		// Mostrar última ejecución (si existe)
		fmt.Println("\n📝 Para ejecutar manualmente ahora: commit-grid run")

		return nil
	},
}

// formatDuration formatea una duración de manera legible
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f segundos", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%d minutos, %d segundos", minutes, seconds)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%d horas, %d minutos", hours, minutes)
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%d días, %d horas", days, hours)
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
