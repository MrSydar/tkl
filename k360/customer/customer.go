package customer

import "fmt"

type Customer struct {
	Id          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Nip         string `json:"VatRegNo,omitempty"`
	CountryCode string `json:"CountryCode,omitempty"`
	Regon       string `json:"RegNo,omitempty"`
	Street      string `json:"Address,omitempty"`
	PostalCode  string `json:"PostalCode,omitempty"`
	City        string `json:"City,omitempty"`
	County      string `json:"County,omitempty"`
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(e.Message)
}
