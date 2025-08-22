package gitutil

import (
	"fmt"
	"os/exec"
	"strings"
)

// ValidateRepo verifica que el repositorio esté en un estado válido
func ValidateRepo(repo string) error {
	// Verificar que el directorio existe y es un repositorio Git
	if err := run(repo, "rev-parse", "--git-dir"); err != nil {
		return fmt.Errorf("el directorio '%s' no es un repositorio Git válido: %w", repo, err)
	}

	// Verificar que hay un remoto configurado
	if err := run(repo, "remote", "-v"); err != nil {
		return fmt.Errorf("no hay repositorio remoto configurado: %w", err)
	}

	// Verificar que el branch actual esté sincronizado
	if err := run(repo, "fetch"); err != nil {
		return fmt.Errorf("error al hacer fetch del repositorio remoto: %w", err)
	}

	// Verificar si el branch actual tiene upstream configurado
	if err := run(repo, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"); err != nil {
		// No hay upstream configurado, configurarlo automáticamente
		fmt.Println("⚠️  El branch actual no tiene upstream configurado. Configurando automáticamente...")
		
		// Obtener el nombre del branch actual
		branchOutput, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
		if err != nil {
			return fmt.Errorf("error al obtener el nombre del branch actual: %w", err)
		}
		branchName := strings.TrimSpace(string(branchOutput))
		
		// Configurar el upstream
		if err := run(repo, "push", "--set-upstream", "origin", branchName); err != nil {
			return fmt.Errorf("error al configurar upstream para el branch %s: %w", branchName, err)
		}
		fmt.Printf("✅ Upstream configurado para el branch %s\n", branchName)
	}

	return nil
}

func run(repo string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo

	// Capturar la salida de error para proporcionar información más detallada
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Si hay salida, incluirla en el error para más contexto
		if len(output) > 0 {
			return fmt.Errorf("git %s falló: %s\nSalida: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
		}
		return fmt.Errorf("git %s falló: %s", strings.Join(args, " "), err)
	}
	return nil
}

func CommitPush(repo, user, email, msg string) error {
	if err := run(repo, "add", "-A"); err != nil {
		return fmt.Errorf("error al agregar archivos: %w", err)
	}

	if err := run(repo, "-c", "user.name="+user, "-c", "user.email="+email, "commit", "-m", msg); err != nil {
		// Verificar si el error es porque no hay cambios
		if strings.Contains(err.Error(), "nothing to commit") || strings.Contains(err.Error(), "no changes added to commit") {
			// No es realmente un error, solo no hay cambios
			return nil
		}
		return fmt.Errorf("error al hacer commit: %w", err)
	}

	if err := run(repo, "push"); err != nil {
		return fmt.Errorf("error al hacer push: %w", err)
	}

	return nil
}
