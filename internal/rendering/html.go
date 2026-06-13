package rendering

import (
	"bytes"
	"cmp"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"website/internal/domain"
)

type contentView struct {
	Style   template.CSS
	Journal []journalSectionView
}

type journalSectionView struct {
	Year uint
	Logs []journalLogView
}

type journalLogView struct {
	Title   string
	Authors []string
	Links   []domain.LabeledLink
}

func RenderContent(w io.Writer, content domain.Content) error {
	view := transform(content)
	log.Println("transformed content into view data")

	/* for _, section := range view.Journal {
		fmt.Println(section.Year)
		for _, entry := range section.Entries {
			fmt.Printf("\t%s\n", entry.Title)
		}
	} */

	tmpl := template.Must(template.ParseGlob("./templates/*.html"))
	log.Println("parsed templates")

	// first pass: render to a temporary file without style
	tmp, err := os.CreateTemp("", "*.html")
	if err != nil {
		return err
	}
	log.Println("created temporary html file for the first pass")
	err = tmpl.ExecuteTemplate(tmp, "Skeleton", &view)
	if err != nil {
		return err
	}
	log.Println("executed the first pass")

	// generate css styles using tailwind based on the first pass
	var style bytes.Buffer
	twCMD := exec.Command("tailwindcss", "-m", "-i", "-", "-o", "-")
	twCMD.Stdin = strings.NewReader(fmt.Sprintf(
		`@import "tailwindcss" source(none); @source "%s";`,
		tmp.Name(),
	))
	twCMD.Stdout = &style
	err = twCMD.Run()
	if err != nil {
		return err
	}
	log.Println("ran tailwind css")

	// clean up the temporary file
	tmpName := tmp.Name()
	err = tmp.Close()
	if err != nil {
		return err
	}
	err = os.Remove(tmpName)
	if err != nil {
		return err
	}
	log.Println("cleaned up temporary html file")

	// second pass: render to the actual writer given as parameter
	var html bytes.Buffer
	view.Style = template.CSS(style.String())
	err = tmpl.ExecuteTemplate(&html, "Skeleton", &view)
	if err != nil {
		return err
	}
	log.Println("executed the second pass")

	minCMD := exec.Command("minhtml", "--minify-css")
	minCMD.Stdin = &html
	minCMD.Stdout = w
	err = minCMD.Run()
	if err != nil {
		return err
	}
	log.Println("minified the result")

	return nil
}

func transform(content domain.Content) contentView {
	var view contentView

	sortedJournal := slices.SortedFunc(
		slices.Values(content.Journal),
		func(a, b domain.JournalLog) int { return cmp.Compare(b.Year, a.Year) },
	)

	// domain to view
	d2v := func(log domain.JournalLog) journalLogView {
		return journalLogView{
			Title:   log.Title,
			Authors: log.Authors,
			Links:   log.Links,
		}
	}

	if len(sortedJournal) > 0 {
		curr := journalSectionView{
			Year: sortedJournal[0].Year,
			Logs: []journalLogView{d2v(sortedJournal[0])},
		}

		for _, entry := range sortedJournal[1:] {
			if entry.Year == curr.Year {
				curr.Logs = append(curr.Logs, d2v(entry))
				continue
			}

			// sort
			slices.SortFunc(curr.Logs, func(a, b journalLogView) int {
				return cmp.Compare(a.Title, b.Title)
			})

			// collect
			view.Journal = append(view.Journal, curr)

			// reset
			curr.Year = entry.Year
			curr.Logs = []journalLogView{d2v(entry)}
		}

		view.Journal = append(view.Journal, curr)
	}

	return view
}
