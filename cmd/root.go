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
		fmt.Printf("\n❌ Error: %s\n\n", err)
		
		// Proporcionar sugerencias basadas en el tipo de error
		if err.Error() == "exit status 128" {
			fmt.Println("💡 Este error sugiere un problema con Git. Verifica:")
			fmt.Println("   • Que el repositorio esté inicializado correctamente")
			fmt.Println("   • Que tengas permisos para hacer push")
			fmt.Println("   • Que la configuración de Git sea correcta")
			fmt.Println("   • Que el repositorio remoto esté configurado")
		}
		
		fmt.Println("🔧 Para más ayuda, ejecuta: commit-grid --help")
		os.Exit(1)
	}
}

func init() {
	// Aquí podrías agregar flags globales si los necesitas
}
