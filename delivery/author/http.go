package author

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"strconv"

	"mytest/models"
	"mytest/service"
)

type Delivery struct {
	service service.Author
}

func New(author service.Author) Delivery {
	return Delivery{author}
}

// Create Request method is to post request
func (d Delivery) Create(c *gofr.Context) (interface{}, error) {
	var author models.Author

	if err := c.Bind(&author); err != nil {
		return models.Author{}, err
	}

	return d.service.Post(c, author)
}

// Update Request method is to update request
func (d Delivery) Update(c *gofr.Context) (interface{}, error) {
	id := c.PathParam("id")

	if id == "" {
		return models.Author{}, errors.MissingParam{Param: []string{id}}
	}

	id2, err := strconv.Atoi(id)
	if err != nil {
		return models.Author{}, errors.InvalidParam{Param: []string{id}}
	}

	var author models.Author

	if err := c.Bind(&author); err != nil {
		return models.Author{}, err
	}

	return d.service.Update(c, id2, author)
}

// Delete method is to delete data from request
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
