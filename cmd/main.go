package main

import (
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ivankwongtszfung/engineering-blog-web-crawler/config"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog/uber"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/infra/repo"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/infra/repo/sqlite"
	"github.com/ivankwongtszfung/engineering-blog-web-crawler/pkg/kvstore"
	"github.com/pkg/errors"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
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
				log.Print(errors.Wrap(err, "Error extracting article"))
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
		log.Printf("Link found: %q -> %s\n", e.Text, link)
		kv.Set(e.Request.AbsoluteURL(link), "")
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://www.uber.com/en-CA/blog/engineering/")
	c.Wait()
}

func main() {

	// kv store init
	var kv kvstore.KVStore = kvstore.NewRedisStore(config.REDIS_DATABASE)
	if err := kv.Ping(); err != nil {
		log.Fatalln("failed to connect to redis", err)
	}
	defer kv.Close()

	// article channel init
	articles := make(chan *blog.Article)
	go func() {
		log.Print("close the article")
		defer close(articles)
		scrapeUberBlogs(kv, articles)
	}()

	// db init
	db, err := sql.Open(config.DB_DRIVER, config.DB_HOST)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	var articleRepo repo.IArticleRepository = sqlite.ArticleRepository{DB: db}
	articleRepo.Create()

	// persist in db
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for article := range articles {
				err = articleRepo.SaveOne(*article)
				if err != nil {
					log.Print(err)
					kv.Delete(article.BlogURL)
				}
			}
		}()
	}
	wg.Wait()
}
