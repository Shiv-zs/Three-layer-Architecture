package author

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"mytest/models"
)

type Datastore struct {
}

func New() Datastore {
	return Datastore{}
}

// Post method is to post the data in Author table
func (d Datastore) Post(c *gofr.Context, auth models.Author) (models.Author, error) {
	// inserting data into db
	_, err := c.DB().Exec("insert into Author(authorId,firstName,lastName,dob,penName) values (?,?,?,?,?)",
		auth.AuthID, auth.FirstName, auth.LastName, auth.Dob, auth.PenName)
	if err != nil {
		return models.Author{}, err
	}

	return auth, nil
}

// Update method is to update the data in Author table
func (d Datastore) Update(c *gofr.Context, id int, auth models.Author) (models.Author, error) {
	_, err := c.DB().Exec("UPDATE Author SET firstName=?, lastName=? , dob=? , penName=? WHERE authorId=?",
		auth.FirstName, auth.LastName, auth.Dob, auth.PenName, id)
	if err != nil {
		return models.Author{}, err
	}

	return auth, nil
}

// Delete method is to delete the data in Author
func (d Datastore) Delete(c *gofr.Context, id int) (int, error) {
	res, err := c.DB().Exec("delete from Author where authorId=?", id)
	if err != nil {
		return 0, err
	}

	rowAffected, err2 := res.RowsAffected()
	if err2 != nil {
		return 0, err2
	}

	return int(rowAffected), nil
}

// IncludeAuthor details by its ID
func (d Datastore) IncludeAuthor(c *gofr.Context, id int) (models.Author, error) {
	row := c.DB().QueryRow("select * from Author where authorId=?", id)

	var author models.Author

	if err := row.Scan(&author.AuthID, &author.FirstName, &author.LastName, &author.Dob, &author.PenName); err != nil {
		return models.Author{}, err
	}

	return author, nil
}

// IsAuthorIDPresent method is to check weather author is present in DB or not
func (d Datastore) IsAuthorIDPresent(c *gofr.Context, id int) bool {
	var author models.Author

	row := c.DB().QueryRow("select * from Author where authorId=?", id)

	if err := row.Scan(&author.AuthID, &author.FirstName, &author.LastName, &author.Dob, &author.PenName); err != nil {
		return true
	}

	return false
}
