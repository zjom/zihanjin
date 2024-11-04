package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/zjom/zihanjin/pkg/blog"
)

func (bh *BlogHandler) RSS(description string) func(c echo.Context) error {
	return func(c echo.Context) error {
		metadatas, err := bh.repo.GetMetaDatas()
		if err != nil {
			return err
		}

		itemsXML := ""
		baseUrl := c.Request().URL.Host
		c.Logger().Debug(baseUrl)
		for _, data := range metadatas {
			itemsXML += metadataToRssItem(data, baseUrl)
		}

		rssFeed := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" ?>
  <rss version="2.0">
    <channel>
        <title>My Portfolio</title>
        <link>%s</link>
        <description>%s</description>
        %s
    </channel>
  </rss>`, baseUrl, description, itemsXML)
		return c.Blob(200, "text/xml", []byte(rssFeed))
	}
}

func metadataToRssItem(data *blog.Metadata, baseUrl string) string {
	return fmt.Sprintf(`<item>
          <title>%s</title>
          <link>%s/blog/%s</link>
          <description>%s</description>
          <pubDate>%s</pubDate>
        </item>`,
		data.Title,
		baseUrl,
		data.Slug,
		data.Description,
		data.CreatedAt.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
}
