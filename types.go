package celeritas

// initPaths represents the paths and folder names to initialize the application
type initPaths struct {
	rootPath    string
	folderNames []string
}

// cookieConfig represents the configuration for the session cookie
type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}
