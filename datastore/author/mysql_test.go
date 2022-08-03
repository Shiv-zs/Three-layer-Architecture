package author

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

// Testing Post Author
func TestAuthor_Post(t *testing.T) {
	testcases := []struct {
		desc         string
		req          models.Author
		resp         models.Author
		lastInsertID int64
		rowAffected  int64
		err          error
	}{
		{desc: "valid details", req: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, resp: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, lastInsertID: 1, rowAffected: 1},
		{desc: "duplicate id", req: models.Author{AuthID: 3, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, lastInsertID: 0, rowAffected: 0, err: errors.New("error")},
	}

	// Customize SQL query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Closing DB after all things done
	defer db.Close()

	for i, v := range testcases {
		// mocking insert exec query
		mock.ExpectExec("insert into Author(authorId,firstName,lastName,dob,penName) values (?,?,?,?,?)").
			WithArgs(v.req.AuthID, v.req.FirstName, v.req.LastName, v.req.Dob, v.req.PenName).
			WillReturnResult(sqlmock.NewResult(v.lastInsertID, v.rowAffected)).WillReturnError(v.err)

		app := gofr.New()
		app.DB().DB = db

		ctx := gofr.NewContext(nil, nil, app)

		datastore := New()

		resp, err := datastore.Post(ctx, v.req)

		// Comparing body
		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		// Comparing errors
		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

// Testing Put Author
func TestAuthor_Put(t *testing.T) {
	testcases := []struct {
		desc string
		id   int
		resp models.Author
		res  driver.Result
		err  error
	}{
		{desc: "valid", id: 1, resp: models.Author{AuthID: 1, FirstName: "Rajan",
			LastName: "Sharma", Dob: "26/04/2001", PenName: "Rajan"}, res: sqlmock.NewResult(1, 1)},
		{desc: "id not exist", id: 11, res: sqlmock.NewResult(0, 0), err: errors.New("error")},
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
		// Mocking Exec for updating data
		mock.ExpectExec("UPDATE Author SET firstName=?, lastName=? , dob=? , penName=? WHERE authorId=?").
			WithArgs(v.resp.FirstName, v.resp.LastName, v.resp.Dob, v.resp.PenName, v.id).WillReturnResult(v.res).WillReturnError(v.err)

		datastore := New()

		resp, err := datastore.Update(ctx, v.id, v.resp)

		// Comparing body
		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		// Getting error
		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

// Testing Delete Author
func TestAuthor_Delete(t *testing.T) {
	testcases := []struct {
		desc string
		id   int
		resp int
		res  driver.Result
		err  error
	}{
		{desc: "valid", id: 1, resp: 1, res: sqlmock.NewResult(0, 1)},
		{desc: "id not exist", id: 11, res: sqlmock.NewResult(0, 0), err: errors.New("error")},
		{desc: "inserted id error", id: 4, res: sqlmock.NewErrorResult(errors.New("error"))},
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
		// Mocking Exec for deleting author
		mock.ExpectExec("delete from Author where authorId=?").WithArgs(v.id).
			WillReturnResult(v.res).WillReturnError(v.err)

		datastore := New()

		resp, err := datastore.Delete(ctx, v.id)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

// TestIncludeAuthor
func TestIncludeAuthor(t *testing.T) {
	testcases := []struct {
		desc string
		id   int
		resp models.Author
		rows *sqlmock.Rows
		err  error
	}{
		{desc: "valid", id: 1, resp: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, rows: sqlmock.NewRows([]string{"authorId", "fistName", "lastName", "dob", "penName"}).
			AddRow(1, "Chetan", "Bhagat", "06/04/2001", "Chetan")},
		{desc: "id not exist", id: 11, resp: models.Author{}, rows: sqlmock.NewRows([]string{"authorId", "fistName", "lastName", "dob", "penName"})},
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
		// Mocking Exec for deleting author
		mock.ExpectQuery("select * from Author where authorId=?").WithArgs(v.id).
			WillReturnRows(v.rows).WillReturnError(v.err)

		datastore := New()

		resp, err := datastore.IncludeAuthor(ctx, v.id)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

func Test_IsAuthorIDPresent(t *testing.T) {
	testCases := []struct {
		desc string
		id   int
		rows *sqlmock.Rows
		resp bool
		err  error
	}{
		{desc: "valid", id: 1, resp: false, rows: sqlmock.NewRows([]string{"authorId", "fistName", "lastName", "dob", "penName"}).
			AddRow(1, "Chetan", "Bhagat", "06/04/2001", "Chetan")},
		{desc: "id not exist", id: 10, rows: sqlmock.NewRows([]string{}), resp: true},
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
		mock.ExpectQuery("select * from Author where authorId=?").WithArgs(v.id).WillReturnRows(v.rows).WillReturnError(v.err)
		datastore := New()

		resp := datastore.IsAuthorIDPresent(ctx, v.id)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}
	}
}
