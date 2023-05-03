package rss

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/ncw/ncwtester/cmd"
	"github.com/ncw/ncwtester/cwgenerator"
	"github.com/ncw/ncwtester/cwplayer"
	"github.com/spf13/cobra"
)

var (
	sampleRate int
	wpm        float64
	frequency  float64
	url        string
	user       string
	pass       string
)

const (
	channelNum      = 2
	bitDepthInBytes = 2
	maxSampleValue  = 32767
)

// subCmd represents the rss ommand
var subCmd = &cobra.Command{
	Use:   "rss",
	Short: "Fetch RSS and turn into morse code",
	Long: `Fetch RSS and turn it into morse code

This fetches an RSS feed, eg http://feeds.bbci.co.uk/news/uk/rss.xml
Parses it and sends the items as morse code.

Most RSS, Atom and JSON feed types are supported.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Root.AddCommand(subCmd)
	flags := subCmd.Flags()
	flags.IntVarP(&sampleRate, "samplerate", "s", 44100, "sample rate")
	flags.Float64VarP(&wpm, "wpm", "", 25.0, "WPM to send at")
	flags.Float64VarP(&frequency, "frequency", "", 600.0, "HZ of morse")
	flags.StringVarP(&url, "url", "", "", "URL to fetch RSS from")
	flags.StringVarP(&url, "user", "", "", "Username for URL (optional)")
	flags.StringVarP(&url, "pass", "", "", "Password for URL (optional)")
}

// Returns a reader to read the RSS from - must be closed afterwards
func fetch() (io.ReadCloser, error) {
	log.Printf("Fetching RSS at %q", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request %q: %w", url, err)
	}
	if user != "" || pass != "" {
		req.SetBasicAuth(user, pass)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch RSS %q: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("bad status %d when RSS %q: %s", resp.StatusCode, url, resp.Status)
	}
	return resp.Body, nil
}

var (
	removeInvalid = regexp.MustCompile(`[^A-Za-z0-9,.: /]+`)
)

// Simplify and play the string
func play(cw *cwplayer.Player, s string) {
	// Remove all unknown characters
	s = removeInvalid.ReplaceAllString(s, "")

	// Replace `:` with more CW friendly ` =`
	s = strings.ReplaceAll(s, ":", " =")

	fmt.Println(s)
	cw.String(s)
	cw.String(" = ")
	cw.SyncPlay()
	cw.Reset()
}

func run() error {
	if url == "" {
		return fmt.Errorf("need --url parameter to fetch from")
	}
	opt := cwgenerator.Options{
		WPM:             wpm,
		Frequency:       frequency,
		SampleRate:      sampleRate,
		ChannelNum:      channelNum,
		BitDepthInBytes: bitDepthInBytes,
		MaxSampleValue:  maxSampleValue,
	}
	cw, err := cwplayer.New(opt)
	if err != nil {
		return fmt.Errorf("failed to make cw generator: %w", err)
	}

	cw.Reset()

	fp := gofeed.NewParser()
	log.Printf("Fetching RSS from %q", url)
	feed, err := fp.ParseURL(url)
	if err != nil {
		return fmt.Errorf("rss fetch and parse failed: %w", err)
	}

	fmt.Printf("Title: %s\n", feed.Title)
	fmt.Printf("Description: %s\n", feed.Description)

	for i, item := range feed.Items {
		play(cw, fmt.Sprintf("NR %d", i+1))
		play(cw, item.Title)
		play(cw, item.Description)
	}

	return nil
}
