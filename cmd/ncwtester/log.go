package ncwtester

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Times written in localtime to the CSV file
const timeFormat = "2006-01-02 15:04:05"

type CSVLog struct {
	logFile string
}

func NewCSVLog(logFile string) *CSVLog {
	l := &CSVLog{
		logFile: logFile,
	}
	l.Create()
	return l
}

// Write a csv line to the file - this should end in \n
func (l *CSVLog) Write(row string) {
	out, err := os.OpenFile(l.logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			log.Fatalf("error closing log file: %v", err)
		}
	}()

	_, err = out.WriteString(row)
	if err != nil {
		log.Fatalf("error writing to log file: %v", err)
	}
	if err != nil {
		log.Fatalf("error writing to log file: %v", err)
	}
}

// Create the file if it doesn't exist
func (l *CSVLog) Create() {
	if l.logFile == "" {
		return
	}
	fi, err := os.Stat(l.logFile)
	if err == nil && fi.Size() > 0 {
		return
	}
	log.Printf("log file %q does not exist -- will create", l.logFile)
	l.Write(`"when","tx","rx","correct","time"` + "\n")
	log.Printf("Created log file %q", l.logFile)
}

// quote a rune for CSV
func quoteRune(r rune) string {
	if r != '"' {
		return string(r)
	}
	return `""`
}

// Add a CSV log line
func (l *CSVLog) Add(tx, rx rune, reactionTime time.Duration) {
	if l.logFile == "" {
		return
	}
	l.Write(fmt.Sprintf(`%s,"%s","%s",%s,%s`+"\n",
		time.Now().Format(timeFormat),
		quoteRune(tx),
		quoteRune(rx),
		fmt.Sprint(rx == tx),
		fmt.Sprintf("%.3f", reactionTime.Seconds()),
	))
}
