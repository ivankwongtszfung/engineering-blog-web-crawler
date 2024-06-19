package sqlite

import (
	"context"
	"database/sql"

	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog"

	"github.com/pkg/errors"
)

type ArticleRepository struct {
	DB *sql.DB
}

const PREPARE_STATEMENT_ERROR = "Cannot create prepare statement"

func (repo ArticleRepository) Create() error {
	sqlStmt := `
	create table IF NOT EXISTS article (
		id text not null primary key,
		title text,
		category text,
		date text,
		blogURL text,
		imageURL text
	);
	`
	_, err := repo.DB.Exec(sqlStmt)
	if err != nil {
		return errors.Wrap(err, "Cannot create article table")
	}
	return nil
}

func (repo ArticleRepository) Get(id string) (*blog.Article, error) {
	query := "SELECT id, title, category, date, blogurl, imageurl FROM article WHERE id = ?"
	row := repo.DB.QueryRow(query, id)

	var article blog.Article
	err := row.Scan(&article.Id, &article.Title, &article.Category, &article.Date, &article.BlogURL, &article.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("id: %s is not found", id)
		}
		return nil, err
	}

	return &article, nil
}

func (repo ArticleRepository) All() ([]*blog.Article, error) {
	query := "SELECT id, title, category, date, blogurl, imageurl FROM article"
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "Error when running the get all query")
	}
	defer rows.Close()

	var articles []*blog.Article
	for rows.Next() {
		var article blog.Article
		err := rows.Scan(&article.Id, &article.Title, &article.Category, &article.Date, &article.BlogURL, &article.ImageURL)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error when getting the rows")
	}

	return articles, nil

}

func (repo ArticleRepository) SaveOne(article blog.Article) error {
	stmt, err := repo.DB.Prepare("insert into article(id, title, category, date, blogURL, imageURL) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, PREPARE_STATEMENT_ERROR)
	}
	defer stmt.Close()

	_, err = stmt.Exec(&article.Id, &article.Title, &article.Category, &article.Date, &article.BlogURL, &article.ImageURL)
	if err != nil {
		return errors.Wrap(err, "Cannot insert article")
	}
	return nil
}

func (repo ArticleRepository) SaveAll(ctx context.Context, articles []blog.Article) error {
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Cannot create transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw after rollback
		} else if err != nil {
			tx.Rollback() // rollback if any error occurred
		}
	}()

	stmt, err := tx.Prepare("insert into article(id, title, category, date, blogURL, imageURL) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "Cannot create prepare statement")
	}
	defer stmt.Close()

	for _, article := range articles {
		_, err = stmt.Exec(article.Id, article.Title, article.Category, article.Date, article.BlogURL, article.ImageURL)
		if err != nil {
			return errors.WithStack(errors.Wrapf(err, "Cannot execute the article statement, %+v", article))
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.WithStack(errors.Wrap(err, "Cannot commit the save all article function"))
	}

	return nil
}
