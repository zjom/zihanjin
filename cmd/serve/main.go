package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
	"github.com/zjom/zihanjin/pkg/components"
)

var (
	inDir  = flag.String("in", "md/in", "the directory containing the markdown files")
	outDir = flag.String("out", "md/out", "the directory to dump the html files")
)

func main() {
	flag.Parse()

	g := blog.NewGenerator(*inDir)
	if err := g.Generate(*outDir); err != nil {
		log.Fatalln(err)
	}

	r := blog.NewRepo("md/out")
	metadatas, err := r.GetMetaDatas()
	if err != nil {
		panic(err)
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
