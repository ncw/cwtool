package rss

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/ncw/cwtool/cmd"
	"github.com/ncw/cwtool/cmd/cwflags"
	"github.com/ncw/cwtool/cw"
	"github.com/spf13/cobra"
)

var (
	url         string
	description bool
)

// subCmd represents the rss ommand
var subCmd = &cobra.Command{
	Use:   "rss",
	Short: "Fetch RSS and turn into morse code",
	Long: `Fetch RSS and turn it into morse code

This fetches an RSS feed, eg http://feeds.bbci.co.uk/news/uk/rss.xml
Parses it and sends the items as morse code.

Most RSS, Atom and JSON feed types are supported.

This plays the title of the RSS feed and then the titles of each item
in the feed.

Each item in the feeed is prefixed with NR x where x counts 1,2,3...

Use --description to add the descriptions of each link in as well as
their titles.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Root.AddCommand(subCmd)
	flags := subCmd.Flags()
	cwflags.Add(flags)
	flags.StringVarP(&url, "url", "", "", "URL to fetch RSS from")
	flags.BoolVarP(&description, "description", "", false, "If set add the description too")
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

func run() error {
	if url == "" {
		return fmt.Errorf("need --url parameter to fetch from")
	}
	opt := cwflags.NewOpt()
	cw, err := cwflags.NewPlayer(opt)
	if err != nil {
		return fmt.Errorf("failed to make cw player: %w", err)
	}

	fp := gofeed.NewParser()
	log.Printf("Fetching RSS from %q", url)
	feed, err := fp.ParseURL(url)
	if err != nil {
		return fmt.Errorf("rss fetch and parse failed: %w", err)
	}

	play(cw, feed.Title)
	if description {
		play(cw, feed.Description)
	}

	for i, item := range feed.Items {
		play(cw, fmt.Sprintf("NR %d", i+1))
		play(cw, item.Title)
		if description {
			play(cw, item.Description)
		}
	}

	return cw.Close()
}
