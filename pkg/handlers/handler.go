package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
	"github.com/zjom/zihanjin/pkg/components"
)

func render(component templ.Component, c echo.Context) error {
	return component.Render(c.Request().Context(), c.Response().Writer)
}

type BlogHandler struct {
	repo blog.Repo
}

func NewBlogHandler(repo blog.Repo) *BlogHandler {
	return &BlogHandler{repo}
}

func (bh *BlogHandler) Index(c echo.Context) error {
	metadatas, err := bh.repo.GetMetaDatas()
	if err != nil {
		return err
	}

	return render(components.Layout(components.Landing(metadatas)), c)
}

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

func (bh *BlogHandler) Register(e *echo.Echo) {
	e.GET("", bh.Index)
	blog := e.Group("blog")
	blog.GET("", bh.BlogHome)
	blog.GET("/:slug", bh.Article)
}
