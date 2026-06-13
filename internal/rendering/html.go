package rendering

import (
	"bytes"
	"cmp"
	"io"
	"os/exec"
	"slices"
	"text/template"
	textemplate "text/template"
	"website/internal/domain"
)

type contentView struct {
	Journal []journalSectionView
}

type journalSectionView struct {
	Year    uint
	Entries []journalEntryView
}

type journalEntryView struct {
	Title   string
	Authors []string
	Links   []domain.LabeledLink
}

func RenderContent(w io.Writer, content domain.Content) error {
	view := transform(content)

	/* for _, section := range view.Journal {
		fmt.Println(section.Year)
		for _, entry := range section.Entries {
			fmt.Printf("\t%s\n", entry.Title)
		}
	} */

	Skeleton, err := textemplate.ParseFiles("./templates/Skeleton.html")
	if err != nil {
		return err
	}

	Journal, err := template.ParseFiles("./templates/Journal.html")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = Journal.Execute(&buf, view.Journal)
	if err != nil {
		return err
	}

	err = Skeleton.Execute(w, buf.String())
	if err != nil {
		return err
	}

	cmd := exec.Command("tailwindcss", "-i", "./tw.css", "-o", "-")
	cmd.Stdout = something
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func transform(content domain.Content) contentView {
	var view contentView

	sortedJournal := slices.SortedFunc(
		slices.Values(content.Journal),
		func(a, b domain.JournalLog) int { return cmp.Compare(b.Year, a.Year) },
	)

	// domain to view
	d2v := func(log domain.JournalLog) journalEntryView {
		return journalEntryView{
			Title:   log.Title,
			Authors: log.Authors,
			Links:   log.Links,
		}
	}

	if len(sortedJournal) > 0 {
		curr := journalSectionView{
			Year:    sortedJournal[0].Year,
			Entries: []journalEntryView{d2v(sortedJournal[0])},
		}

		for _, entry := range sortedJournal[1:] {
			if entry.Year == curr.Year {
				curr.Entries = append(curr.Entries, d2v(entry))
				continue
			}

			// sort
			slices.SortFunc(curr.Entries, func(a, b journalEntryView) int {
				return cmp.Compare(a.Title, b.Title)
			})

			// collect
			view.Journal = append(view.Journal, curr)

			// reset
			curr.Year = entry.Year
			curr.Entries = []journalEntryView{d2v(entry)}
		}

		view.Journal = append(view.Journal, curr)
	}

	return view
}
