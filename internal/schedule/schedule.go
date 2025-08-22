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

// marcador único para identificar entradas de cron de esta app
const cronMarker = "# COMMIT_GRID_DRAW"

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
	if util.IsLinux() {
		return disableCron()
	}
	if util.IsMac() {
		return disableLaunchd()
	}
	return fmt.Errorf("OS no soportado")
}

func Status() (string, error) {
	if util.IsLinux() {
		return cronStatus()
	}
	if util.IsMac() {
		return launchdStatus()
	}
	return "", fmt.Errorf("OS no soportado")
}

// --- Linux: cron ---
func enableCron(bin string, h, m int) error {
	line := fmt.Sprintf("%d %d * * * %s run >> ~/.local/state/commit-grid-draw/commit-grid.log 2>&1 %s", m, h, bin, cronMarker)
	// leer crontab actual
	cur, _ := exec.Command("crontab", "-l").Output()
	if bytes.Contains(cur, []byte(cronMarker)) {
		// reemplazar
		out := replaceCronLine(string(cur), line)
		cmd := exec.Command("crontab", "-")
		cmd.Stdin = bytes.NewBufferString(out)
		return cmd.Run()
	}
	// agregar
	buf := bytes.NewBuffer(cur)
	if len(cur) > 0 && cur[len(cur)-1] != '\n' {
		buf.WriteByte('\n')
	}
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
	if err != nil {
		return "sin crontab o error consultando", nil
	}
	if bytes.Contains(cur, []byte(cronMarker)) {
		return "habilitado (cron)", nil
	}
	return "deshabilitado (cron)", nil
}

func replaceCronLine(cur, newLine string) string {
	return dropCronLines(cur) + newLine + "\n"
}
func dropCronLines(cur string) string {
	lines := bytes.Split([]byte(cur), []byte("\n"))
	out := bytes.Buffer{}
	for _, ln := range lines {
		if !bytes.Contains(ln, []byte(cronMarker)) && len(bytes.TrimSpace(ln)) > 0 {
			out.Write(ln)
			out.WriteByte('\n')
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
	if bytes.Contains(out, []byte("com.commitgrid.draw")) {
		return "habilitado (launchd)", nil
	}
	return "deshabilitado (launchd)", nil
}
