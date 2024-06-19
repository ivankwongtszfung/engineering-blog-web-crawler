package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog/uber"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/infra/repo"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/infra/repo/sqlite"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/pkg/kvstore"

	"github.com/gocolly/colly"
)

const ENGINEERING_PAGE_PREFIX string = "https://www.uber.com/en-CA/blog/engineering/page/"

func scrapeUberBlogs(kv kvstore.KVStore, articles chan<- *blog.Article) {
	// Instantiate default collector
	c := colly.NewCollector(colly.MaxDepth(1), colly.Async(true))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*uber.*",
		Parallelism: 3,
		RandomDelay: 2 * time.Second,
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html")
		log.Println("Visiting", r.URL.String())
		log.Println("Visiting", r.Headers)
	})

	c.OnHTML("div[data-baseweb]", func(e *colly.HTMLElement) {
		className := e.Attr("data-baseweb")
		if className == "flex-grid-item" {
			article, err := uber.ExtractArticle(e)
			if err != nil {
				fmt.Printf("Error extracting article: %v\n", err)
				return
			}
			if val, _ := kv.Exist(article.BlogURL); val {
				return
			}
			kv.Set(article.BlogURL, "")
			log.Printf("Found New article: %v\n", article.Title)
			articles <- &article.Article
		}
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = e.Request.AbsoluteURL(link)
		if !strings.HasPrefix(link, ENGINEERING_PAGE_PREFIX) {
			return
		}
		if val, _ := kv.Exist(link); val {
			return
		}
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		kv.Set(e.Request.AbsoluteURL(link), "")
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://www.uber.com/en-CA/blog/engineering/")
	c.Wait()
	fmt.Print("close the article")
}

func main() {

	kv := kvstore.NewRedisStore("localhost:6379")
	if err := kv.Ping(); err != nil {
		log.Fatalln("failed to connect to redis", err)
	}

	articles := make(chan *blog.Article)
	go func() {
		defer close(articles)
		scrapeUberBlogs(kv, articles)
	}()

	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	var articleRepo repo.IArticleRepository = sqlite.ArticleRepository{DB: db}

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for article := range articles {
				articleRepo.SaveOne(*article)
			}
		}()
	}
	wg.Wait()
}
