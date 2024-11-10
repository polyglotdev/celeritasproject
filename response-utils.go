package celeritas

import (
	"encoding/json"
	"net/http"
)

// WriteJSON provides a consistent way to send JSON responses across the application, handling
// common response needs like custom headers and content-type setting. It abstracts away the
// repetitive tasks of JSON marshaling and header management that would otherwise be scattered
// throughout handler code.
//
// The method uses MarshalIndent instead of Marshal to generate formatted JSON that's easier
// to read during development and debugging, accepting the minor performance trade-off for
// improved developer experience.
//
// Headers are made optional through variadic parameters to support both simple responses
// and cases requiring custom header inclusion without complicating the common case.
//
// Example usage with custom headers:
//
//	headers := http.Header{
//	    "X-Custom-Header": []string{"value"},
//	}
//	err := c.WriteJSON(w, http.StatusOK, data, headers)
//
// Example simple usage:
//
//	err := c.WriteJSON(w, http.StatusOK, data)
func (c *Celeritas) WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers ...http.Header,
) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}
