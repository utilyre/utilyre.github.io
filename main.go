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

}

type RenderData struct {
	Journal []JournalSection
}

type JournalSection struct {
	Year    uint
	Entries []JournalEntry
}

type JournalEntry struct{
	Title string
	Authors []string
	Links map[string]string
}

func cookContent(content Content) RenderData {
	j := slices.SortedFunc(
		slices.Values(content.Journal),
		func(a, b Material) int { return cmp.Compare(b.Year, a.Year) },
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
	Journal []Material
}

type Material struct {
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
