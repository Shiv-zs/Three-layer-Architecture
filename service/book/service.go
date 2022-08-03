package book

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"mytest/datastore"
	"mytest/models"

	"strconv"
	"strings"
)

type Service struct {
	datastoreBook   datastore.Book
	datastoreAuthor datastore.Author
}

func New(book datastore.Book, author datastore.Author) Service {
	return Service{book, author}
}

// Post method is to post Book details
func (s Service) Post(c *gofr.Context, book *models.Book) (models.Book, error) {
	if book.BookID <= 0 {
		return models.Book{}, errors.Error("invalid id")
	}

	// missing book fields
	if isBookFieldsMissing(book) {
		return models.Book{}, errors.Error("missing book fields")
	}

	if !isValidPublishedDate(book.PublishedDate) {
		return models.Book{}, errors.Error("invalid publishedDate")
	}

	if !isValidPublication(book.Publication) {
		return models.Book{}, errors.Error("invalid publication")
	}

	var err error

	book.Auth, err = s.datastoreAuthor.IncludeAuthor(c, book.AuthorID)
	if err != nil {
		return models.Book{}, errors.EntityNotFound{Entity: "Author", ID: strconv.Itoa(book.AuthorID)}
	}

	_, err = s.datastoreBook.Post(c, book)
	if err != nil {
		return models.Book{}, err
	}

	return *book, nil
}

// GetByID method is to get Book details by id
func (s Service) GetByID(c *gofr.Context, id int) (models.Book, error) {
	// Checking invalid id
	if id <= 0 {
		return models.Book{}, errors.Error("invalid id")
	}

	check := s.datastoreBook.IsBookPresent(c, id)

	if check {
		return models.Book{}, errors.EntityNotFound{Entity: "Book", ID: strconv.Itoa(id)}
	}

	book, err := s.datastoreBook.GetByID(c, id)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

// Update method is to update Book details
func (s Service) Update(c *gofr.Context, id int, book *models.Book) (models.Book, error) {
	if id <= 0 {
		return models.Book{}, errors.Error("Invalid Id")
	}

	// missing book fields
	if isBookFieldsMissing(book) {
		return models.Book{}, errors.Error("missing book fields")
	}

	if !isValidPublishedDate(book.PublishedDate) {
		return models.Book{}, errors.Error("invalid publishedDate")
	}

	if !isValidPublication(book.Publication) {
		return models.Book{}, errors.Error("invalid publication")
	}

	check := s.datastoreBook.IsBookPresent(c, id)
	if check {
		return models.Book{}, errors.EntityNotFound{Entity: "Book", ID: strconv.Itoa(id)}
	}

	var err error
	book.Auth, err = s.datastoreAuthor.IncludeAuthor(c, book.AuthorID)
	if err != nil {
		return models.Book{}, errors.EntityNotFound{Entity: "Author", ID: strconv.Itoa(id)}
	}

	bk, err2 := s.datastoreBook.Update(c, id, book)
	if err2 != nil {
		return models.Book{}, err2
	}

	return bk, nil
}

// Delete method is to delete Book details
func (s Service) Delete(c *gofr.Context, id int) (int, error) {
	// Checking invalid id
	if id <= 0 {
		return 0, errors.Error("invalid id")
	}

	check := s.datastoreBook.IsBookPresent(c, id)

	if check {
		return 0, errors.EntityNotFound{Entity: "Book", ID: strconv.Itoa(id)}
	}

	rowAffected, err := s.datastoreBook.Delete(c, id)
	if err != nil {
		return 0, err
	}

	return rowAffected, nil
}

// GetAll method is to get the details of book according to title and author details
func (s Service) GetAll(c *gofr.Context, title, includeAuthor string) ([]models.Book, error) {
	// To store book details
	var books []models.Book

	var err error

	if title != "" {
		books, err = s.datastoreBook.GetBookByTitle(c, title)
	} else {
		books, err = s.datastoreBook.GetAll(c)
	}

	if err != nil {
		return []models.Book{}, err
	}

	if includeAuthor == "true" {
		for i := range books {
			author, err := s.datastoreAuthor.IncludeAuthor(c, books[i].AuthorID)
			if err != nil {
				return []models.Book{}, err
			}

			books[i].Auth = author
		}
	}

	return books, nil
}

func isValidPublishedDate(date string) bool {
	p := strings.Split(date, "/")
	day, _ := strconv.Atoi(p[0])
	month, _ := strconv.Atoi(p[1])
	year, _ := strconv.Atoi(p[2])

	switch {
	case day < 0 || day > 31:
		return false
	case month < 0 || month > 12:
		return false
	case year > 2022 || year < 1880:
		return false
	}

	return true
}

func isValidPublication(publication string) bool {
	pub := strings.ToLower(publication)

	if pub == "scholastic" || pub == "arihant" || pub == "penguin" {
		return true
	}

	return false
}

func isBookFieldsMissing(book *models.Book) bool {
	if book.PublishedDate == "" || book.Publication == "" || book.Title == "" {
		return true
	}

	return false
}
