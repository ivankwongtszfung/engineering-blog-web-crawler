package sqlite

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"web_crawler/entity/blog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func TestSqlite(t *testing.T) {
	defer os.Remove("./test.db")
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	repo := ArticleRepository{DB: db}

	t.Run("create the db", func(t *testing.T) {
		if err := repo.Create(); err != nil {
			t.Errorf("Cannot create DB table - article: %v", err)
		}
	})

	blogPost := blog.Article{
		Id:       "asdasdfasdfwer214",
		Title:    "Introduction to Web Development",
		Category: "Web Development",
		Date:     "2023-03-10",
		BlogURL:  "https://example.com/web-development-intro",
		ImageURL: "https://example.com/images/web-development-intro.jpg",
	}

	blogPosts := []blog.Article{
		{
			Id:       "asdfasfasfasfd",
			Title:    "Understanding Golang Structs",
			Category: "Programming",
			Date:     "2023-01-15",
			BlogURL:  "https://example.com/golang-structs",
			ImageURL: "https://example.com/images/golang-structs.jpg",
		},
		{
			Id:       "adfasdfasfasdf",
			Title:    "Advanced Techniques in Golang",
			Category: "Programming",
			Date:     "2023-02-20",
			BlogURL:  "https://example.com/advanced-golang",
			ImageURL: "https://example.com/images/advanced-golang.jpg",
		},
	}

	t.Run("save one article", func(t *testing.T) {
		if err := repo.SaveOne(blogPost); err != nil {
			t.Errorf("Cannot save article: %v %v", blogPost, err)
		}
		actualArticle, err := repo.Get(blogPost.Id)
		if err != nil {
			t.Error(errors.WithStack(err))
		}
		if !blogPost.Compare(*actualArticle) {
			t.Errorf("article mistmatch\n Expected: %v\n Actual: %v\n", blogPost, actualArticle)
		}

	})

	t.Run("save all articles", func(t *testing.T) {
		ctx := context.Background()
		if err := repo.SaveAll(ctx, blogPosts); err != nil {
			t.Error("Cannot save all articles", err)
		}

		for _, expectedArticle := range blogPosts {
			actualArticle, err := repo.Get(expectedArticle.Id)
			if err != nil {
				t.Errorf("Cannot get article %s: %v", expectedArticle.Id, err)
			}
			if !expectedArticle.Compare(*actualArticle) {
				t.Errorf("article mismatch\n Expected: %v\n Actual: %v\n", blogPost, actualArticle)
			}

		}

	})

}
