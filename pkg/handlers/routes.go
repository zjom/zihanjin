package handlers

import "github.com/labstack/echo/v4"

func (bh *BlogHandler) Register(e *echo.Echo) {
	e.GET("", bh.Index)
	e.GET("/rss", bh.RSS("a blog about random shit idk"))
	blog := e.Group("blog")
	blog.GET("", bh.BlogHome)
	blog.GET("/:slug", bh.Article)
}
