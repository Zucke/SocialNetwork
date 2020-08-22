package errorstatus

import "errors"

//Some error status
var (
	ErrorUserExist   = errors.New("the user exist")
	ErrorNotFount    = errors.New("Not found")
	ErrorBadInfo     = errors.New("Bad Information")
	ErrorAccesDenied = errors.New("Access Denied")
)
