# Commit Grid Drawer 🎨📊

Commit Grid Drawer es una herramienta CLI multiplataforma que automatiza commits diarios a GitHub para "dibujar" patrones personalizados en tu gráfico de contribuciones.
Proporciona una experiencia de onboarding moderna con TUI, programación flexible y múltiples estrategias para la intensidad de commits.

## ✨ Características

- **Commits automáticos diarios** a tu repositorio GitHub
- **Soporte para dibujar patrones** (fijo, aleatorio o basado en CSV)
- **Onboarding interactivo (TUI)** para configurar repo, usuario, zona horaria y programación
- **Programación multiplataforma**:
  - Linux → cron
  - macOS → launchd
- **Intensidad configurable** (# de commits por día)
- **Ligero** (binario Go, sin daemons, inicio instantáneo)
- **CLI amigable** con UX moderna (stack Charmbracelet)

## 🚀 Inicio Rápido

### 1) Clonar y construir

```bash
# Clonar el repositorio
git clone https://github.com/laskydev/commit-grid-drawer.git
cd commit-grid-drawer

# Construir el binario
go build -o commit-grid .
```

### 2) Ejecutar onboarding

```bash
./commit-grid init
```

### 3) Activar el programador

```bash
./commit-grid enable
```

### 4) Verificar estado

```bash
./commit-grid status
```

### 5) Probar manualmente

```bash
./commit-grid run
```

### 6) Ver configuración

```bash
./commit-grid config get
```

## ⚙️ Configuración

El archivo de configuración se guarda en `~/.config/commit-grid-draw/config.yaml`:

```yaml
repo_path: "./drawing"                    # Ruta al repositorio Git
git_user: "tu-usuario"                    # Nombre de usuario para Git
git_email: "tu-email@example.com"        # Email para Git
timezone: "America/Monterrey"             # Zona horaria (opcional)
hour_24: 10                              # Hora de ejecución (0-23)
minute: 0                                 # Minuto de ejecución (0-59)
intensity_strategy: "fixed"               # Estrategia: fixed | random | pattern
intensity_value: 1                        # Número de commits por día (para fixed)
pattern_file: "data/pattern.csv"          # Archivo de patrón (para pattern)
```

## 🧱 Stack Tecnológico

- **Lenguaje**: Go 1.22+
- **CLI**: spf13/cobra
- **TUI**: bubbletea, bubbles, lipgloss, glamour, huh
- **Configuración**: YAML en `~/.config/commit-grid-draw/config.yaml`
- **Programador**: cron (Linux), launchd (macOS)
- **Logs**:
  - Linux → `~/.local/state/commit-grid-draw/commit-grid.log`
  - macOS → `~/Library/Logs/commit-grid.log`

## 🕒 Programación

### Linux (cron)
```bash
0 10 * * * /ruta/al/binario/commit-grid run >> ~/.local/state/commit-grid-draw/commit-grid.log 2>&1
```

### macOS (launchd)
El archivo se crea automáticamente en `~/Library/LaunchAgents/com.commitgrid.draw.plist`

## 🤖 Cómo Funciona

1. **Carga** la configuración del usuario y zona horaria
2. **Determina** la intensidad del día (número de commits)
3. **Asegura** que el repo esté limpio
4. **Actualiza** `data/grid.csv` con la entrada de hoy
5. **Hace N commits** con mensajes como:
   ```
   grid: 2025-08-21 (1/3)
   ```
6. **Hace push** de los commits al remoto, actualizando tu gráfico de contribuciones

## 📋 Comandos Disponibles

- `commit-grid init` - Asistente interactivo de configuración
- `commit-grid enable` - Activa el programador diario
- `commit-grid disable` - Desactiva el programador diario
- `commit-grid status` - Muestra el estado del programador
- `commit-grid run` - Ejecuta la tarea del día manualmente
- `commit-grid config get` - Lee la configuración actual
- `commit-grid completion` - Genera script de autocompletado

## ⚠️ Solución de Problemas

### Error "exit status 128"
Este error típicamente indica un problema con Git. Verifica:

1. **El repositorio existe** y es válido
2. **Tienes permisos** para hacer push al remoto
3. **El remoto está configurado** correctamente
4. **Tu autenticación Git** está funcionando

### Cambiar usuario de Git
Si necesitas cambiar el usuario de Git configurado:

1. Edita manualmente `~/.config/commit-grid-draw/config.yaml`
2. Cambia `git_user` y `git_email`
3. O ejecuta `./commit-grid init` para reconfigurar

## 🧪 Notas de Calidad

- **Idempotente**: habilitar reemplaza entradas previas de cron/launchd
- **Seguro**: solo hace commits dentro del repo elegido
- **Portable**: binario Go estático, sin CGO
- **Logs**: toda la actividad queda registrada para debugging

## 📌 Descargo de Responsabilidad

⚠️ **Las estrategias de alta intensidad generan múltiples commits por día y pueden considerarse spam.** Úsalas responsablemente para mantener tu gráfico divertido y significativo.

## 📄 Licencia

![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)