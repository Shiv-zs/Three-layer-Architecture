package book

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"log"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"mytest/datastore"
	"mytest/models"
)

var author = models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"}

// TestBook_Post function is to test post author details
func TestBook_Post(t *testing.T) {
	testcases := []struct {
		desc             string
		req              models.Book
		response         models.Book
		includeAuthorErr error
		PostErr          error
	}{
		{desc: "valid details", req: models.Book{BookID: 1, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"},
			response: models.Book{BookID: 1, AuthorID: 1,
				Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
				Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, includeAuthorErr: nil, PostErr: nil},
		{desc: "invalid id", req: models.Book{BookID: -11, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, response: models.Book{}, includeAuthorErr: nil, PostErr: nil},
		{desc: "invalid publication", req: models.Book{BookID: 1, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Lenin", PublishedDate: "16/03/2016"}, response: models.Book{}, includeAuthorErr: nil, PostErr: nil},
		{desc: "invalid publishedDate : year", req: models.Book{BookID: 2, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2061"}},
		{desc: "invalid publishedDate : month", req: models.Book{BookID: 7, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/33/2016"}},
		{desc: "invalid publishedDate : day", req: models.Book{BookID: 8, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "116/03/2011"}},
		{desc: "missing title", req: models.Book{BookID: 1, AuthorID: 1,
			Auth:        models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Publication: "Scholastic", PublishedDate: "16/03/2016"}, response: models.Book{}, includeAuthorErr: nil, PostErr: nil},
		{desc: "missing publication", req: models.Book{BookID: 3, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", PublishedDate: "16/03/2016"}, PostErr: errors.Error("missing book fields")},
		{desc: "missing publishedDate", req: models.Book{BookID: 4, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic"}, response: models.Book{}, includeAuthorErr: nil, PostErr: nil},
		{desc: "error in includeAuthor", req: models.Book{BookID: 5, AuthorID: 11,
			Auth:  models.Author{AuthID: 11, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, response: models.Book{}, PostErr: nil,
			includeAuthorErr: errors.Error("error in includeAuthor")},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockAuthor.EXPECT().IncludeAuthor(c, v.req.AuthorID).Return(v.response.Auth, v.includeAuthorErr).AnyTimes()
		mockBook.EXPECT().Post(c, &v.req).Return(v.response, v.PostErr).AnyTimes()

		resp, err := service.Post(c, &v.req)

		if !reflect.DeepEqual(resp, v.response) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.response)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err)
		}
	}
}

// TestBook_PostErr function is to test post author details
func TestBook_PostErr(t *testing.T) {
	testcases := []struct {
		desc             string
		req              models.Book
		response         models.Book
		includeAuthorErr error
		PostErr          error
	}{
		{desc: "error in datastore post", req: models.Book{BookID: 6, AuthorID: 21,
			Auth:  models.Author{AuthID: 21, FirstName: "Chetan", LastName: "Sharma", Dob: "26/04/2001", PenName: "Sharma"},
			Title: "3 States", Publication: "Scholastic", PublishedDate: "26/03/2016"}, response: models.Book{},
			PostErr: errors.Error("error in post"), includeAuthorErr: nil},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockAuthor.EXPECT().IncludeAuthor(c, v.req.AuthorID).Return(v.response.Auth, v.includeAuthorErr).AnyTimes()
		mockBook.EXPECT().Post(c, &v.req).Return(v.response, v.PostErr).AnyTimes()

		resp, err := service.Post(c, &v.req)

		if !reflect.DeepEqual(resp, v.response) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.response)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err)
		}
	}
}

// TestBook_GetAll function is to test for getting all books
func TestBook_GetAll(t *testing.T) {
	testcases := []struct {
		desc          string
		title         string
		includeAuthor string
		resp          []models.Book
		getTitleErr   error
		getAuthorErr  error
		getAllErr     error
	}{

		{
			"valid details ", "", "false",
			[]models.Book{{BookID: 1, AuthorID: 1,
				Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
				Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"}},
			nil, nil, nil,
		},
		{
			"get by title", "States", "false",
			[]models.Book{{BookID: 2, AuthorID: 1, Auth: models.Author{}, Title: "States", Publication: "Scholastic",
				PublishedDate: "16/03/2016"}},
			nil,
			nil,
			nil,
		},
		{
			"error in get title", "StateOfAmerica", "false",
			[]models.Book{}, errors.Error("error in get title"),
			nil, nil,
		},
		{
			"include Author", "Enough", "true",
			[]models.Book{{BookID: 3, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Chetan",
				LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"}, Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"}},
			nil, nil, nil,
		},
		{
			"error in include Author", "Village", "true",
			[]models.Book{{4, 1, models.Author{1, "shiv",
				"Bhagat", "06/04/1990", "shiv"}, "nothing", "Scholastic", "16/03/2016"}},
			nil, errors.Error("error in includeAuthor"), nil,
		},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockBook.EXPECT().GetBookByTitle(c, v.title).Return(v.resp, v.getTitleErr).AnyTimes()
		mockAuthor.EXPECT().IncludeAuthor(c, author.AuthID).Return(author, v.getAuthorErr).AnyTimes()
		mockBook.EXPECT().GetAll(c).Return(v.resp, v.getAllErr).AnyTimes()

		resp, err := service.GetAll(c, v.title, v.includeAuthor)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("Desc : %v,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

// TestBook_GetAllIncludeAuthorErr function is to test for getting all books
func TestBook_GetAllIncludeAuthorErr(t *testing.T) {
	testcases := []struct {
		desc          string
		title         string
		includeAuthor string
		resp          []models.Book
		getTitleErr   error
		getAuthorErr  error
		getAllErr     error
	}{
		{
			"error in include Author", "Village", "true",
			[]models.Book{},
			nil, errors.Error("error in includeAuthor"), nil,
		},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockBook.EXPECT().GetBookByTitle(c, v.title).Return(v.resp, v.getTitleErr).AnyTimes()
		mockAuthor.EXPECT().IncludeAuthor(c, author.AuthID).Return(author, v.getAuthorErr).AnyTimes()
		mockBook.EXPECT().GetAll(c).Return(v.resp, v.getAllErr).AnyTimes()

		resp, err := service.GetAll(c, v.title, v.includeAuthor)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("Desc : %v,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

// TestBook_GetByID function is to test for get a book
func TestBook_GetByID(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		resp      models.Book
		checkBook bool
		isBookErr error
		getIdErr  error
	}{

		{desc: "valid detail", id: 1, resp: models.Book{BookID: 1, AuthorID: 1,
			Auth:  models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, checkBook: false, isBookErr: nil, getIdErr: nil},
		{desc: "invalid id", id: -11, getIdErr: errors.Error("invalid id"), resp: models.Book{}, checkBook: false, isBookErr: nil},
		{desc: "error in IsBookPresent", id: 15, resp: models.Book{}, isBookErr: errors.Error("error in isBook"), checkBook: true, getIdErr: nil},
		{desc: "error in datastore get", id: 23, resp: models.Book{}, getIdErr: errors.Error("error in Get"), checkBook: false, isBookErr: nil},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockBook.EXPECT().IsBookPresent(c, v.id).Return(v.checkBook).AnyTimes()
		mockBook.EXPECT().GetByID(c, v.id).Return(v.resp, v.getIdErr).AnyTimes()

		resp, err := service.GetByID(c, v.id)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("Desc : %v,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err)
		}
	}
}

// TestBook_Put function is to test for updating Book
func TestBook_Put(t *testing.T) {
	testcases := []struct {
		desc             string
		id               int
		req              models.Book
		resp             models.Book
		checkBook        bool
		includeAuthorErr error
		putErr           error
	}{
		{
			desc: "valid ",
			id:   1,
			req: models.Book{BookID: 1, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"},
			resp: models.Book{BookID: 1, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"},
			checkBook:        false,
			putErr:           nil,
			includeAuthorErr: nil,
		},
		{
			desc: "invalid id",
			req: models.Book{BookID: 3, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "Arihant", PublishedDate: "17/03/2016"},
			resp:             models.Book{},
			id:               -11,
			checkBook:        false,
			putErr:           nil,
			includeAuthorErr: nil,
		},
		{
			desc: "invalid publication",
			id:   4,
			req: models.Book{BookID: 4, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "lenin", PublishedDate: "17/03/2016"},
			resp:             models.Book{},
			checkBook:        false,
			putErr:           nil,
			includeAuthorErr: nil,
		},
		{
			desc: "invalid publishedDate",
			id:   5,
			req: models.Book{BookID: 5, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2061"},
			resp:             models.Book{},
			checkBook:        false,
			putErr:           nil,
			includeAuthorErr: nil,
		},
		{
			desc: "missing book fields",
			id:   6,
			req: models.Book{BookID: 6, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Publication: "lenin", PublishedDate: "17/03/2016"},
			resp:             models.Book{},
			checkBook:        false,
			putErr:           nil,
			includeAuthorErr: nil,
		},
		{
			desc: "error in isBookPresent",
			id:   7,
			req: models.Book{BookID: 7, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"},
			resp:             models.Book{},
			checkBook:        true,
			putErr:           nil,
			includeAuthorErr: nil,
		},
		{
			desc: "error in put Author",
			id:   9,
			req: models.Book{BookID: 9, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Singh",
				Dob: "07/04/2001", PenName: "Gaurav"}, Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"},
			resp:             models.Book{},
			checkBook:        false,
			putErr:           errors.Error("error in putAuthor"),
			includeAuthorErr: nil,
		},
		{
			"error in includeAuthor",
			8,
			models.Book{BookID: 8, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Gaurav", LastName: "Chandra",
				Dob: "07/04/2001", PenName: "Chandra"}, Title: "Days", Publication: "Penguin", PublishedDate: "17/03/2016"},
			models.Book{},
			false,
			errors.Error("error in include author"),
			nil,
		},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context
		mockBook.EXPECT().IsBookPresent(c, v.id).Return(v.checkBook).AnyTimes()
		mockAuthor.EXPECT().IncludeAuthor(c, author.AuthID).Return(author, v.includeAuthorErr).AnyTimes()
		mockBook.EXPECT().Update(c, v.id, &v.req).Return(v.resp, v.putErr).AnyTimes()

		resp, err := service.Update(c, v.id, &v.req)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}

// TestBook_Delete function is to test for deleting a valid book
func TestBook_Delete(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		rowAffected int
		checkBook   bool
		err         error
	}{
		{desc: "valid", id: 1, rowAffected: 1, checkBook: false, err: nil},
		{desc: "invalid id", id: -2, rowAffected: 0, checkBook: false, err: nil},
		{desc: "error isBookPresent", id: 3, rowAffected: 0, checkBook: true, err: nil},
		{desc: "error in datastore delete", id: 4, rowAffected: 0, checkBook: false, err: errors.Error("error in delete")},
	}

	ctr := gomock.NewController(t)
	mockBook := datastore.NewMockBook(ctr)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockBook, mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockBook.EXPECT().IsBookPresent(c, v.id).Return(v.checkBook).AnyTimes()
		mockBook.EXPECT().Delete(c, v.id).Return(v.rowAffected, v.err).AnyTimes()

		resp, err := service.Delete(c, v.id)

		if !reflect.DeepEqual(resp, v.rowAffected) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.rowAffected)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err)
		}
	}
}
