package book

import (
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"Three-Layer-Architecture/models"
	"Three-Layer-Architecture/service"
)

// TestPostBook function is to test Post book method for posting books
func TestPostBook(t *testing.T) {
	testcase := []struct {
		desc       string
		req        interface{}
		resp       models.Book
		statusCode int
		err        error
	}{
		{desc: "valid details ", req: &models.Book{BookID: 1, AuthorID: 1,
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"},
			resp: models.Book{BookID: 1, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Chetan",
				LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"}},
			statusCode: http.StatusOK, err: nil},
		{desc: "error from svc", req: &models.Book{BookID: -11, AuthorID: 1,
			Title: "2 States", Publication: "Scholastic", PublishedDate: "16/03/2016"},
			statusCode: http.StatusBadRequest, err: nil, resp: models.Book{}},
		{desc: "error in bind", req: &[]models.Book{}, resp: models.Book{},
			statusCode: http.StatusBadRequest, err: nil},
	}

	ctr := gomock.NewController(t)
	mockBook := service.NewMockBook(ctr)
	delivery := New(mockBook)
	k := gofr.New()

	for i, v := range testcase {
		body, err := json.Marshal(v.req)
		if err != nil {
			log.Printf("Not able to marshal : %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/book", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		mockBook.EXPECT().Post(ctx, v.req).Return(v.resp, v.err).AnyTimes()

		book, err2 := delivery.Create(ctx)

		if !reflect.DeepEqual(book, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, book, v.resp)
		}

		if err2 != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v", v.desc, i+1, err2)
		}
	}
}

// TestGetAllBooks function is to test GetAll method for fetching details of books
func TestGetAllBooks(t *testing.T) {
	testcases := []struct {
		desc          string
		title         string
		includeAuthor string
		output        []models.Book
		statusCode    int
		err           error
	}{
		{
			desc: "valid details", output: []models.Book{
				{BookID: 1, AuthorID: 1, Auth: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
					Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"}}, statusCode: http.StatusOK,
			title: "", includeAuthor: "", err: nil,
		},
	}

	ctr := gomock.NewController(t)
	mockBook := service.NewMockBook(ctr)
	delivery := New(mockBook)
	k := gofr.New()

	for i, v := range testcases {
		r := httptest.NewRequest(http.MethodGet, "/books?"+"title="+v.title+"&"+"includeAuthor="+v.includeAuthor, nil)

		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		mockBook.EXPECT().GetAll(ctx, v.title, v.includeAuthor).Return(v.output, v.err).AnyTimes()

		output, err := delivery.GetAll(ctx)

		if !reflect.DeepEqual(output, v.output) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, output, v.output)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v", v.desc, i+1, err)
		}
	}
}

// TestGetBook function is to test GetByID method for fetching a book
func TestGetBook(t *testing.T) {
	testcases := []struct {
		desc       string
		id         string
		resp       models.Book
		statusCode int
		err        error
	}{
		{desc: "valid details", id: "1", resp: models.Book{BookID: 1, AuthorID: 1,
			Title: "States", Publication: "Scholastic", PublishedDate: "16/03/2016"}, statusCode: http.StatusOK, err: nil},
		{desc: "missing param", id: "", resp: models.Book{}, statusCode: http.StatusBadRequest,
			err: errors.Error("missing param")},
		{desc: "invalid param", id: "abc", resp: models.Book{}, statusCode: http.StatusBadRequest,
			err: errors.Error("invalid param")},
	}

	ctr := gomock.NewController(t)
	mockBook := service.NewMockBook(ctr)
	delivery := New(mockBook)
	k := gofr.New()

	for i, v := range testcases {
		params := url.Values{}
		params.Add("bookId", v.id)

		r := httptest.NewRequest(http.MethodGet, "/books?"+params.Encode(), nil)

		w := httptest.NewRecorder()

		r = mux.SetURLVars(r, map[string]string{"id": v.id})

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		id2, err := strconv.Atoi(v.id)
		if err != nil {
			log.Printf("error in converting string : %v", err)
		}

		mockBook.EXPECT().GetByID(ctx, id2).Return(v.resp, v.err).AnyTimes()

		book, err2 := delivery.GetByID(ctx)

		if !reflect.DeepEqual(book, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, book, v.resp)
		}

		if err2 != nil {
			log.Printf("desc : %v, [TEST%d} Got %v\n", v.desc, i+1, err)
		}
	}
}

// TestUpdateBook function is to test Put method for updating details of book
func TestUpdateBook(t *testing.T) {
	testcases := []struct {
		desc       string
		id         string
		req        interface{}
		resp       models.Book
		statusCode int
		err        error
	}{
		{desc: "valid", id: "1", req: &models.Book{BookID: 1, AuthorID: 1, Title: "300 Days", Publication: "Penguin",
			PublishedDate: "17/03/2016"}, resp: models.Book{BookID: 1, AuthorID: 1, Title: "300 Days", Publication: "Penguin",
			PublishedDate: "17/03/2016"}, statusCode: http.StatusOK, err: nil},
		{desc: "missing param", id: "", req: &models.Book{BookID: 2, AuthorID: 1,
			Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"}, err: errors.Error("missing param"),
			statusCode: http.StatusBadRequest, resp: models.Book{}},
		{desc: "invalid param", id: "abc", req: &models.Book{BookID: 3, AuthorID: 1,
			Title: "300 Days", Publication: "Penguin", PublishedDate: "17/03/2016"}, resp: models.Book{}, statusCode: http.StatusBadRequest,
			err: errors.Error("invalid param")},
		{desc: "error in bind", id: "11", req: &[]models.Book{}, resp: models.Book{}, statusCode: http.StatusBadRequest,
			err: errors.Error("error in bind")},
	}

	ctr := gomock.NewController(t)
	mockBook := service.NewMockBook(ctr)
	delivery := New(mockBook)
	k := gofr.New()

	for i, v := range testcases {
		params := url.Values{}
		params.Add("bookId", v.id)

		body, err := json.Marshal(v.req)
		if err != nil {
			log.Printf("Not able to marshal : %v", err)
		}

		r := httptest.NewRequest(http.MethodGet, "/books?"+params.Encode(), bytes.NewReader(body))

		w := httptest.NewRecorder()

		r = mux.SetURLVars(r, map[string]string{"id": v.id})

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		id2, err2 := strconv.Atoi(v.id)
		if err != nil {
			log.Printf("error in converting string to int : %v", err2)
		}

		mockBook.EXPECT().Update(ctx, id2, v.req).Return(v.resp, v.err).AnyTimes()

		book, err3 := delivery.Update(ctx)

		if !reflect.DeepEqual(book, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, book, v.resp)
		}

		if err3 != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err2)
		}
	}
}

// TestDeleteBook function is to test delete method to remove any book
func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc        string
		id          string
		rowAffected int
		statusCode  int
		err         error
	}{
		{desc: "valid", id: "1", rowAffected: 1, statusCode: http.StatusOK, err: nil},
		{desc: "missing param", id: "", rowAffected: 0, err: errors.Error("missing param"), statusCode: http.StatusBadRequest},
		{desc: "invalid param", id: "abc", rowAffected: 0, err: errors.Error("invalid param"), statusCode: http.StatusBadRequest},
	}

	ctr := gomock.NewController(t)
	mockBook := service.NewMockBook(ctr)
	delivery := New(mockBook)
	k := gofr.New()

	for i, v := range testcases {
		params := url.Values{}
		params.Add("bookId", v.id)

		r := httptest.NewRequest(http.MethodGet, "/books?"+params.Encode(), nil)

		w := httptest.NewRecorder()

		r = mux.SetURLVars(r, map[string]string{"id": v.id})

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		id2, err2 := strconv.Atoi(v.id)
		if err2 != nil {
			log.Printf("error in converting string to int : %v", err2)
		}

		mockBook.EXPECT().Delete(ctx, id2).Return(v.rowAffected, v.err).AnyTimes()

		rowAffected, err := delivery.Delete(ctx)

		if !reflect.DeepEqual(rowAffected, v.rowAffected) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, rowAffected, v.rowAffected)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err)
		}
	}
}
