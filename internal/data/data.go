package data

import (
	"errors"
)

//DBNAME is the name of the database we are working
const DBNAME = "socialnetwork"

//all names of collections
const (
	UserColletion         = "user"
	PublicationsColletion = "publications"
)

//Some error status
var (
	ErrorUserExist   = errors.New("the user exist")
	ErrorNotFount    = errors.New("Not found")
	ErrorBadInfo     = errors.New("Bad Information")
	ErrorAccesDenied = errors.New("Access Denied")
)

//FieldValidation is the data that has a validations his field
type FieldValidation interface {
	IsValidFields() bool
}
