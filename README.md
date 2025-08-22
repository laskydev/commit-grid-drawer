Commit Grid Draw 🎨📊

Commit Grid Draw is a cross-platform CLI tool that automates daily commits to GitHub in order to draw custom patterns on your contribution graph.
It provides a modern TUI onboarding experience, flexible scheduling, and multiple strategies for commit intensity.

✨ Features

Daily automated commits to your GitHub repository.

Pattern drawing support (fixed, random, or CSV-driven).

Interactive onboarding (TUI) to configure repo, user, timezone, and schedule.

Cross-platform scheduling:

Linux → cron

macOS → launchd (or cron)

Configurable intensity (# of commits per day).

Lightweight (Go binary, no daemons, instant startup).

Human-friendly CLI with modern UX (Charmbracelet stack).

🚀 Quick Start
# 1) Build
go mod init commit-grid
go get github.com/spf13/cobra github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss github.com/charmbracelet/huh gopkg.in/yaml.v3
go build -o commit-grid .

# 2) Run onboarding
./commit-grid init

# 3) Enable scheduler
./commit-grid enable

# 4) Check status
./commit-grid status

# 5) Run manually (test)
./commit-grid run

# 6) View config
./commit-grid config get

⚙️ Configuration (YAML)

Example ~/.config/commit-grid-draw/config.yaml:

repo_path: "/home/user/projects/commit-art"
git_user: "grid-bot"
git_email: "grid-bot@users.noreply.github.com"
timezone: "America/Monterrey"
hour_24: 05
minute: 00
intensity_strategy: "fixed"   # fixed | pattern | random
intensity_value: 3
pattern_file: "data/pattern.csv"

🧱 Tech Stack

Language: Go 1.22+

CLI: spf13/cobra

TUI: bubbletea
, bubbles
, lipgloss
, glamour
, huh

Config: YAML in ~/.config/commit-grid-draw/config.yaml

Scheduler: cron (Linux), launchd (macOS)

Logging:

Linux → ~/.local/state/commit-grid-draw/commit-grid.log

macOS → ~/Library/Logs/commit-grid.log

🕒 Scheduling
Linux (cron)
0 5 * * * /usr/local/bin/commit-grid run >> ~/.local/state/commit-grid-draw/commit-grid.log 2>&1

macOS (launchd)

~/Library/LaunchAgents/com.commitgrid.draw.plist:

<dict>
  <key>ProgramArguments</key>
  <array><string>/usr/local/bin/commit-grid</string><string>run</string></array>
  <key>StartCalendarInterval</key>
  <dict><key>Hour</key><integer>5</integer><key>Minute</key><integer>0</integer></dict>
</dict>

🤖 How It Works

Loads user config and timezone.

Determines today’s intensity (number of commits).

Ensures repo is clean (optional pull).

Updates data/grid.csv with today’s entry.

Makes N commits with messages like:

grid: 2025-08-21 (1/3)


Pushes commits to remote, updating your GitHub contribution graph.

🧪 Quality Notes

Idempotent: enabling replaces old cron/launchd entries.

Safe: commits only inside your chosen repo.

Portable: static Go binary, no CGO.

Logs: all activity recorded for debugging.

📌 Disclaimer

⚠️ High-intensity strategies generate multiple commits per day and may be considered spammy. Use responsibly to keep your graph fun and meaningful.

![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)