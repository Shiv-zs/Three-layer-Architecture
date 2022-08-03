package book

import (
	"database/sql/driver"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"mytest/models"
)

// Test_Post Book
func Test_Post(t *testing.T) {
	testcases := []struct {
		desc         string
		req          models.Book
		response     models.Book
		lastInsertID int64
		rowAffected  int64
		err          error
	}{
		{desc: "valid details", req: models.Book{BookID: 1, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, response: models.Book{BookID: 1,
			AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, lastInsertID: 1, rowAffected: 1},
		{desc: "duplicate id", req: models.Book{BookID: 1, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, err: errors.New(" Duplicate entry '1' for key 'PRIMARY'")},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// Mocking insert query for book
		mock.ExpectExec("insert into Book(bookId,title,authorId,Publication,PublishedDate) values (?,?,?,?,?)").
			WithArgs(v.req.BookID, v.req.Title, v.req.AuthorID, v.req.Publication, v.req.PublishedDate).
			WillReturnResult(sqlmock.NewResult(v.lastInsertID, v.rowAffected)).WillReturnError(v.err)

		// injecting mock db
		datastore := New()

		resp, err := datastore.Post(ctx, &v.req)

		// comparing body
		if !reflect.DeepEqual(resp, v.response) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.response)
		}

		// comparing error
		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// Test_GetAll all book
func Test_GetAll(t *testing.T) {
	testcases := []struct {
		desc string
		resp []models.Book
		rows *sqlmock.Rows
		err  error
	}{
		{desc: "valid details ", resp: []models.Book{
			{BookID: 1, AuthorID: 1,
				Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"},
			{BookID: 2, AuthorID: 1,
				Title: "3 States", Publication: "Penguin", PublishedDate: "11/03/2016"}},
			rows: sqlmock.NewRows([]string{"bookId", "title", "authorId", "Publication", "PublishedDate"}).
				AddRow(1, "States", 1, "Scholastic", "16/03/2016").
				AddRow(2, "3 States", 1, "Penguin", "11/03/2016"),
		},
		{desc: "error in scanning", resp: []models.Book{}, rows: sqlmock.NewRows([]string{"bookId", "title",
			"authorId", "Publication", "PublishedDate"}).AddRow("abc", "States", 1, "Scholastic", "16/03/2016"),
		},
		{desc: "error in select all ", rows: sqlmock.NewRows([]string{}), err: errors.New("error in select all")},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// Mocking select all books query
		mock.ExpectQuery("SELECT * FROM Book").WillReturnRows(v.rows).WillReturnError(v.err)

		// injecting mock db
		datastore := New()

		resp, err := datastore.GetAll(ctx)

		// Comparing body
		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("Desc : %v,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// Test_GetByID Testing book Get by id
func Test_GetByID(t *testing.T) {
	testcases := []struct {
		desc string
		id   int
		resp models.Book
		rows *sqlmock.Rows
		err  error
	}{
		{desc: "valid", id: 1, resp: models.Book{BookID: 1, AuthorID: 1, Title: "States",
			Publication: "Scholastic", PublishedDate: "16/03/2016"}, rows: sqlmock.NewRows([]string{"bookId", "title", "authorId", "Publication", "PublishedDate"}).
			AddRow(1, "States", 1, "Scholastic", "16/03/2016")},
		{desc: "error in scanning", id: 11, rows: sqlmock.NewRows([]string{"bookId", "title", "authorId", "Publication", "PublishedDate"}).
			AddRow("ac", "States", 1, "Scholastic", "16/03/2016"),
			err: errors.New("error in scanning")},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// Mocking Query for reading book
		mock.ExpectQuery("select * from Book where bookId=?").WithArgs(v.id).
			WillReturnRows(v.rows).WillReturnError(v.err)

		// Injecting mock DB
		datastore := New()

		resp, err := datastore.GetByID(ctx, v.id)

		// Comparing body
		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("Desc : %v,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		// Comparing errors
		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// Test_Put book
func Test_Put(t *testing.T) {
	testcases := []struct {
		desc   string
		id     int
		req    models.Book
		resp   models.Book
		result driver.Result
		err    error
	}{
		{desc: "valid", id: 1, req: models.Book{BookID: 1, AuthorID: 1,
			Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"}, result: sqlmock.NewResult(1, 1),
			resp: models.Book{BookID: 1, AuthorID: 1, Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"}},
		{desc: "error in exec", id: 11, req: models.Book{BookID: 1, AuthorID: 1,
			Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"}, result: sqlmock.NewResult(0, 0),
			err: errors.New("sql: no rows in result set")},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// Mocking Exec query for updating data
		mock.ExpectExec("UPDATE Book SET title=?, Publication=? , PublishedDate=? WHERE bookId=?").
			WithArgs(v.resp.Title, v.resp.Publication, v.resp.PublishedDate, v.id).
			WillReturnResult(v.result).WillReturnError(v.err)

		// Injecting mock Db
		datastore := New()

		resp, err := datastore.Update(ctx, v.id, &v.req)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// Test_Delete book
func Test_Delete(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		rowAffected int
		result      driver.Result
		err         error
	}{
		{desc: "valid", id: 1, rowAffected: 1, result: sqlmock.NewResult(1, 1)},
		{desc: "error in exec", id: 1, rowAffected: 0, result: sqlmock.NewResult(0, 0), err: errors.New("error in exec")},
		{desc: "error in rowAffected", id: 11, result: sqlmock.NewErrorResult(errors.New("sql: no rows in result set"))},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// Mocking delete query from book
		mock.ExpectExec("DELETE FROM Book where bookId=?").WithArgs(v.id).
			WillReturnResult(v.result).WillReturnError(v.err)

		// Injecting mock DB
		datastore := New()

		resp, err := datastore.Delete(ctx, v.id)

		// Comparing body
		if !reflect.DeepEqual(resp, v.rowAffected) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.rowAffected)
		}

		// Comparing errors
		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// Test_GetByTitle Testing book Get by id
func Test_GetByTitle(t *testing.T) {
	testcases := []struct {
		desc  string
		title string
		resp  []models.Book
		rows  *sqlmock.Rows
		err   error
	}{
		{desc: "valid details ", title: "States", resp: []models.Book{
			{BookID: 1, AuthorID: 1,
				Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"},
			{BookID: 2, AuthorID: 1,
				Title: "States", Publication: "Penguin", PublishedDate: "11/03/2016"}},
			rows: sqlmock.NewRows([]string{"bookId", "title", "authorId", "Publication", "PublishedDate"}).
				AddRow(1, "States", 1, "Scholastic", "16/03/2016").
				AddRow(2, "States", 1, "Penguin", "11/03/2016"),
		},
		{desc: "error in scanning", title: "States", resp: []models.Book{}, rows: sqlmock.NewRows([]string{"bookId", "title",
			"authorId", "Publication", "PublishedDate"}).AddRow("abc", "States", 1, "Scholastic", "16/03/2016"),
		},
		{desc: "error in select all ", resp: []models.Book{}, rows: sqlmock.NewRows([]string{}), err: errors.New("error in select all by title")},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// Mocking select all books query
		mock.ExpectQuery("select * from Book where title=?").WithArgs(v.title).
			WillReturnRows(v.rows).WillReturnError(v.err)

		// injecting mock db
		datastore := New()

		resp, err := datastore.GetBookByTitle(ctx, v.title)

		// Comparing body
		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("Desc : %v,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// Test_IsBookPresent is to check for book existence
func Test_IsBookPresent(t *testing.T) {
	testCases := []struct {
		desc string
		id   int
		rows *sqlmock.Rows
		resp bool
		err  error
	}{
		{desc: "valid", id: 1, resp: false, rows: sqlmock.NewRows([]string{"bookId", "title", "authorId", "Publication", "PublishedDate"}).
			AddRow(1, "States", 1, "Scholastic", "16/03/2016")},
		{desc: "error in scanning", id: 10, rows: sqlmock.NewRows([]string{}), resp: true},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("error mocking:%v", err)
	}

	app := gofr.New()
	app.DB().DB = db

	ctx := gofr.NewContext(nil, nil, app)

	defer db.Close()

	for i, v := range testCases {
		mock.ExpectQuery("select * from Book where bookId=?").WillReturnRows(v.rows).WillReturnError(v.err)
		datastore := New()

		resp := datastore.IsBookPresent(ctx, v.id)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}
	}
}
