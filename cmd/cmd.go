package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Global vars
var (
	Debug bool
)

// Root represents the base command when called without any subcommands
var Root = &cobra.Command{
	Use:   "ncwtester",
	Short: "Show help for ncwtester commands.",
	Long: `
Ncwtester provides a suite of morse code tools.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := Root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	pflag.BoolVarP(&Debug, "verbose", "v", false, "Verbose debugging")
}

// Debugf writes format to stderr if -v was set
func Debugf(format string, a ...interface{}) {
	if !Debug {
		return
	}
	fmt.Fprintf(os.Stderr, format, a...)
	os.Stderr.WriteString("\n")
}
