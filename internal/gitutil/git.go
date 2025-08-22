package gitutil

import (
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
