package book

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"strconv"

	"mytest/models"
	"mytest/service"
)

type Delivery struct {
	service service.Book
}

func New(book service.Book) Delivery {
	return Delivery{service: book}
}

// Create method is post details of Book
func (d Delivery) Create(c *gofr.Context) (interface{}, error) {
	var book models.Book

	if err := c.Bind(&book); err != nil {
		return models.Book{}, err
	}

	return d.service.Post(c, &book)
}

// GetAll method is get all details of Books
func (d Delivery) GetAll(c *gofr.Context) (interface{}, error) {
	title := c.Param("title")
	includeAuthor := c.Param("includeAuthor")

	// Getting all books
	return d.service.GetAll(c, title, includeAuthor)
}

// GetByID method is get the book by its id
func (d Delivery) GetByID(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	if id == "" {
		return models.Book{}, errors.MissingParam{Param: []string{id}}
	}

	id2, err := strconv.Atoi(id)
	if err != nil {
		return models.Book{}, errors.InvalidParam{Param: []string{id}}
	}

	return d.service.GetByID(c, id2)
}

// Update method is to update details of Book
func (d Delivery) Update(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	if id == "" {
		return models.Book{}, errors.MissingParam{Param: []string{id}}
	}

	id2, err := strconv.Atoi(id)
	if err != nil {
		return models.Book{}, errors.InvalidParam{Param: []string{id}}
	}

	var book models.Book

	if err := c.Bind(&book); err != nil {
		return models.Book{}, err
	}

	return d.service.Update(c, id2, &book)
}

// Delete method is to delete details of Book by its id
func (d Delivery) Delete(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	if id == "" {
		return 0, errors.MissingParam{Param: []string{id}}
	}

	id2, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.InvalidParam{Param: []string{id}}
	}

	return d.service.Delete(c, id2)
}
