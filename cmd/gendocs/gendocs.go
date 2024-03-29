// Package gendocs provides the gendocs command.
package gendocs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ncw/cwtool/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	cmd.Root.AddCommand(subCmd)
}

var (
	headingRe = regexp.MustCompile(`(?m)^## (.*?)$`)
	autoGenRe = regexp.MustCompile(`(?m)^#.*?Auto generated.*?$`)
	linkRe    = regexp.MustCompile(`\[.*?\]\(.*?\.md\)`)
)

// munge the output files into a single file
func mungeFile(out *strings.Builder, fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file %q: %w", fileName, err)
	}
	help := string(data)
	uname := filepath.Base(fileName)
	uname = uname[:len(uname)-len(filepath.Ext(uname))]
	// Add ID into main heading - DOESN't work on GitHub
	// use headings with `-` instead of ` ` instead
	// help = headingRe.ReplaceAllString(help, "## $1 {#"+uname+"}\n")
	// Remove auto generated lines
	help = autoGenRe.ReplaceAllString(help, "")
	// Munge links to be internal
	help = linkRe.ReplaceAllStringFunc(help, func(s string) string {
		i := strings.IndexRune(s, '(')
		text := s[:i]
		link := s[i+1 : len(s)-1]
		link = strings.TrimSuffix(link, ".md")
		link = strings.ReplaceAll(link, "_", "-")
		return text + "(#" + link + ")"
	})
	// Send to output
	_, err = out.WriteString(help)
	return err
}

// munge the output files into a single file
func mungeFiles(out *strings.Builder, outputDir string) error {
	// List output and prepare to munge
	fis, err := ioutil.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("read output directory %q: %w", outputDir, err)
	}
	var files []string
	for _, fi := range fis {
		files = append(files, fi.Name())
	}
	// Sort by underscores then alphabetically
	sort.Slice(files, func(i, j int) bool {
		iUnders := strings.Count(files[i], "_")
		jUnders := strings.Count(files[j], "_")
		if iUnders != jUnders {
			return iUnders < jUnders
		}
		return files[i] < files[j]
	})
	// Now munge each individual file
	for _, fileName := range files {
		fileName = filepath.Join(outputDir, fileName)
		err = mungeFile(out, fileName)
		if err != nil {
			return fmt.Errorf("munge file %q: %w", fileName, err)
		}
	}
	return nil
}

const (
	readmeFile    = "README.md"
	readmeTemp    = readmeFile + ".new"
	autoGenerated = "<!-- auto generated from here on -->"
)

// Replace the help string in README.md
func mungeREADME(help string) (err error) {
	data, err := ioutil.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("read %q: %w", readmeFile, err)
	}
	readme := string(data)

	// Find auto generated tag and replace everything after
	i := strings.Index(readme, autoGenerated)
	if i < 0 {
		return fmt.Errorf("couldn't find %q in %q", autoGenerated, readmeFile)
	}
	readme = readme[:i+len(autoGenerated)] + "\n\n" + help

	// write the README out
	err = ioutil.WriteFile(readmeTemp, []byte(readme), 0666)
	if err != nil {
		return fmt.Errorf("write %q: %w", readmeTemp, err)
	}

	// do atomic rename
	err = os.Rename(readmeTemp, readmeFile)
	if err != nil {
		return fmt.Errorf("rename %q -> %q: %w", readmeTemp, readmeFile, err)
	}

	return nil
}

var subCmd = &cobra.Command{
	Use:   "gendocs",
	Short: `Output markdown docs for cwtool.`,
	Long: `
This produces markdown docs for the cwtool and inserts them into README.md.`,
	Hidden: true,
	RunE: func(command *cobra.Command, args []string) error {
		outputDir, err := os.MkdirTemp("", "cwtool.gendocs.")
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
		defer os.RemoveAll(outputDir)
		err = doc.GenMarkdownTree(cmd.Root, outputDir)
		if err != nil {
			return fmt.Errorf("generate markdown tree: %w", err)
		}
		var out strings.Builder
		err = mungeFiles(&out, outputDir)
		if err != nil {
			return fmt.Errorf("munge files: %w", err)
		}
		err = mungeREADME(out.String())
		if err != nil {
			return fmt.Errorf("munge README: %w", err)
		}
		return nil
	},
}
