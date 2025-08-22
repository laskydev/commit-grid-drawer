Project: Commit Grid Draw

Features:
    -Daily commit to GitHub
    -Draw patterns


Requirimients (spanish - TODO english versio)

🎯 Objetivo

Un CLI commit-grid que:

commit-grid init → Onboarding interactivo (TUI) para configurar:

Repo Git a usar

Nombre/email del bot

Zona horaria

Hora diaria de ejecución

Estrategia de “intensidad” (n° de commits por día)

commit-grid enable / disable → instala/quita cron (Linux) o launchd (macOS).

commit-grid status → muestra si está activo, próximo run y últimos logs.

commit-grid run → ejecuta la lógica del día (para probar manualmente).

commit-grid config get/set → lee/escribe config.

commit-grid pattern → utilidades para patrones (p. ej. setear la intensidad de hoy).

RAM: Go + proceso corto (sin daemon).
Arranque: instantáneo (binario único).
Soporte: Linux (cron), macOS (launchd; si prefieres, también cron).

🧱 Stack

Go 1.22+

CLI: spf13/cobra

TUI/estilo: charmbracelet/bubbletea, bubbles, lipgloss, glamour, huh (forms)

Config: os.UserConfigDir() → ~/.config/commit-grid-draw/config.yaml

Git: shell out a git (más liviano que go-git)

Logs: ~/.local/state/commit-grid-draw/commit-grid.log (Linux) / ~/Library/Logs/commit-grid.log (macOS)

📦 Estructura
commit-grid/
  cmd/
    root.go
    init.go
    enable.go
    disable.go
    status.go
    run.go
    config.go
    pattern.go
  internal/
    config/config.go
    schedule/schedule.go
    gitutil/git.go
    pattern/pattern.go
    tui/init_form.go
    util/platform.go
  main.go
  go.mod

⚙️ Config (YAML)
repo_path: "/home/user/projects/commit-art"
git_user: "grid-bot"
git_email: "grid-bot@users.noreply.github.com"
timezone: "America/Monterrey"
hour_24: 05
minute: 00
intensity_strategy: "fixed"  # fixed|pattern|random
intensity_value: 3           # para fixed
pattern_file: "data/pattern.csv"

🕒 Scheduling
Linux (cron)

Entrada en crontab del usuario:

0 5 * * * /usr/local/bin/commit-grid run >> ~/.local/state/commit-grid-draw/commit-grid.log 2>&1

macOS (launchd)

~/Library/LaunchAgents/com.commitgrid.draw.plist:

ProgramArguments: ["/usr/local/bin/commit-grid","run"]

StartCalendarInterval: { Hour=5, Minute=0 }

StandardOutPath/StandardErrorPath a ~/Library/Logs/commit-grid.log

Cargar con: launchctl load -w ~/Library/LaunchAgents/com.commitgrid.draw.plist

El CLI hace esto automáticamente en enable, detectando OS.

🤖 Lógica de “run”

Carga config y TZ.

Decide intensidad del día (1..N) según estrategia.

Asegura repo_path limpio (pull opcional).

Actualiza un archivo (p.ej. data/grid.csv) con la entrada de hoy.

Hace N commits:

git add -A

git -c user.name=... -c user.email=... commit -m "grid: YYYY-MM-DD (#i/N)"

git push

Nota: Intensidades altas generan varios commits; úsalo con moderación.

🚀 Código base (recorta/pega)
main.go
package main

import "commit-grid/cmd"

func main() { cmd.Execute() }

cmd/root.go
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

internal/config/config.go
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RepoPath          string `yaml:"repo_path"`
	GitUser           string `yaml:"git_user"`
	GitEmail          string `yaml:"git_email"`
	Timezone          string `yaml:"timezone"`
	Hour24            int    `yaml:"hour_24"`
	Minute            int    `yaml:"minute"`
	IntensityStrategy string `yaml:"intensity_strategy"`
	IntensityValue    int    `yaml:"intensity_value"`
	PatternFile       string `yaml:"pattern_file"`
}

func defaultPaths() (cfgPath, stateDir string, err error) {
	cfgBase, err := os.UserConfigDir()
	if err != nil { return "", "", err }
	stateBase, err := os.UserHomeDir()
	if err != nil { return "", "", err }
	cfgDir := filepath.Join(cfgBase, "commit-grid-draw")
	stateDir = filepath.Join(stateBase, ".local", "state", "commit-grid-draw")
	_ = os.MkdirAll(cfgDir, 0o755)
	_ = os.MkdirAll(stateDir, 0o755)
	return filepath.Join(cfgDir, "config.yaml"), stateDir, nil
}

func Load() (*Config, string, error) {
	cfgFile, _, err := defaultPaths()
	if err != nil { return nil, "", err }
	b, err := os.ReadFile(cfgFile)
	if err != nil { return nil, "", err }
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil { return nil, "", err }
	return &c, cfgFile, nil
}

func Save(c *Config) (string, error) {
	cfgFile, _, err := defaultPaths()
	if err != nil { return "", err }
	b, err := yaml.Marshal(c)
	if err != nil { return "", err }
	return cfgFile, os.WriteFile(cfgFile, b, 0o644)
}

func Exists() bool {
	cfgFile, _, err := defaultPaths()
	if err != nil { return false }
	_, err = os.Stat(cfgFile)
	return err == nil
}

func Debug(c *Config) string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

internal/util/platform.go
package util

import "runtime"

func IsMac() bool   { return runtime.GOOS == "darwin" }
func IsLinux() bool { return runtime.GOOS == "linux" }

internal/schedule/schedule.go
package schedule

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"commit-grid/internal/util"
)

const plistTmpl = `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC '-//Apple//DTD PLIST 1.0//EN' 'http://www.apple.com/DTDs/PropertyList-1.0.dtd'>
<plist version='1.0'><dict>
  <key>Label</key><string>com.commitgrid.draw</string>
  <key>ProgramArguments</key>
  <array><string>{{.Bin}}</string><string>run</string></array>
  <key>StartCalendarInterval</key>
  <dict><key>Hour</key><integer>{{.Hour}}</integer><key>Minute</key><integer>{{.Minute}}</integer></dict>
  <key>StandardOutPath</key><string>{{.Log}}</string>
  <key>StandardErrorPath</key><string>{{.Log}}</string>
  <key>RunAtLoad</key><true/>
</dict></plist>`

func Enable(binPath string, hour, minute int) error {
	if util.IsLinux() {
		return enableCron(binPath, hour, minute)
	}
	if util.IsMac() {
		return enableLaunchd(binPath, hour, minute)
	}
	return fmt.Errorf("OS no soportado")
}

func Disable() error {
	if util.IsLinux() { return disableCron() }
	if util.IsMac()   { return disableLaunchd() }
	return fmt.Errorf("OS no soportado")
}

func Status() (string, error) {
	if util.IsLinux() { return cronStatus() }
	if util.IsMac()   { return launchdStatus() }
	return "", fmt.Errorf("OS no soportado")
}

// --- Linux: cron ---
func enableCron(bin string, h, m int) error {
	line := fmt.Sprintf("%d %d * * * %s run >> ~/.local/state/commit-grid-draw/commit-grid.log 2>&1", m, h, bin)
	// leer crontab actual
	cur, _ := exec.Command("crontab", "-l").Output()
	if bytes.Contains(cur, []byte("commit-grid run")) {
		// reemplazar
		out := replaceCronLine(string(cur), line)
		cmd := exec.Command("crontab", "-")
		cmd.Stdin = bytes.NewBufferString(out)
		return cmd.Run()
	}
	// agregar
	buf := bytes.NewBuffer(cur)
	if len(cur) > 0 && cur[len(cur)-1] != '\n' { buf.WriteByte('\n') }
	buf.WriteString(line + "\n")
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = buf
	return cmd.Run()
}

func disableCron() error {
	cur, _ := exec.Command("crontab", "-l").Output()
	out := dropCronLines(string(cur))
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = bytes.NewBufferString(out)
	return cmd.Run()
}

func cronStatus() (string, error) {
	cur, err := exec.Command("crontab", "-l").Output()
	if err != nil { return "sin crontab o error consultando", nil }
	if bytes.Contains(cur, []byte("commit-grid run")) { return "habilitado (cron)", nil }
	return "deshabilitado (cron)", nil
}

func replaceCronLine(cur, newLine string) string {
	return dropCronLines(cur) + newLine + "\n"
}
func dropCronLines(cur string) string {
	lines := bytes.Split([]byte(cur), []byte("\n"))
	out := bytes.Buffer{}
	for _, ln := range lines {
		if !bytes.Contains(ln, []byte("commit-grid run")) && len(bytes.TrimSpace(ln)) > 0 {
			out.Write(ln); out.WriteByte('\n')
		}
	}
	return out.String()
}

// --- macOS: launchd ---
func enableLaunchd(bin string, h, m int) error {
	home, _ := os.UserHomeDir()
	log := filepath.Join(home, "Library", "Logs", "commit-grid.log")
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.commitgrid.draw.plist")
	_ = os.MkdirAll(filepath.Dir(plistPath), 0o755)
	var buf bytes.Buffer
	t := template.Must(template.New("p").Parse(plistTmpl))
	if err := t.Execute(&buf, map[string]any{"Bin": bin, "Hour": h, "Minute": m, "Log": log}); err != nil {
		return err
	}
	if err := os.WriteFile(plistPath, buf.Bytes(), 0o644); err != nil {
		return err
	}
	exec.Command("launchctl", "unload", plistPath).Run()
	return exec.Command("launchctl", "load", "-w", plistPath).Run()
}

func disableLaunchd() error {
	home, _ := os.UserHomeDir()
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.commitgrid.draw.plist")
	_ = exec.Command("launchctl", "unload", plistPath).Run()
	return os.Remove(plistPath)
}

func launchdStatus() (string, error) {
	out, _ := exec.Command("launchctl", "list").Output()
	if bytes.Contains(out, []byte("com.commitgrid.draw")) { return "habilitado (launchd)", nil }
	return "deshabilitado (launchd)", nil
}

internal/gitutil/git.go
package gitutil

import (
	"fmt"
	"os/exec"
)

func run(repo string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	return cmd.Run()
}

func CommitPush(repo, user, email, msg string) error {
	if err := run(repo, "add", "-A"); err != nil { return err }
	if err := run(repo, "-c", "user.name="+user, "-c", "user.email="+email, "commit", "-m", msg); err != nil {
		// si no hay cambios, git devuelve error: ignorarlo
		return nil
	}
	return run(repo, "push")
}

internal/pattern/pattern.go
package pattern

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Strategy string
const (
	Fixed   Strategy = "fixed"
	Random  Strategy = "random"
	Pattern Strategy = "pattern"
)

func IntensityToday(str Strategy, fixedVal int) int {
	switch str {
	case Fixed:
		if fixedVal < 1 { return 1 }
		return fixedVal
	case Random:
		return 1 + rand.Intn(4) // 1..4
	default:
		return 1 // TODO: leer de archivo pattern.csv si lo deseas
	}
}

func TouchToday(repo string) error {
	_ = os.MkdirAll(filepath.Join(repo, "data"), 0o755)
	f := filepath.Join(repo, "data", "grid.csv")
	now := time.Now().UTC().Format("2006-01-02")
	entry := fmt.Sprintf("%s,level=1\n", now)
	fh, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil { return err }
	defer fh.Close()
	_, err = fh.WriteString(entry)
	return err
}

cmd/run.go
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

cmd/enable.go
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

internal/tui/init_form.go (onboarding con Charmbracelet Huh)
package tui

import (
	"commit-grid/internal/config"
	"fmt"

	"github.com/charmbracelet/huh"
)

func Onboard() (*config.Config, error) {
	var repo, user, email, tz string
	var hour, minute, intensity int
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
			huh.NewInput().Title("Intensidad fija (si aplica)").Prompt("🔥 ").Value(&intensity),
			huh.NewInput().Title("Hora (0-23)").Prompt("🕒 ").Value(&hour),
			huh.NewInput().Title("Minuto (0-59)").Prompt("🕒 ").Value(&minute),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil { return nil, err }

	if repo == "" || user == "" || email == "" {
		return nil, fmt.Errorf("campos obligatorios vacíos")
	}
	cfg := &config.Config{
		RepoPath: repo, GitUser: user, GitEmail: email,
		Timezone: tz, Hour24: hour, Minute: minute,
		IntensityStrategy: strategy, IntensityValue: intensity,
		PatternFile: "data/pattern.csv",
	}
	return cfg, nil
}

cmd/init.go
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

cmd/config.go
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

🏁 Uso
# 1) Compila
go mod init commit-grid
go get github.com/spf13/cobra github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/charmbracelet/huh gopkg.in/yaml.v3
go build -o commit-grid .

# 2) Onboarding
./commit-grid init

# 3) Activa scheduler
./commit-grid enable

# 4) Ver estado
./commit-grid status

# 5) Probar manual
./commit-grid run

# 6) Ver/editar config
./commit-grid config get

🧪 Notas de calidad

Idempotencia: enable reemplaza la entrada previa de cron/launchd.

Logs: quedan en la ruta indicada; status se puede ampliar para leer “último run OK/ERROR”.

Seguridad: usa tu git ya autenticado (SSH o PAT). Puedes añadir GIT_SSH_COMMAND si quieres llaves dedicadas.

Portabilidad: sin CGO, binario estático (en Linux con -tags netgo si deseas).

Animaciones: huh ya te da una UI suave; si quieres spinners durante run, mete bubbles/spinner en cmd/run.go.

Si quieres, te lo empaqueto con un Makefile, fórmulas de Homebrew e instalador .deb. ¿Quieres que te lo deje listo con un patrón ejemplo (7×53) y comandos para “pintar” semanas completas?
