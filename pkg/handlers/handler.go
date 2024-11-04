package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
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
