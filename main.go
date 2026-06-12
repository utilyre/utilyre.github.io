package main

import (
	"bytes"
	"cmp"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	textemplate "text/template"

	"github.com/goccy/go-yaml"
)

func main() {
	content, err := readContent()
	if err != nil {
		log.Fatalln(err)
	}

	journal := groupJournal(content.Journal)

	/* for _, section := range journal {
		if section.Year == 0 {
			fmt.Println("n.d.")
		} else {
			fmt.Println(section.Year)
		}

		for _, entry := range section.Entries {
			fmt.Println("\t", entry.Title)
		}
	} */

	Skeleton, err := textemplate.ParseFiles("./templates/Skeleton.html")
	if err != nil {
		log.Fatalln(err)
	}

	Journal, err := template.ParseFiles("./templates/Journal.html")
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	err = Journal.Execute(&buf, journal)
	if err != nil {
		log.Fatalln(err)
	}

	err = Skeleton.Execute(os.Stdout, buf.String())
	if err != nil {
		log.Fatalln(err)
	}
}

func groupJournal(journal []JournalEntry) []JournalSection {
	j := slices.SortedFunc(
		slices.Values(journal),
		func(a, b JournalEntry) int { return cmp.Compare(b.Year, a.Year) },
	)

	var sections []JournalSection
	if len(j) > 0 {
		curr := JournalSection{
			Year:    j[0].Year,
			Entries: []JournalEntry{j[0]},
		}

		for _, entry := range j[1:] {
			if entry.Year == curr.Year {
				curr.Entries = append(curr.Entries, entry)
			} else {
				// sort
				slices.SortFunc(
					curr.Entries,
					func(a, b JournalEntry) int { return strings.Compare(a.Title, b.Title) },
				)

				// collect
				sections = append(sections, curr)

				// reset
				curr.Year = entry.Year
				curr.Entries = []JournalEntry{entry}
			}
		}

		sections = append(sections, curr)
	}

	return sections
}

type Content struct {
	Journal []JournalEntry
}

type JournalSection struct {
	Year    uint
	Entries []JournalEntry
}

type JournalEntry struct {
	Title   string
	Year    uint
	Authors []string
	Links   map[string]string
}

func readContent() (Content, error) {
	f, err := os.Open("./content.yml")
	if err != nil {
		return Content{}, fmt.Errorf("os.Open: %v", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return Content{}, fmt.Errorf("io.ReadAll: %v", err)
	}

	var content Content
	err = yaml.Unmarshal(data, &content)
	if err != nil {
		return Content{}, fmt.Errorf("yaml.Unmarshal: %v", err)
	}

	return content, nil
}
