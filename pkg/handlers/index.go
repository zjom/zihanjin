package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/components"
)

func (bh *BlogHandler) Index(c echo.Context) error {
	metadatas, err := bh.repo.GetMetaDatas()
	if err != nil {
		return err
	}

	return render(components.Layout(components.Landing(metadatas)), c)
}
