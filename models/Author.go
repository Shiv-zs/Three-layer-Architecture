package models

type Author struct {
	AuthID    int    `json:"authID,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Dob       string `json:"dob,omitempty"`
	PenName   string `json:"penName,omitempty"`
}
