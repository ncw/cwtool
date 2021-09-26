package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Hold stats about a letter or group of letters
type Stat struct {
	Name          string
	ReactionTimes []float64
	Errors        int
	Total         int
	ConfusedWith  []string

	sorted bool
	sum    float64 // sum of ReactionTimes
}

// NewStat instantiates a new statistic
func NewStat(Name string) *Stat {
	s := &Stat{
		Name: Name,
	}
	return s
}

func (s *Stat) Add(tx, rx string, reactionTime float64) {
	s.sorted = false
	if rx == tx {
		s.ReactionTimes = append(s.ReactionTimes, reactionTime)
		s.sum += reactionTime
	} else {
		s.Errors++
		s.ConfusedWith = append(s.ConfusedWith, string(rx))
	}
	s.Total++
}

func (s *Stat) sort() {
	if !s.sorted {
		sort.Float64s(s.ReactionTimes)
		s.sorted = true
	}
}

func (s *Stat) Min() float64 {
	if len(s.ReactionTimes) == 0 {
		return 0
	}
	s.sort()
	return s.ReactionTimes[0]
}

func (s *Stat) Max() float64 {
	if len(s.ReactionTimes) == 0 {
		return 0
	}
	s.sort()
	return s.ReactionTimes[len(s.ReactionTimes)-1]
}

func (s *Stat) Avg() float64 {
	if len(s.ReactionTimes) == 0 {
		return 0
	}
	return s.sum / float64(len(s.ReactionTimes))
}

func (s *Stat) Percentile(p float64) float64 {
	if len(s.ReactionTimes) == 0 {
		return 0
	}
	s.sort()
	i := int(math.Round((p / 100) * float64(len(s.ReactionTimes))))
	if i < 0 {
		i = 0
	} else if i >= len(s.ReactionTimes) {
		i = len(s.ReactionTimes) - 1
	}
	return s.ReactionTimes[i]
}

func (s *Stat) PercentageErrors() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Errors) / float64(s.Total) * 100.0
}

type Stats struct {
	StatByLetter map[string]*Stat
	Total        *Stat

	// Working variables - not serialized
	csvFile    string
	timeCutoff time.Duration
	rows       int
}

// NewStats loads the stats from the fileName if found otherwise
// returns empty stats
func NewStats(csvFile string, timeCutoff time.Duration) *Stats {
	s := &Stats{
		csvFile:      csvFile,
		timeCutoff:   timeCutoff,
		StatByLetter: map[string]*Stat{},
		Total:        NewStat("total"),
	}
	s.Load()
	return s
}

func (s *Stats) Add(tx, rx string, reactionTime float64) {
	letter := string(tx)
	stat := s.StatByLetter[letter]
	if stat == nil {
		stat = NewStat(letter)
		s.StatByLetter[letter] = stat
	}
	stat.Add(tx, rx, reactionTime)
	s.Total.Add(tx, rx, reactionTime)
	s.rows++
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
	// Order by average
	var letters []string
	for letter := range s.StatByLetter {
		letters = append(letters, letter)
	}
	sort.Slice(letters, func(i, j int) bool {
		iv := s.StatByLetter[letters[i]]
		jv := s.StatByLetter[letters[j]]
		//return iv.Avg() < jv.Avg()
		return iv.Percentile(75) < jv.Percentile(75)
	})

	// Find the max reaction time for scaling
	maxValue := s.Total.Max()

	for _, letter := range letters {
		stat := s.StatByLetter[letter]
		if stat == nil {
			continue
		}
		minT, avgT, maxT := stat.Min(), stat.Avg(), stat.Max()
		p50, p75, p90 := stat.Percentile(50), stat.Percentile(75), stat.Percentile(90)
		fmt.Printf("%s: min %6.3f p50 %6.3f p75 %6.3f p90 %6.3f errors %5.1f%%",
			letter,
			minT,
			p50,
			p75,
			p90,
			stat.PercentageErrors(),
		)
		b := bar(60, maxValue, minT, avgT, maxT)
		fmt.Printf(" %s", b)
		if len(stat.ConfusedWith) > 0 {
			fmt.Printf(" confused with %q", strings.Join(stat.ConfusedWith, ""))
		}
		fmt.Println()
	}

	fmt.Printf("----\n")
	fmt.Printf("%s: min %6.3f p50 %6.3f p75 %6.3f p90 %6.3f errors %5.1f%%\n",
		s.Total.Name,
		s.Total.Min(),
		s.Total.Percentile(50),
		s.Total.Percentile(75),
		s.Total.Percentile(90),
		s.Total.PercentageErrors(),
	)
}

// Load loads the stats from s.csvFile if set
func (s *Stats) Load() {
	in, err := os.Open(s.csvFile)
	if err != nil {
		log.Fatalf("error opening statsfile: %v", err)
	}
	defer in.Close()

	r := csv.NewReader(in)
	first := true
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read csv log: %v", err)
		}
		if first {
			first = false
			continue
		}
		if len(row) != 5 {
			log.Fatalf("Ignoring bad line in csv log: %q", row)
			continue
		}
		when, err := time.Parse(timeFormat, row[0])
		if err != nil {
			log.Fatalf("Failed to parse time %q from csv log: %v", row[0], err)
		}
		// Don't load rows older than cutoff
		if s.timeCutoff > 0 {
			ago := time.Since(when)
			if ago > s.timeCutoff {
				continue
			}
		}
		tx := row[1]
		rx := row[2]
		reactionTime, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			log.Fatalf("Failed to parse duration %q from csv log: %v", row[4], err)
		}
		s.Add(tx, rx, reactionTime)
	}

	log.Printf("loaded %d rows from statsfile %q", s.rows, s.csvFile)
}
