package repo

import (
	"context"

	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog"
)

type IArticleRepository interface {
	All() ([]*blog.Article, error)
	Create() error
	Get(string) (*blog.Article, error)
	SaveAll(context.Context, []blog.Article) error
	SaveOne(blog.Article) error
}
