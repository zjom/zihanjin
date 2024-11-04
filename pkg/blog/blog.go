package blog

import "time"

type Article struct {
	*Metadata

	Content []byte
}

type Metadata struct {
	Title       string
	Slug        string
	Description string
	CreatedAt   time.Time
	ModifiedAt  time.Time
}
