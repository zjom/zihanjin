// compile/main.go recursively traverses a directory of markdown
// files, parsing and converting them to html, dumping them in the specified directory
package main

import (
	"flag"
	"log"

	"github.com/zjom/zihanjin/pkg/blog"
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
}
