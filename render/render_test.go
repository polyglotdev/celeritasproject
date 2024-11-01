package render

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CloudyKit/jet/v6"
)

func TestRender_Page(t *testing.T) {
	// Create temporary test directory and template
	tmpDir := t.TempDir()
	viewsDir := filepath.Join(tmpDir, "views")
	if err := os.MkdirAll(viewsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test templates for both Go and Jet
	goTmpl := `<h1>Hello {{.ServerName}}</h1>`
	if err := os.WriteFile(filepath.Join(viewsDir, "test.page.tmpl"), []byte(goTmpl), 0644); err != nil {
		t.Fatal(err)
	}

	jetTmpl := `<h1>Hello {{ .ServerName }}</h1>`
	if err := os.WriteFile(filepath.Join(viewsDir, "test.jet"), []byte(jetTmpl), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		renderer  string
		view      string
		data      any
		variables any
		wantErr   bool
		errMsg    string
		wantBody  string
		setupFunc func() *Render
	}{
		{
			name:     "valid go renderer",
			renderer: "go",
			view:     "test",
			wantErr:  false,
			wantBody: "<h1>Hello </h1>",
			setupFunc: func() *Render {
				return &Render{
					Renderer: "go",
					RootPath: tmpDir,
				}
			},
		},
		{
			name:     "valid jet renderer",
			renderer: "jet",
			view:     "test",
			data:     &TemplateData{ServerName: "TestServer"},
			wantErr:  false,
			wantBody: "<h1>Hello TestServer</h1>",
			setupFunc: func() *Render {
				views := jet.NewSet(
					jet.NewOSFileSystemLoader(filepath.Join(tmpDir, "views")),
					jet.InDevelopmentMode(),
				)
				return &Render{
					Renderer: "jet",
					RootPath: tmpDir,
					JetViews: views,
				}
			},
		},
		{
			name:     "invalid renderer",
			renderer: "invalid",
			view:     "test",
			wantErr:  true,
			errMsg:   "no renderer found",
			setupFunc: func() *Render {
				return &Render{
					Renderer: "invalid",
					RootPath: tmpDir,
				}
			},
		},
	}

	for _, test := range tests {
		testCase := test // Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			c := testCase.setupFunc()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			err := c.Page(w, r, testCase.view, testCase.variables, testCase.data)
			if (err != nil) != testCase.wantErr {
				t.Errorf("Page() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if testCase.wantErr && err != nil && err.Error() != testCase.errMsg {
				t.Errorf("Page() error message = %v, want %v", err, testCase.errMsg)
				return
			}
			if !testCase.wantErr && testCase.wantBody != "" && w.Body.String() != testCase.wantBody {
				t.Errorf("Page() body = %v, want %v", w.Body.String(), testCase.wantBody)
			}
		})
	}
}

func TestRender_GoPage(t *testing.T) {
	tmpDir := t.TempDir()
	viewsDir := filepath.Join(tmpDir, "views")
	if err := os.MkdirAll(viewsDir, 0755); err != nil {
		t.Fatal(err)
	}

	tmpl := `<h1>Hello {{.ServerName}}</h1>`
	if err := os.WriteFile(filepath.Join(viewsDir, "test.page.tmpl"), []byte(tmpl), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		view      string
		data      any
		wantErr   bool
		wantBody  string
		errMsg    string
		setupFunc func() *Render
	}{
		{
			name:     "valid template without data",
			view:     "test",
			wantErr:  false,
			wantBody: "<h1>Hello </h1>",
			setupFunc: func() *Render {
				return &Render{
					RootPath: tmpDir,
				}
			},
		},
		{
			name: "valid template with data",
			view: "test",
			data: &TemplateData{
				ServerName: "TestServer",
			},
			wantErr:  false,
			wantBody: "<h1>Hello TestServer</h1>",
			setupFunc: func() *Render {
				return &Render{
					RootPath: tmpDir,
				}
			},
		},
		{
			name:    "invalid template path",
			view:    "nonexistent",
			wantErr: true,
			errMsg:  "no such file or directory",
			setupFunc: func() *Render {
				return &Render{
					RootPath: tmpDir,
				}
			},
		},
		{
			name:    "invalid template data type",
			view:    "test",
			data:    "invalid",
			wantErr: true,
			errMsg:  "invalid template data type",
			setupFunc: func() *Render {
				return &Render{
					RootPath: tmpDir,
				}
			},
		},
	}

	for _, test := range tests {
		testCase := test // Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			c := testCase.setupFunc()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			err := c.GoPage(w, r, testCase.view, testCase.data)
			if (err != nil) != testCase.wantErr {
				t.Errorf("GoPage() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr && err != nil {
				if !strings.Contains(err.Error(), testCase.errMsg) {
					t.Errorf("GoPage() error = %v, want error containing %v", err, testCase.errMsg)
				}
				return
			}

			if !testCase.wantErr && w.Body.String() != testCase.wantBody {
				t.Errorf("GoPage() body = %v, want %v", w.Body.String(), testCase.wantBody)
			}
		})
	}
}

func TestRender_JetPage(t *testing.T) {
	// Create temporary test directory and template
	tmpDir := t.TempDir()
	viewsDir := filepath.Join(tmpDir, "views")
	if err := os.MkdirAll(viewsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Set up a basic Jet template
	tmpl := `<h1>Hello {{ .ServerName }}</h1>`
	if err := os.WriteFile(filepath.Join(viewsDir, "test.jet"), []byte(tmpl), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		view      string
		data      any
		variables any
		wantErr   bool
		errMsg    string
		wantBody  string
		setupFunc func() *Render
	}{
		{
			name:      "valid template with nil variables",
			view:      "test",
			data:      &TemplateData{ServerName: "TestServer"},
			variables: nil,
			wantErr:   false,
			wantBody:  "<h1>Hello TestServer</h1>",
			setupFunc: func() *Render {
				views := jet.NewSet(
					jet.NewOSFileSystemLoader(filepath.Join(tmpDir, "views")),
					jet.InDevelopmentMode(),
				)
				return &Render{
					JetViews: views,
					RootPath: tmpDir,
				}
			},
		},
		{
			name:      "invalid variables type",
			view:      "test",
			data:      nil,
			variables: "invalid type",
			wantErr:   true,
			errMsg:    "invalid variables type",
			setupFunc: func() *Render {
				views := jet.NewSet(
					jet.NewOSFileSystemLoader(filepath.Join(tmpDir, "views")),
					jet.InDevelopmentMode(),
				)
				return &Render{
					JetViews: views,
					RootPath: tmpDir,
				}
			},
		},
		{
			name:      "invalid data type",
			view:      "test",
			data:      "invalid",
			variables: make(jet.VarMap),
			wantErr:   true,
			errMsg:    "invalid template data type",
			setupFunc: func() *Render {
				views := jet.NewSet(
					jet.NewOSFileSystemLoader(filepath.Join(tmpDir, "views")),
					jet.InDevelopmentMode(),
				)
				return &Render{
					JetViews: views,
					RootPath: tmpDir,
				}
			},
		},
		{
			name:      "template not found",
			view:      "nonexistent",
			data:      nil,
			variables: make(jet.VarMap),
			wantErr:   true,
			errMsg:    "no template found",
			setupFunc: func() *Render {
				views := jet.NewSet(
					jet.NewOSFileSystemLoader(filepath.Join(tmpDir, "views")),
					jet.InDevelopmentMode(),
				)
				return &Render{
					JetViews: views,
					RootPath: tmpDir,
				}
			},
		},
	}

	for _, test := range tests {
		testCase := test // Capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			c := testCase.setupFunc()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			err := c.JetPage(w, r, testCase.view, testCase.variables, testCase.data)
			if (err != nil) != testCase.wantErr {
				t.Errorf("JetPage() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if testCase.wantErr && err != nil {
				if !strings.Contains(err.Error(), testCase.errMsg) {
					t.Errorf("JetPage() error = %v, want error containing %v", err, testCase.errMsg)
				}
			}
		})
	}
}
