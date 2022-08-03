package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	datastoreauthor "mytest/datastore/author"
	datastorebook "mytest/datastore/book"
	deliveryauthor "mytest/delivery/author"
	deliverybook "mytest/delivery/book"
	serviceauthor "mytest/service/author"
	servicebook "mytest/service/book"
)

func main() {
	authorDatastore := datastoreauthor.New()
	authorService := serviceauthor.New(authorDatastore)
	authorHandler := deliveryauthor.New(authorService)

	bookDatastore := datastorebook.New()
	bookService := servicebook.New(bookDatastore, authorDatastore)
	bookHandler := deliverybook.New(bookService)

	r := gofr.New()

	// Author endpoint
	r.POST("/author", authorHandler.Create)
	r.PUT("/author/{id}", authorHandler.Update)
	r.DELETE("/author/{id}", authorHandler.Delete)

	// Book endpoints
	r.POST("/book", bookHandler.Create)
	r.GET("/books", bookHandler.GetAll)
	r.GET("/book/{id}", bookHandler.GetByID)
	r.PUT("/book/{id}", bookHandler.Update)
	r.DELETE("/book/{id}", bookHandler.Delete)

	r.Start()

}
