package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
	"github.com/zjom/zihanjin/pkg/components"
)

var (
	inDir    = flag.String("in", "md/in", "the directory containing the markdown files")
	outDir   = flag.String("out", "md/out", "the directory to dump the html files")
	generate = flag.Bool("generate", true, "whether or not to generate the static files")
)

func main() {
	flag.Parse()

	g := blog.NewGenerator(*inDir)
	if err := g.Generate(*outDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate blog, error: %s\n", err)
		os.Exit(1)
		return
	}

	r := blog.NewRepo("md/out")
	metadatas, err := r.GetMetaDatas()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get metadatas, error: %s\n", err)
		os.Exit(1)
		return
	}

	e := echo.New()
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return components.Layout(components.Landing(metadatas)).
			Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/blog", func(c echo.Context) error {
		return components.Layout(components.BlogPageHome(metadatas)).
			Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/blog/:slug", func(c echo.Context) error {
		slug := c.Param("slug")
		a, err := r.GetArticle(slug)
		if err != nil {
			return err
		}
		return components.BlogPageArticle(a).Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
