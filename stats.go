package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Stats struct {
	StatByLetter map[string]*Stat

	// Working variables - not serialized
	statsFile string
}

type Stat struct {
	Letter        string
	ReactionTimes []float64
	Errors        int
	ConfusedWith  []string

	// Working variables - not serialized
	min, avg, max float64
}

// NewStats loads the stats from the fileName if found otherwise
// returns empty stats
func NewStats(statsFile string) *Stats {
	s := &Stats{
		statsFile:    statsFile,
		StatByLetter: map[string]*Stat{},
	}
	s.Load()
	return s
}

func (s *Stats) Add(tx, rx rune, reactionTime float64) {
	letter := string(tx)
	stat := s.StatByLetter[letter]
	if stat == nil {
		stat = new(Stat)
		s.StatByLetter[letter] = stat
		stat.Letter = letter
	}
	stat.ReactionTimes = append(stat.ReactionTimes, reactionTime)
	if rx != tx {
		stat.Errors++
		stat.ConfusedWith = append(stat.ConfusedWith, string(rx))
	}
}

func min(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func max(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func avg(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	xsCopy := make([]float64, len(xs))
	copy(xsCopy, xs)
	xs = xsCopy
	sort.Float64s(xs)
	// Trim off top and bottom 10%
	trim := len(xs) / 10
	xs = xs[trim : len(xs)-trim]
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func bar(width int, maxValue, min, avg, max float64) string {
	b := make([]rune, width)
	scale := maxValue / float64(width)
	for i := range b {
		b[i] = ' '
	}
	plot := func(v float64, symbol rune) {
		for i := range b {
			if float64(i)*scale < v {
				b[i] = symbol
			} else {
				break
			}
		}
	}
	plot(max, '>')
	plot(avg, '<')
	plot(min, '*')
	return string(b)
}

// Summary digests the stats and shows them
func (s *Stats) Summary() {
	// Calculate min/max/avg
	for _, stat := range s.StatByLetter {
		stat.min, stat.avg, stat.max = min(stat.ReactionTimes), avg(stat.ReactionTimes), max(stat.ReactionTimes)
	}

	// Order by average
	var letters []string
	for letter := range s.StatByLetter {
		letters = append(letters, letter)
	}
	sort.Slice(letters, func(i, j int) bool {
		iv := s.StatByLetter[letters[i]]
		jv := s.StatByLetter[letters[j]]
		return iv.avg < jv.avg
	})

	// Find the max reaction time for scaling
	maxValue := 0.0
	for _, stat := range s.StatByLetter {
		for _, x := range stat.ReactionTimes {
			if x > maxValue {
				maxValue = x
			}
		}
	}

	for _, letter := range letters {
		stat := s.StatByLetter[letter]
		if stat == nil {
			continue
		}
		minT, avgT, maxT := min(stat.ReactionTimes), avg(stat.ReactionTimes), max(stat.ReactionTimes)
		fmt.Printf("%s: min %6.3f avg %6.3f max %6.3f errors %d",
			letter,
			minT,
			avgT,
			maxT,
			stat.Errors,
		)
		b := bar(60, maxValue, minT, avgT, maxT)
		fmt.Printf(" %s", b)
		if len(stat.ConfusedWith) > 0 {
			fmt.Printf(" confused with %q", strings.Join(stat.ConfusedWith, ""))
		}
		fmt.Println()
	}
}

// Store stores the stats to s.statsFile if set
func (s *Stats) Store() {
	if s.statsFile == "" {
		return
	}
	out, err := os.Create(s.statsFile)
	if err != nil {
		log.Fatalf("error opening statsfile: %v", err)
	}
	defer out.Close()
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "\t")
	// Encode into the buffer
	err = encoder.Encode(&s)
	if err != nil {
		log.Fatalf("error writing stats file: %v", err)
	}
}

// Load loads the stats from s.statsFile if set
func (s *Stats) Load() {
	if s.statsFile == "" {
		return
	}
	_, err := os.Stat(s.statsFile)
	if err != nil {
		log.Printf("statsfile %q does not exist -- will create", s.statsFile)
		return
	}
	in, err := os.Open(s.statsFile)
	if err != nil {
		log.Fatalf("error opening statsfile: %v", err)
	}
	defer in.Close()
	decoder := json.NewDecoder(in)
	err = decoder.Decode(s)
	if err != nil {
		log.Fatalf("error reading statsfile: %v", err)
	}
	log.Printf("loaded statsfile %q", s.statsFile)
}
