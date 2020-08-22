package data

//DBNAME is the name of the database we are working
const DBNAME = "socialnetwork"

//all names of collections
const (
	UserColletion         = "user"
	PublicationsColletion = "publications"
)

//FieldValidation is the data that has a validations his field
type FieldValidation interface {
	IsValidFields() bool
}
