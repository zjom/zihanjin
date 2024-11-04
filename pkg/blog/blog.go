package blog

import "time"

type Article struct {
	*Metadata

	Content []byte
}

type Metadata struct {
	Title      string
	Slug       string
	CreatedAt  time.Time
	ModifiedAt time.Time
}
