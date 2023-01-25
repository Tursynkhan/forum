package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors Err
}

func New(data url.Values) *Form {
	return &Form{
		data,
		Err(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) UserNameValid(field string) {
	value := f.Get(field)

	if utf8.RuneCountInString(value) < 8 {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d character)", 8))
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too lonf (maximum is %d character)", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Minlength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d character)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) ErrorField(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) ErrorMatchesPattern(field string) {
	str := `Minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character:`
	f.Errors.Add(field, str)
}

func (f *Form) IsExist(field string, str string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	f.Errors.Add(field, str)
}
