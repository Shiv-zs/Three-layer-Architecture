package author

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"log"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"Three-Layer-Architecture/datastore"
	"Three-Layer-Architecture/models"
)

// TestAuthor_Post function is to test post author details for valid conditions
func TestAuthor_Post(t *testing.T) {
	testcases := []struct {
		desc     string
		req      models.Author
		response models.Author
		err      error
	}{
		{desc: "valid details", req: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			response: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"}},
		{desc: "missing first name", req: models.Author{AuthID: 1, LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			err: errors.Error("missing fields")},
		{desc: "invalid id", req: models.Author{AuthID: -11, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Chetan"},
			err: errors.Error("invalid id")},
		{desc: "missing last name", req: models.Author{AuthID: 1, FirstName: "Chetan", Dob: "06/04/2001", PenName: "Chetan"},
			err: errors.Error("missing fields")},
		{desc: "missing dob", req: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", PenName: "Chetan"},
			err: errors.Error("missing fields")},
		{desc: "missing penName", req: models.Author{AuthID: 1, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001"},
			err: errors.Error("missing fields")},
		{desc: "error in post", req: models.Author{AuthID: 12, FirstName: "Chetan", LastName: "Bhagat", Dob: "06/04/2001", PenName: "Bhagat"}, response: models.Author{},
			err: errors.Error("error in post")},
	}

	ctr := gomock.NewController(t)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockAuthor.EXPECT().Post(c, v.req).Return(v.response, v.err).AnyTimes()

		resp, err := service.Post(c, v.req)

		if !reflect.DeepEqual(resp, v.response) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.response)
		}

		if !reflect.DeepEqual(err, v.err) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// TestAuthor_PutValidID function is to test update valid author details
func TestAuthor_PutValidID(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		req         models.Author
		resp        models.Author
		checkAuthor bool
		err         error
	}{
		{desc: "valid", id: 1, req: models.Author{AuthID: 1, FirstName: "Rajan", LastName: "Sharma",
			Dob: "26/04/2001", PenName: "Rajan"}, resp: models.Author{AuthID: 1, FirstName: "Rajan", LastName: "Sharma",
			Dob: "26/04/2001", PenName: "Rajan"}, checkAuthor: false},
		{desc: "missing fields", id: 1, req: models.Author{AuthID: 1, FirstName: "Rajan", LastName: "Sharma",
			Dob: "26/04/2001", PenName: ""}, checkAuthor: false},
		{desc: "invalid id", id: -11, err: errors.Error("invalid id"), checkAuthor: false},
		{desc: "error in Update", id: 11, req: models.Author{AuthID: 14, FirstName: "Rajan", LastName: "Sharma",
			Dob: "26/04/2001", PenName: "Rajan"}, err: errors.Error("error in update"), checkAuthor: false},
	}

	ctr := gomock.NewController(t)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockAuthor.EXPECT().Update(c, v.id, v.req).Return(v.resp, v.err).AnyTimes()
		mockAuthor.EXPECT().IsAuthorIDPresent(c, v.id).Return(v.checkAuthor).AnyTimes()

		resp, err := service.Update(c, v.id, v.req)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// TestAuthor_PutInvalidID function is to test update valid author details
func TestAuthor_PutInvalidID(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		req         models.Author
		resp        models.Author
		checkAuthor bool
		err         error
	}{
		{desc: "isAuthorPresent", req: models.Author{AuthID: 12, FirstName: "shiv", LastName: "chandra",
			Dob: "01/01/2001", PenName: "shiv"}, id: 1, checkAuthor: true},
	}

	ctr := gomock.NewController(t)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context

		mockAuthor.EXPECT().Update(c, v.id, v.req).Return(v.resp, v.err).AnyTimes()
		mockAuthor.EXPECT().IsAuthorIDPresent(c, v.id).Return(v.checkAuthor).AnyTimes()

		resp, err := service.Update(c, v.id, v.req)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.resp)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d] Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// TestAuthor_DeleteValidID function is to test for remove author
func TestAuthor_DeleteValidID(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		rowAffected int
		checkAuthor bool
		err         error
	}{
		{desc: "valid", id: 1, rowAffected: 1},
		{desc: "invalid id", id: -11, err: errors.Error("invalid id")},
		{desc: "error in delete", id: 12, err: errors.Error("error in delete")},
	}

	ctr := gomock.NewController(t)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context
		mockAuthor.EXPECT().Delete(c, v.id).Return(v.rowAffected, v.err).AnyTimes()
		mockAuthor.EXPECT().IsAuthorIDPresent(c, v.id).Return(v.checkAuthor).AnyTimes()

		resp, err := service.Delete(c, v.id)

		if !reflect.DeepEqual(resp, v.rowAffected) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.rowAffected)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}

// TestAuthor_DeleteValidID function is to test for remove author
func TestAuthor_DeleteInvalidID(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		rowAffected int
		checkAuthor bool
		err         error
	}{
		{desc: "isAuthorPresent", id: 1, checkAuthor: true},
	}

	ctr := gomock.NewController(t)
	mockAuthor := datastore.NewMockAuthor(ctr)
	service := New(mockAuthor)

	for i, v := range testcases {
		var c *gofr.Context
		mockAuthor.EXPECT().Delete(c, v.id).Return(v.rowAffected, v.err).AnyTimes()
		mockAuthor.EXPECT().IsAuthorIDPresent(c, v.id).Return(v.checkAuthor).AnyTimes()

		resp, err := service.Delete(c, v.id)

		if !reflect.DeepEqual(resp, v.rowAffected) {
			t.Errorf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, resp, v.rowAffected)
		}

		if err != nil {
			log.Printf("desc : %v ,[TEST%d]Failed. Got %v\tExpected %v\n", v.desc, i+1, err, v.err)
		}
	}
}
