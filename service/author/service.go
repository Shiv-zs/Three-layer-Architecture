package author

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"mytest/datastore"
	"mytest/models"
	"strconv"
)

type Service struct {
	datastore datastore.Author
}

func New(author datastore.Author) Service {
	return Service{author}
}

// Post Author details
func (s Service) Post(c *gofr.Context, auth models.Author) (models.Author, error) {
	// Checking for invalid id
	if auth.AuthID <= 0 {
		return models.Author{}, errors.Error("invalid id")
	}

	if isMissingFields(auth) {
		return models.Author{}, errors.Error("missing fields")
	}

	author, err := s.datastore.Post(c, auth)
	if err != nil {
		return models.Author{}, err
	}

	return author, nil
}

// Update Author details
func (s Service) Update(c *gofr.Context, id int, auth models.Author) (models.Author, error) {
	// Checking invalid id
	if id <= 0 {
		return models.Author{}, errors.Error("invalid id")
	}

	if isMissingFields(auth) {
		return models.Author{}, errors.Error("missing fields")
	}

	// Checking author ID present or not
	check := s.datastore.IsAuthorIDPresent(c, id)
	if check {
		return models.Author{}, errors.EntityNotFound{Entity: "Author", ID: strconv.Itoa(id)}
	}

	author, err := s.datastore.Update(c, id, auth)
	if err != nil {
		return models.Author{}, err
	}

	return author, nil
}

// Delete Author by its ID
func (s Service) Delete(c *gofr.Context, id int) (int, error) {
	// Checking for invalid id
	if id <= 0 {
		return 0, errors.Error("invalid id")
	}

	// Checking author ID present or not
	check := s.datastore.IsAuthorIDPresent(c, id)

	if check {
		return 0, errors.Error("author id is not valid")
	}

	rowAffected, err := s.datastore.Delete(c, id)
	if err != nil {
		return 0, err
	}

	return rowAffected, nil
}

func isMissingFields(auth models.Author) bool {
	if auth.FirstName == "" || auth.LastName == "" || auth.PenName == "" || auth.Dob == "" {
		return true
	}

	return false
}
