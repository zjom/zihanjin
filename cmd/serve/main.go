package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
	"github.com/zjom/zihanjin/pkg/handlers"
)

var (
	inDir    = flag.String("in", "md/in", "the directory containing the markdown files")
	outDir   = flag.String("out", "md/out", "the directory to dump the html files")
	generate = flag.Bool("generate", true, "whether or not to generate the static files")
)

func main() {
	flag.Parse()

	if *generate {
		g := blog.NewGenerator(*inDir)
		if err := g.Generate(*outDir); err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate blog, error: %s\n", err)
			os.Exit(1)
			return
		}
	}

	r := blog.NewRepo(*outDir)
	h := handlers.NewBlogHandler(r)

	e := echo.New()
	e.Static("/static", "static")
	h.Register(e)

	e.Logger.Fatal(e.Start(":1323"))
}
