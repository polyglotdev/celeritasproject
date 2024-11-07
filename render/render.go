package render

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
)

// Render provides template rendering capabilities for web applications.
// It supports multiple rendering engines and manages template-related configuration.
//
// The zero value is not ready to use. All fields should be initialized before use.
type Render struct {
	Renderer   string   // Template engine to use ("go" or "jet")
	RootPath   string   // Base directory for template files
	Secure     bool     // Whether to use HTTPS
	Port       string   // Server port for URL generation
	ServerName string   // Server name for URL generation
	JetViews   *jet.Set // Jet template engine
	Session    *scs.SessionManager
}

// TemplateData holds all dynamic data needed for template rendering.
// It provides a structured way to pass data from handlers to templates.
//
// Fields are initialized to their zero values by default.
type TemplateData struct {
	IsAuthenticated bool               // Whether the current user is authenticated
	IntMap          map[string]int     // Integer data for template rendering
	StringMap       map[string]string  // String data for template rendering
	FloatMap        map[string]float32 // Float data for template rendering
	Data            map[string]any     // Generic data storage for template rendering
	CSRFToken       string             // Cross-Site Request Forgery protection token
	Port            string             // Server port for URL generation
	ServerName      string             // Server name for URL generation
	Secure          bool               // Whether to use HTTPS in generated URLs
	Flash           string             // Flash message to display once
	Warning         string             // Warning message to display
	Error           string             // Error message to display
}

// defaultData returns a pointer to a TemplateData struct with default values.
func (c *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Secure = c.Secure
	td.ServerName = c.ServerName
	td.Port = c.Port
	if c.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = true
	}
	return td
}

// Page renders a template using the configured rendering engine.
// It determines which rendering engine to use based on the Renderer field
// and delegates to the appropriate rendering method.
//
// Parameters:
//   - w: The HTTP response writer to render to
//   - r: The HTTP request being processed
//   - view: The name of the template to render
//   - variables: Additional variables for the template (currently unused)
//   - data: Optional TemplateData for the template
//
// Returns an error if rendering fails or if no valid renderer is configured.
func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data any) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variables, data)
	default:
		return errors.New("no renderer found")
	}
}

// GoPage renders a template using Go's standard template package.
// It loads the template from disk, processes it with the provided data,
// and writes the result to the HTTP response writer.
//
// The template file must have a .page.tmpl extension and be located in
// the views directory under RootPath.
//
// Parameters:
//   - w: The HTTP response writer to render to
//   - r: The HTTP request being processed
//   - view: The name of the template to render (without extension)
//   - data: Optional TemplateData for the template
//
// Returns an error if template loading or execution fails.
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data any) error {
	tmpl, err := template.ParseFiles(filepath.Join(c.RootPath, "views", fmt.Sprintf("%s.page.tmpl", view)))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		templateData, ok := data.(*TemplateData)
		if !ok {
			return fmt.Errorf("invalid template data type: expected *TemplateData, got %T", data)
		}
		td = templateData
	}

	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet template engine.
// This is a placeholder for future implementation.
//
// Parameters:
//   - w: The HTTP response writer to render to
//   - r: The HTTP request being processed
//   - view: The name of the template to render
//   - variables: Optional variables for the template
//   - data: Optional data for the template
//
// Currently returns nil as it is not implemented.
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, view string, variables, data any) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		var ok bool
		vars, ok = variables.(jet.VarMap)
		if !ok {
			return fmt.Errorf("invalid variables type: expected jet.VarMap, got %T", variables)
		}
	}

	td := &TemplateData{}
	if data != nil {
		templateData, ok := data.(*TemplateData)
		if !ok {
			return fmt.Errorf("invalid template data type: expected *TemplateData, got %T", data)
		}
		td = templateData
	}
	td = c.defaultData(td, r)

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		return fmt.Errorf("no template found: %w", err)
	}

	if err := t.Execute(w, vars, td); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
