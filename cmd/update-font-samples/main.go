// update-font-samples is used by the maintainer of the repos to update the
// README.md and images directories of all the go-fonts* repos.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	fontDirs = "abcdefghijklmnopqrstuvwyz" // no 'x' fonts currently
	splitStr = "## Font samples\n"
)

var (
	samplesTemplate = template.Must(template.New("samplesTemplateStr").Funcs(funcMap).Parse(samplesTemplateStr))
)

func main() {
	masterSamples, err := filepath.Glob("images/sample_*.png")
	must(err)

	if len(masterSamples) == 0 {
		log.Fatal("No images/sample_*.png files found. Aborting.")
	}

	allFonts := make([]string, 0, len(masterSamples))
	for _, sample := range masterSamples {
		font := strings.TrimSuffix(strings.TrimPrefix(sample, "images/sample_"), ".png")
		allFonts = append(allFonts, font)
	}

	var allSamples bytes.Buffer
	must(samplesTemplate.Execute(&allSamples, allFonts))

	c := &client{
		allSamples:    allSamples.String(),
		masterSamples: map[string]bool{},
	}
	for _, name := range masterSamples {
		c.masterSamples[name] = true
	}

	c.updateReadme("README.md")

	for _, r := range fontDirs {
		c.updateFontDir("../go-fonts-" + string(r))
	}
}

type client struct {
	allSamples    string
	masterSamples map[string]bool
}

func (c *client) updateFontDir(dirPrefix string) {
	dirName := filepath.Join(dirPrefix, "images/sample_*.png")

	c.updateReadme(filepath.Join(dirPrefix, "README.md"))

	fontSamples, err := filepath.Glob(dirName)
	must(err)

	seen := map[string]bool{}
	for _, fullname := range fontSamples {
		filename := strings.TrimPrefix(fullname, dirPrefix+"/")
		seen[filename] = true
		if _, ok := c.masterSamples[filename]; !ok {
			log.Fatalf("Found %v which is missing in go-fonts directory! Please fix this.", fullname)
		}
	}

	for filename := range c.masterSamples {
		if _, ok := seen[filename]; !ok {
			log.Printf("Copying %v to %v/images ...", filename, dirPrefix)
		}
		b, err := os.ReadFile(filename)
		must(err)
		outFile := filepath.Join(dirPrefix, filename)
		must(os.WriteFile(outFile, b, 0644))
	}
}

func (c *client) updateReadme(filename string) {
	b, err := os.ReadFile(filename)
	must(err)
	parts := strings.Split(string(b), splitStr)
	if len(parts) != 2 {
		log.Fatalf("Error parsing %v, got %v parts", filename, len(parts))
	}
	outStr := fmt.Sprintf("%v%v%v", parts[0], splitStr, c.allSamples)
	must(os.WriteFile(filename, []byte(outStr), 0644))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var funcMap = map[string]any{
	"firstLetter": firstLetter,
}

func firstLetter(s string) string {
	return s[0:1]
}

var samplesTemplateStr = `{{ range . }}[![{{ . }}](images/sample_{{ . }}.png)](https://github.com/gmlewis/go-fonts-{{ . | firstLetter }}/tree/master/fonts/{{ . }})
{{ end }}`
