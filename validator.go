package celeritas

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (v *Validation) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

func (v *Validation) Value() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a given field and returns true if added
func (v *Validation) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// Has checks if a given field has an error
func (v *Validation) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	_, exists := v.Errors[field]
	return exists
}

func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "this field is required")
		}
	}
}

func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// IsEmail takes field and value and returns true if the value is a valid email
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "invalid email address")
	}
}
