package blog

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Repo interface {
	GetArticle(slug string) (*Article, error)
	GetMetaDatas() ([]*Metadata, error)
}

func NewRepo(dir string) Repo {
	fr := &fileRepo{dir, map[string]*Article{}, nil, nil}

	return fr
}

type fileRepo struct {
	dir                     string
	loadedArticles          map[string]*Article
	metadatas               []*Metadata
	creationSortedMetadatas []*Metadata
}

func (fr *fileRepo) GetArticle(slug string) (*Article, error) {
	if a, ok := fr.loadedArticles[slug]; ok {
		return a, nil
	}

	m, found := find(fr.metadatas, slug)
	if !found {
		return nil, ErrSlugNotFound
	}

	f, err := os.Open(filepath.Join(fr.dir, slug+".html"))
	if err != nil {
		return nil, err
	}
	p, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &Article{m, p}, nil
}

func (fr *fileRepo) GetMetaDatas() ([]*Metadata, error) {
	if fr.metadatas == nil {
		if err := fr.loadMetaData(); err != nil {
			return nil, err
		}
	}

	return fr.creationSortedMetadatas, nil
}

func (fr *fileRepo) loadMetaData() error {
	f, err := os.Open(filepath.Join(fr.dir, "articles.json"))
	if err != nil {
		return err
	}

	p, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var metadatas []*Metadata
	if err := json.Unmarshal(p, &metadatas); err != nil {
		return err
	}

	slices.SortFunc(metadatas, func(a, b *Metadata) int {
		return strings.Compare(a.Slug, b.Slug)
	})
	fr.metadatas = metadatas

	creationSorted := slices.Clone(metadatas)
	slices.SortFunc(creationSorted, func(a, b *Metadata) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		}
		if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})
	fr.creationSortedMetadatas = creationSorted
	return nil
}

var (
	ErrSlugNotFound = errors.New("error slug not found")
)

func find(m []*Metadata, slug string) (*Metadata, bool) {
	idx, found := slices.BinarySearchFunc(m, slug, func(val *Metadata, target string) int {
		return strings.Compare(val.Slug, target)
	})
	if !found {
		return nil, false
	}

	return m[idx], true
}
