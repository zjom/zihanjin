package blog

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Repo interface {
	GetArticle(slug string) (*Article, error)
	GetMetaDatas() ([]*Metadata, error)
}

func NewRepo(dir string) Repo {
	fr := &fileRepo{dir, map[string]*Article{}, nil}

	return fr
}

type fileRepo struct {
	dir            string
	loadedArticles map[string]*Article
	metadatas      map[string]*Metadata
}

func (fr *fileRepo) GetArticle(slug string) (*Article, error) {
	if a, ok := fr.loadedArticles[slug]; ok {
		return a, nil
	}

	m, ok := fr.metadatas[slug]
	if !ok {
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

	return values(fr.metadatas), nil
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

	m := map[string]*Metadata{}
	for _, md := range metadatas {
		m[md.Slug] = md
	}

	fr.metadatas = m
	return nil
}

var (
	ErrSlugNotFound = errors.New("error slug not found")
)

func values(m map[string]*Metadata) []*Metadata {
	retv := make([]*Metadata, 0)
	for _, v := range m {
		retv = append(retv, v)
	}

	return retv
}
