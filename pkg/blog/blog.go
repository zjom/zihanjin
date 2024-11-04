package blog

import "time"

type Blog struct {
	Metadata

	Content string
}

type Metadata struct {
	Title       string
	Slug        string
	PublishedAt time.Time
	EditedAt    time.Time
}
