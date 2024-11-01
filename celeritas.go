package celeritas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/polyglotdev/celeritasproject/render"
)

const (
	Version = "1.0.0"
)

// Celeritas represents the core application framework, providing essential services
// for web applications including routing, logging, and configuration management.
// It coordinates all the main components needed to run a web server.
//
// The zero value is not usable. Use New to create a new Celeritas instance.
//
// A Celeritas instance manages:
//   - Application identification and versioning
//   - Debug mode configuration
//   - Structured logging (info and error levels)
//   - File system organization via RootPath
//   - HTTP routing with Chi router
//   - Server configuration (ports, rendering options)
//   - Jet template engine for rendering views
//   - Go template engine for rendering views
//
// Celeritas is safe for use by a single goroutine at a time.
type Celeritas struct {
	AppName  string         // Application name used in logging and identification
	Debug    bool           // Debug mode flag for detailed logging and error handling
	Version  string         // Application version for deployment tracking
	ErrorLog *log.Logger    // Structured error logging
	InfoLog  *log.Logger    // Structured information logging
	RootPath string         // Base directory for application files and folders
	Routes   *chi.Mux       // HTTP router for handling web requests
	config   config         // Internal server configuration settings
	Render   *render.Render // Rendering engine
	JetViews *jet.Set       // Jet template engine
}

type config struct {
	port     string
	renderer string
}

// New returns a new Celeritas application
func (c *Celeritas) New(rootPath string) error {
	c.AppName = "celeritas"
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers",
			"migrations",
			"views",
			"data",
			"public",
			"tmp",
			"logs",
			"middleware",
		},
	}

	if err := c.Init(pathConfig); err != nil {
		return err
	}

	if err := c.checkDotEnv(rootPath); err != nil {
		return err
	}

	// read .env file
	if err := godotenv.Load(fmt.Sprintf("%s/.env", rootPath)); err != nil {
		return err
	}

	// start loggers
	infoLog, errorLog := c.StartLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errorLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = Version
	c.RootPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)

	c.JetViews = views

	c.createRenderer()

	return nil
}

// Init takes a initPaths and returns an error
func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		if err := c.CreateDirIfNotExist(filepath.Join(root, path)); err != nil {
			return err
		}
	}
	return nil
}

// checkDotEnv creates a .env file in the specified root directory if one doesn't exist.
// It takes a path string and returns an error if the file creation fails.
// If the .env file already exists, it does nothing and returns nil.
func (c *Celeritas) checkDotEnv(p string) error {
	if err := c.CreateFileIfNotExists(fmt.Sprintf("%s/.env", p)); err != nil {
		return err
	}
	return nil
}

// createServer creates and configures an HTTP server with the application settings.
// It returns the configured server and an error if configuration fails.
//
// The server is configured with:
//   - Port from application config
//   - Application error logger
//   - Chi router as handler
//   - Idle timeout of 30 seconds
//   - Read timeout of 30 seconds
//   - Write timeout of 600 seconds
//
// If the port is not configured (empty), createServer returns an error.
func (c *Celeritas) createServer() (*http.Server, error) {
	if c.config.port == "" {
		return nil, fmt.Errorf("port cannot be empty")
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", c.config.port),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	c.InfoLog.Printf("Starting %s on port %s", c.AppName, c.config.port)
	return srv, nil
}

// ListenAndServe starts and runs the HTTP server until it encounters an error.
// It configures the server with:
//
//   - Idle timeout of 30 seconds
//   - Read timeout of 30 seconds
//   - Write timeout of 600 seconds
//
// The server listens on the port specified in the application config.
// If the server encounters an error during startup or operation, it logs
// the error and terminates the application.
func (c *Celeritas) ListenAndServe() error {
	srv, err := c.createServer()
	if err != nil {
		return err
	}
	return srv.ListenAndServe()
}

// StartLoggers initializes the application's logging system with two loggers:
// an InfoLog for general information (in green) and an ErrorLog for error messages (in red).
// Both loggers write to standard output with different prefixes, formats, and colors.
// It returns two loggers for info and error logging respectively.
func (c *Celeritas) StartLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	// Create colored prefixes using fatih/color
	infoPrefix := color.New(color.FgGreen).Sprint("INFO\t")
	errorPrefix := color.New(color.FgRed).Sprint("ERROR\t")

	infoLog = log.New(os.Stdout, infoPrefix, log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, errorPrefix, log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (c *Celeritas) createRenderer() {
	myRenderer := render.Render{
		Renderer: c.config.renderer,
		RootPath: c.RootPath,
		Port:     c.config.port,
		JetViews: c.JetViews,
	}
	c.Render = &myRenderer
}
