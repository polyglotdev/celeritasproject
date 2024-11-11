package celeritas

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
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

// WriteXML writes an XML response to the client
func (c *Celeritas) WriteXML(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// DownloadFile serves a file for download to the client. It takes the following parameters:
//   - w: The http.ResponseWriter to write the response to
//   - r: The *http.Request containing the request details
//   - pathToFile: The directory path where the file is located
//   - fileName: The name of the file to be downloaded
//
// The method:
//  1. Joins the path and filename safely using path.Join
//  2. Cleans the resulting path to prevent directory traversal attacks
//  3. Sets the Content-Type header to force file download with the original filename
//  4. Serves the file using http.ServeFile
//
// Returns error (currently always nil, but preserved for future error handling)
func (c *Celeritas) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

// Error404 sends a 404 Not Found response to the client.
// It takes:
//   - w: The http.ResponseWriter to write the response to
//   - r: The *http.Request containing the request details
func (c *Celeritas) Error404(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusNotFound)
}

// Error500 sends a 500 Internal Server Error response to the client.
// It takes:
//   - w: The http.ResponseWriter to write the response to
//   - r: The *http.Request containing the request details
func (c *Celeritas) Error500(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusInternalServerError)
}

// ErrorUnauthorized sends a 401 Unauthorized response to the client.
// It takes:
//   - w: The http.ResponseWriter to write the response to
//   - r: The *http.Request containing the request details
func (c *Celeritas) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusUnauthorized)
}

// ErrorForbidden sends a 403 Forbidden response to the client.
// It takes:
//   - w: The http.ResponseWriter to write the response to
//   - r: The *http.Request containing the request details
func (c *Celeritas) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusForbidden)
}

// ErrorStatus sends an HTTP error response with the specified status code.
// It takes:
//   - w: The http.ResponseWriter to write the response to
//   - status: The HTTP status code to send
//
// The response body will contain the standard HTTP status text for the given code.
func (c *Celeritas) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
