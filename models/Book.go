package models

type Book struct {
	BookID        int    `json:"bookID"`
	AuthorID      int    `json:"authID"`
	Auth          Author `json:"auth,omitempty"`
	Title         string `json:"title"`
	Publication   string `json:"publication"`
	PublishedDate string `json:"publishedDate"`
}
