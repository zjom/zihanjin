package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/components"
)

func customerHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	he, ok := err.(*echo.HTTPError)
	if !ok {
		c.Logger().Error(err)
		if c.JSON(500, err); err != nil {
			c.Logger().Error(err)
		}
		return
	}
	switch he.Code {
	case http.StatusNotFound:
		if err := render(components.Layout(components.NotFound()), c); err != nil {
			c.Logger().Error(err)
			return
		}
		break
	default:
		if err := c.JSON(he.Code, he.Message); err != nil {
			c.Logger().Error(err)
			return
		}
	}
}
