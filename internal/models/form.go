package models

import "net/url"

type FormErr struct {
	FormValue url.Values
	Err       map[string][]string
}

type Form struct {
	Form FormErr
}
