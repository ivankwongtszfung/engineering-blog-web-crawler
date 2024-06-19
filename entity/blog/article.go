package blog

import (
	"encoding/json"

	"github.com/gocolly/colly"
)

// Article struct holds the article data
type Article struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Date     string `json:"date"`
	BlogURL  string `json:"blogurl"`
	ImageURL string `json:"imageurl"`
}

func (a Article) Compare(a2 Article) bool {
	return a.Id == a2.Id &&
		a.Title == a2.Title &&
		a.Category == a2.Category &&
		a.Date == a2.Date &&
		a.BlogURL == a2.BlogURL &&
		a.ImageURL == a2.ImageURL
}

func (a *Article) JsonDump() (string, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

// IArticle interface defines the ExtractArticle method
type IArticle interface {
	ExtractArticle(e *colly.HTMLElement) (*Article, error)
	JsonDump() (string, error)
}
