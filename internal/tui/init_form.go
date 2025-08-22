package tui

import (
	"commit-grid/internal/config"
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
)

func Onboard() (*config.Config, error) {
	var repo, user, email, tz string
	var hourStr, minuteStr, intensityStr string
	strategy := "fixed"

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Ruta del repo Git").Value(&repo).Prompt("📁 "),
			huh.NewInput().Title("Git user.name").Value(&user).Prompt("👤 "),
			huh.NewInput().Title("Git user.email").Value(&email).Prompt("✉️ "),
			huh.NewInput().Title("Zona horaria (IANA)").Value(&tz).Placeholder("America/Monterrey"),
			huh.NewSelect[string]().
				Title("Estrategia de intensidad").
				Options(huh.NewOptions("fixed","random","pattern")...).
				Value(&strategy),
			huh.NewInput().Title("Intensidad fija (si aplica)").Prompt("🔥 ").Value(&intensityStr),
			huh.NewInput().Title("Hora (0-23)").Prompt("🕒 ").Value(&hourStr),
			huh.NewInput().Title("Minuto (0-59)").Prompt("🕒 ").Value(&minuteStr),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil { return nil, err }

	if repo == "" || user == "" || email == "" {
		return nil, fmt.Errorf("campos obligatorios vacíos")
	}

	hour, _ := strconv.Atoi(hourStr)
	minute, _ := strconv.Atoi(minuteStr)
	intensity, _ := strconv.Atoi(intensityStr)

	cfg := &config.Config{
		RepoPath: repo, GitUser: user, GitEmail: email,
		Timezone: tz, Hour24: hour, Minute: minute,
		IntensityStrategy: strategy, IntensityValue: intensity,
		PatternFile: "data/pattern.csv",
	}
	return cfg, nil
}
