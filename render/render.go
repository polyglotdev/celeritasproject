package render

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]any
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Flash           string
	Warning         string
	Error           string
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data any) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, data)
	default:
		return errors.New("no renderer found")
	}
}

func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data any) error {
	tmpl, err := template.ParseFiles(filepath.Join(c.RootPath, "views", fmt.Sprintf("%s.page.tmpl", view)))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, view string, data any) error {

	return nil
}
