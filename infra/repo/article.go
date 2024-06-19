package repo

import (
	"context"
	"web_crawler/entity/blog"
)

type IArticleRepository interface {
	All() ([]*blog.Article, error)
	Create() error
	Get(string) (*blog.Article, error)
	SaveAll(context.Context, []blog.Article) error
	SaveOne(blog.Article) error
}
