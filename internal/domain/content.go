package domain

type Content struct {
	Journal []JournalLog
}

type JournalLog struct {
	Title   string
	Year    uint
	Authors []string
	Links   []LabeledLink
}

type LabeledLink struct {
	Label string
	URL   string
}
