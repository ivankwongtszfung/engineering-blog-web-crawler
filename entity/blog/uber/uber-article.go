package uber

import (
	"fmt"
	"strings"

	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog"

	"github.com/gocolly/colly"
)

// UberArticle struct embeds Article and implements IExtract interface
type UberArticle struct {
	blog.Article
}

// ExtractArticle method extracts an article from the given HTMLElement
func ExtractArticle(e *colly.HTMLElement) (*UberArticle, error) {
	title := strings.TrimSpace(e.DOM.Find("h5").Text())
	if title == "" {
		return nil, fmt.Errorf("missing title")
	}

	category := strings.TrimSpace(e.DOM.Find("div > div > div").First().Text())
	if category == "" {
		return nil, fmt.Errorf("missing category")
	}

	date := strings.TrimSpace(e.DOM.Find("p").Text())
	if date == "" {
		return nil, fmt.Errorf("missing date")
	}

	blogURL := e.Request.AbsoluteURL(strings.TrimSpace(e.DOM.Find("a").AttrOr("href", "")))
	if blogURL == "" {
		return nil, fmt.Errorf("missing blog URL")
	}

	imageURL := strings.TrimSpace(e.DOM.Find("img").AttrOr("src", ""))
	if imageURL == "" {
		return nil, fmt.Errorf("missing image URL")
	}

	return &UberArticle{
		Article: blog.Article{
			Title:    title,
			Category: category,
			Date:     date,
			BlogURL:  blogURL,
			ImageURL: imageURL,
		},
	}, nil
}
