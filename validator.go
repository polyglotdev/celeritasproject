package celeritas

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// Validation struct holds the data and errors
type Validation struct {
	Data   url.Values
	Errors map[string]string
}

// Validator creates a new Validation instance
func (c *Celeritas) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

// Valid returns true if there are no errors
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a given field and returns true if added
func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Has checks if a given field has an error
func (v *Validation) Has(field string, r *http.Request) bool {
	return r.Form.Get(field) != ""
}

// Required checks if the required fields are not empty
func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "this field is required")
		}
	}
}

// Check checks if the condition is true and adds an error if not
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

// IsInt checks if the value is an integer
func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "This field must be an integer")
	}
}

// IsFloat checks if the value is a float
func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "This field must be a floating point number")
	}
}

// IsDateISO checks if the value is a valid ISO date
func (v *Validation) IsDateISO(field, value string) {
	validFormat := "2006-03-04"
	_, err := time.Parse(validFormat, value)
	if err != nil {
		v.AddError(field, "this field must be a valid date in format YYYY-MM-DD")
	}
}

// NoSpaces removes all whitespace from the value
func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "this field cannot contain spaces")
	}
}
