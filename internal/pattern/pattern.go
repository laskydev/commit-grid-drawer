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
