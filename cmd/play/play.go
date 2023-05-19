package play

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ncw/cwtool/cmd"
	"github.com/ncw/cwtool/cmd/cwflags"
	"github.com/ncw/cwtool/cw"
	"github.com/spf13/cobra"
)

var (
	file  string
	stdin bool
)

// subCmd represents the rss ommand
var subCmd = &cobra.Command{
	Use:   "play",
	Short: "Play morse code from the command line or file",
	Long: strings.ReplaceAll(`

This plays morse code from the command line or from a file with the
|--file| flag or from stdin with the |--stdin| flag.

`, "|", "`"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(args)
	},
}

func init() {
	cmd.Root.AddCommand(subCmd)
	flags := subCmd.Flags()
	cwflags.Add(flags)
	flags.StringVarP(&file, "file", "", "", "File to play Morse from (optional)")
	flags.BoolVarP(&stdin, "stdin", "", false, "If set play Morse from stdin")
}

var (
	removeInvalid = regexp.MustCompile(`[^A-Za-z0-9,.: /]+`)
)

// Simplify and play the string
func play(cw cw.CW, s string) {
	// Remove all unknown characters
	s = removeInvalid.ReplaceAllString(s, "")

	// Replace `:` with more CW friendly ` =`
	s = strings.ReplaceAll(s, ":", " =")

	fmt.Println(s)
	cw.String(s)
	cw.String(" = ")
	cw.Sync()
}

func run(args []string) error {
	opt := cwflags.NewOpt()
	opt.Title = strings.Join(args, " ")
	cw, err := cwflags.NewPlayer(opt)
	if err != nil {
		return fmt.Errorf("failed to make cw player: %w", err)
	}

	for _, arg := range args {
		cw.String(arg)
		cw.Rune(' ')
		cw.Sync()
	}

	return cw.Close()
}
