package book

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"mytest/models"
)

type Datastore struct {
}

func New() Datastore {
	return Datastore{}
}

// Post method is to Post data in Book
func (d Datastore) Post(c *gofr.Context, book *models.Book) (models.Book, error) {
	// inserting data into Db
	_, err := c.DB().Exec("insert into Book(bookId,title,authorId,Publication,PublishedDate) values (?,?,?,?,?)",
		book.BookID, book.Title, book.AuthorID, book.Publication, book.PublishedDate)
	if err != nil {
		return models.Book{}, err
	}

	return *book, nil
}

// GetAll method is to get all Books with Author
func (d Datastore) GetAll(c *gofr.Context) ([]models.Book, error) {
	// reading all books from Db
	allRows, err := c.DB().Query("SELECT * FROM Book")
	if err != nil {
		return nil, err
	}

	// Closing db.query
	defer allRows.Close()

	// To store books
	var book []models.Book

	// Iterating to each book
	for allRows.Next() {
		var b models.Book

		err2 := allRows.Scan(&b.BookID, &b.Title, &b.AuthorID, &b.Publication, &b.PublishedDate)
		if err2 != nil {
			return []models.Book{}, err2
		}

		book = append(book, b)
	}

	return book, nil
}

// GetByID method is to get book by its ID
func (d Datastore) GetByID(c *gofr.Context, id int) (models.Book, error) {
	// reading all data of book with given id
	row := c.DB().QueryRow("select * from Book where bookId=?", id)

	// to store d book
	var book models.Book

	// fetching data of book at given id and storing in book
	if err := row.Scan(&book.BookID, &book.Title, &book.AuthorID, &book.Publication, &book.PublishedDate); err != nil {
		return models.Book{}, err
	}

	return book, nil
}

// Update method is to change data of Particular book
func (d Datastore) Update(c *gofr.Context, id int, book *models.Book) (models.Book, error) {
	_, err := c.DB().Exec("UPDATE Book SET title=?, Publication=? , PublishedDate=? WHERE bookId=?",
		book.Title, book.Publication, book.PublishedDate, id)
	if err != nil {
		return models.Book{}, err
	}

	return *book, nil
}

// Delete method is remove Book by its ID
func (d Datastore) Delete(c *gofr.Context, id int) (int, error) {
	res, err := c.DB().Exec("DELETE FROM Book where bookId=?", id)
	if err != nil {
		return 0, err
	}

	rowAffected, err2 := res.RowsAffected()
	if err2 != nil {
		return 0, err
	}

	return int(rowAffected), nil
}

// GetBookByTitle method is to get all the books according to given title
func (d Datastore) GetBookByTitle(c *gofr.Context, title string) ([]models.Book, error) {
	rows, err := c.DB().Query("select * from Book where title=?", title)
	if err != nil {
		return []models.Book{}, err
	}

	// Closing rows
	defer rows.Close()

	var books []models.Book

	// Iterate to all books
	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.BookID, &b.Title, &b.AuthorID, &b.Publication, &b.PublishedDate)
		if err != nil {
			return []models.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

// IsBookPresent method is to find weather a book is present or not
func (d Datastore) IsBookPresent(c *gofr.Context, id int) bool {
	var book models.Book

	row := c.DB().QueryRow("select * from Book where bookId=?", id)

	if err := row.Scan(&book.BookID, &book.Title, &book.AuthorID, &book.Publication, &book.PublishedDate); err != nil {
		return true
	}

	return false
}
