package datastore

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"mytest/models"
)

type Book interface {
	Post(c *gofr.Context, book *models.Book) (models.Book, error)
	GetAll(c *gofr.Context) ([]models.Book, error)
	GetByID(c *gofr.Context, id int) (models.Book, error)
	Update(c *gofr.Context, id int, book *models.Book) (models.Book, error)
	Delete(c *gofr.Context, id int) (int, error)
	GetBookByTitle(c *gofr.Context, title string) ([]models.Book, error)
	IsBookPresent(c *gofr.Context, id int) bool
}

type Author interface {
	Post(c *gofr.Context, auth models.Author) (models.Author, error)
	Update(c *gofr.Context, id int, author models.Author) (models.Author, error)
	Delete(c *gofr.Context, id int) (int, error)
	IncludeAuthor(c *gofr.Context, id int) (models.Author, error)
	IsAuthorIDPresent(c *gofr.Context, id int) bool
}
