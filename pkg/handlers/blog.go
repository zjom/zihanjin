package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/components"
)

func (bh *BlogHandler) BlogHome(c echo.Context) error {
	metadatas, err := bh.repo.GetMetaDatas()
	if err != nil {
		return err
	}

	return render(components.Layout(components.BlogPageHome(metadatas)), c)
}

func (bh *BlogHandler) Article(c echo.Context) error {
	slug := c.Param("slug")
	a, err := bh.repo.GetArticle(slug)
	if err != nil {
		return err
	}

	return render(components.BlogPageArticle(a), c)
}
