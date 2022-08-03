package author

import (
	"Three-Layer-Architecture/models"
	"Three-Layer-Architecture/service"
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

// TestPostAuthor function is to test post method
func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		req        interface{}
		resp       models.Author
		StatusCode int
		err        error
	}{
		{desc: "valid", req: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, resp: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, StatusCode: http.StatusCreated},
		{desc: "error in bind", req: "Sujeet", StatusCode: http.StatusBadRequest},
		{desc: "errors from svc", req: models.Author{AuthID: -21, FirstName: "Sagar", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, StatusCode: http.StatusBadRequest, err: errors.New("invalid id")},
	}

	ctr := gomock.NewController(t)
	mockAuthor := service.NewMockAuthor(ctr)
	delivery := New(mockAuthor)
	k := gofr.New()

	for i, v := range testcases {
		body, err := json.Marshal(v.req)
		if err != nil {
			log.Printf("Not able to marshal : %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/author", bytes.NewReader(body))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		mockAuthor.EXPECT().Post(ctx, v.req).Return(v.resp, v.err).AnyTimes()

		author, err2 := delivery.Create(ctx)

		if !reflect.DeepEqual(author, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, author, v.resp)
		}

		if err2 != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err2)
		}

	}
}

// TestUpdateAuthor function is to test put method
func TestUpdateAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		id         string
		req        interface{}
		resp       models.Author
		StatusCode int
		err        error
	}{
		{desc: "valid case", id: "1", req: models.Author{AuthID: 1, FirstName: "Rajan", LastName: "Sharma",
			Dob: "26/04/2001", PenName: "Sharma"}, StatusCode: http.StatusOK},
		{desc: "missing param", id: "", StatusCode: http.StatusBadRequest},
		{desc: "error in bind", id: "1", req: "something", StatusCode: http.StatusBadRequest},
		{desc: "errors from svc", id: "11", req: models.Author{AuthID: -21, FirstName: "Sagar", LastName: "Bhagat", Dob: "06/04/2001",
			PenName: "Chetan"}, StatusCode: http.StatusBadRequest, err: errors.New("invalid id")},
		{desc: "error in strconv", id: "abc", req: models.Author{AuthID: 14, FirstName: "Rajan", LastName: "Sharma",
			Dob: "26/04/2001", PenName: "Sharma"}, StatusCode: http.StatusBadRequest},
	}

	ctr := gomock.NewController(t)
	mockAuthor := service.NewMockAuthor(ctr)
	delivery := New(mockAuthor)
	k := gofr.New()

	for i, v := range testcases {
		params := url.Values{}
		params.Add("bookId", v.id)

		body, err := json.Marshal(v.req)
		if err != nil {
			log.Printf("Not able to marshal : %v", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/author?"+params.Encode(), bytes.NewReader(body))
		w := httptest.NewRecorder()

		r = mux.SetURLVars(r, map[string]string{"id": v.id})

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		id, err2 := strconv.Atoi(v.id)
		if err2 != nil {
			log.Printf("error in string conversion : %v\n", err2)
		}

		mockAuthor.EXPECT().Update(ctx, id, v.req).Return(v.resp, v.err).AnyTimes()

		// Mocking Update
		author, err3 := delivery.Update(ctx)

		if !reflect.DeepEqual(author, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, author, v.resp)
		}

		if err3 != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err3)
		}
	}
}

// TestDeleteAuthor function is to test delete method
func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc        string
		id          string
		statusCode  int
		rowAffected int
		err         error
	}{
		{desc: "valid case", id: "1", statusCode: http.StatusNoContent, rowAffected: 1},
		{desc: "error from svc", id: "-11", statusCode: http.StatusBadRequest},
		{desc: "missing params", id: "", statusCode: http.StatusBadRequest},
		{desc: "error in strconv", id: "abc", statusCode: http.StatusBadRequest},
	}

	ctr := gomock.NewController(t)
	mockAuthor := service.NewMockAuthor(ctr)
	delivery := New(mockAuthor)
	k := gofr.New()

	for i, v := range testcases {
		params := url.Values{}
		params.Add("bookId", v.id)

		r := httptest.NewRequest(http.MethodPost, "/author/{id}"+v.id, nil)
		w := httptest.NewRecorder()

		r = mux.SetURLVars(r, map[string]string{"id": v.id})

		req := request.NewHTTPRequest(r)
		resp := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(resp, req, k)

		id, err2 := strconv.Atoi(v.id)
		if err2 != nil {
			log.Printf("error in converting string : %v", err2)
		}

		mockAuthor.EXPECT().Delete(ctx, id).Return(v.rowAffected, v.err).AnyTimes()

		rowAffected, err := delivery.Delete(ctx)

		if rowAffected != v.rowAffected {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, rowAffected, v.rowAffected)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\n", v.desc, i+1, err)
		}

	}
}
