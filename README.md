# Commit Grid Drawer 🎨📊

Commit Grid Drawer is a cross-platform CLI tool that automates daily commits to GitHub to "draw" custom patterns on your contribution graph.
It provides a modern TUI onboarding experience, flexible scheduling, and multiple strategies for commit intensity.

## ✨ Features

- **Daily automated commits** to your GitHub repository
- **Pattern drawing support** (fixed, random, or CSV-based)
- **Interactive onboarding (TUI)** to configure repo, user, timezone, and schedule
- **Cross-platform scheduling**:
  - Linux → cron
  - macOS → launchd
- **Configurable intensity** (# of commits per day)
- **Lightweight** (Go binary, no daemons, instant startup)
- **User-friendly CLI** with modern UX (Charmbracelet stack)

## 🚀 Quick Start

### 1) Clone and build

```bash
# Clone the repository
git clone https://github.com/laskydev/commit-grid-drawer.git
cd commit-grid-drawer

# Build the binary
go build -o commit-grid .
```

### 2) Run onboarding

```bash
./commit-grid init
```

### 3) Enable scheduler

```bash
./commit-grid enable
```

### 4) Check status

```bash
./commit-grid status
```

### 5) Test manually

```bash
./commit-grid run
```

### 6) View configuration

```bash
./commit-grid config get
```

## ⚙️ Configuration

The configuration file is saved in `~/.config/commit-grid-draw/config.yaml`:

```yaml
repo_path: "./drawing" # Path to Git repository
git_user: "your-username" # Git username
git_email: "your-email@example.com" # Git email
timezone: "America/Monterrey" # Timezone (optional)
hour_24: 10 # Execution hour (0-23)
minute: 0 # Execution minute (0-59)
intensity_strategy: "fixed" # Strategy: fixed | random | pattern
intensity_value: 1 # Number of commits per day (for fixed)
pattern_file: "data/pattern.csv" # Pattern file (for pattern)
```

## 🧱 Tech Stack

- **Language**: Go 1.22+
- **CLI**: spf13/cobra
- **TUI**: bubbletea, bubbles, lipgloss, glamour, huh
- **Configuration**: YAML in `~/.config/commit-grid-draw/config.yaml`
- **Scheduler**: cron (Linux), launchd (macOS)
- **Logs**:
  - Linux → `~/.local/state/commit-grid-draw/commit-grid.log`
  - macOS → `~/Library/Logs/commit-grid.log`

## 🕒 Scheduling

### Linux (cron)

```bash
0 10 * * * /path/to/binary/commit-grid run >> ~/.local/state/commit-grid-draw/commit-grid.log 2>&1
```

### macOS (launchd)

The file is automatically created in `~/Library/LaunchAgents/com.commitgrid.draw.plist`

## 🤖 How It Works

1. **Loads** user configuration and timezone
2. **Determines** today's intensity (number of commits)
3. **Ensures** the repo is clean
4. **Updates** `data/grid.csv` with today's entry
5. **Makes N commits** with messages like:
   ```
   grid: 2025-08-21 (1/3)
   ```
6. **Pushes** commits to remote, updating your contribution graph

## 📋 Available Commands

- `commit-grid init` - Interactive configuration wizard
- `commit-grid enable` - Enables daily scheduler
- `commit-grid disable` - Disables daily scheduler
- `commit-grid status` - Shows scheduler status
- `commit-grid run` - Manually executes daily task
- `commit-grid config get` - Reads current configuration
- `commit-grid completion` - Generates autocompletion script

## ⚠️ Troubleshooting

### Error "exit status 128"

This error typically indicates a Git problem. Check:

1. **The repository exists** and is valid
2. **You have permissions** to push to remote
3. **The remote is configured** correctly
4. **Your Git authentication** is working

### Change Git user

If you need to change the configured Git user:

1. Manually edit `~/.config/commit-grid-draw/config.yaml`
2. Change `git_user` and `git_email`
3. Or run `./commit-grid init` to reconfigure

## 🧪 Quality Notes

- **Idempotent**: enabling replaces previous cron/launchd entries
- **Safe**: only commits within your chosen repo
- **Portable**: static Go binary, no CGO
- **Logs**: all activity is recorded for debugging

## 📌 Disclaimer

⚠️ **High-intensity strategies generate multiple commits per day and may be considered spam.** Use them responsibly to keep your graph fun and meaningful.

## 📄 License

![License](https://github.com/laskydev/commit-grid-drawer/blob/main/LICENSE)
