package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
	"github.com/zjom/zihanjin/pkg/components"
)

var metadatas = []*blog.Metadata{
	{
		Title:       "creating this blog",
		Slug:        "creating-this-blog",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
	{
		Title:       "why neovim",
		Slug:        "why-neovim",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
	{
		Title:       "functional programming",
		Slug:        "functional-programming",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
	{
		Title:       "how to leetcode",
		Slug:        "how-to-leetcode",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
	{
		Title:       "thoughts on the future of software",
		Slug:        "thoughts-on-the-future-of-software",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
	{
		Title:       "the changing world order",
		Slug:        "the-changing-world-order",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
	{
		Title:       "why is everything lower case?",
		Slug:        "why-is-everything-lower-case?",
		PublishedAt: time.Now(),
		EditedAt:    time.Now(),
	},
}

func main() {
	e := echo.New()
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return components.Layout(components.Landing(metadatas)).
			Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/blog", func(c echo.Context) error {
		return components.Layout(components.BlogPage(metadatas)).
			Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
