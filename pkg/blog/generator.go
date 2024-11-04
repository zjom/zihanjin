package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	img64 "github.com/tenkoh/goldmark-img64"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// Generator converts all markdown files in a given directory into html files
// and saves them to the outDir
type Generator interface {
	Generate(outDir string) error
}

func NewGenerator(dir string) Generator {
	return &generator{
		inDirectory: dir,
		// jobs:        make(chan []byte),
	}
}

type generator struct {
	inDirectory string
	// jobs        chan []byte
}

func (g *generator) Generate(out string) error {
	markdownPaths, err := findAllMarkdownFiles(g.inDirectory)
	if err != nil {
		return err
	}

	articles := make([]*Article, 0)
	for _, path := range markdownPaths {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %s, error: %s\n", path, err)
			continue
		}
		p, err := io.ReadAll(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %s, error: %s\n", path, err)
			continue
		}
		a, err := g.parse(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		articles = append(articles, a)
	}

	return save(articles, out)
}

func save(articles []*Article, dir string) error {
	if !folderExists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	metadatas := getMetadatas(articles)
	metadatas_json, err := json.Marshal(metadatas)
	if err != nil {
		return err
	}
	if err := writeToFile(filepath.Join(dir, "articles.json"), metadatas_json); err != nil {
		return err
	}

	for _, a := range articles {
		if err := writeToFile(filepath.Join(dir, a.Slug+".html"), a.Content); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s, error: %s\n", a.Slug, err)
			continue
		}
	}

	return nil
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func writeToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func (g *generator) parse(source []byte) (*Article, error) {
	gm := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.LeftSingleQuote:  []byte("&sbquo;"),
					extension.RightSingleQuote: nil, // nil disables a substitution
				}),
			),
			meta.Meta,
			img64.Img64,
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(img64.WithParentPath(g.inDirectory)),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := gm.Convert(source, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}
	metadata_interface := meta.Get(context)
	metadata, err := getMetadata(metadata_interface)
	if err != nil {
		return nil, err
	}

	return &Article{
		Metadata: metadata,
		Content:  buf.Bytes(),
	}, nil
}

func getMetadata(metaData map[string]interface{}) (*Metadata, error) {
	title_interface, ok := metaData["title"]
	if !ok {
		return nil, ErrTitleNotFound
	}
	title, ok := title_interface.(string)
	if !ok {
		return nil, ErrInvalidTitle
	}

	created_at_str, ok := metaData["createdAt"].(string)
	var (
		created_at = time.Now()
		err        error
	)
	if ok {
		created_at, err = time.Parse(created_at_str, "2006-02-01")
		if err != nil {
			created_at = time.Now()
		}
	}

	var (
		modified_at = time.Now()
	)
	modified_at_str, ok := metaData["modifiedAt"].(string)
	if ok {
		modified_at, err = time.Parse(modified_at_str, "2006-02-01")
		if err != nil {
			modified_at = time.Now()
		}
	}

	return &Metadata{
		Title:      title,
		Slug:       slug.Make(title),
		CreatedAt:  created_at,
		ModifiedAt: modified_at,
	}, nil
}

var ErrTitleNotFound = errors.New("failed to parse: title not found")
var ErrInvalidTitle = errors.New("failed to parse: title not string")

func findAllMarkdownFiles(root string) ([]string, error) {
	var mdFiles []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			mdFiles = append(mdFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mdFiles, nil
}

func getMetadatas(articles []*Article) []*Metadata {
	res := make([]*Metadata, len(articles))
	for i, a := range articles {
		res[i] = a.Metadata
	}

	return res
}
