package celeritas

const (
	Version = "1.0.0"
)

// Celeritas is the main application struct
type Celeritas struct {
	AppName string
	Debug   bool
	Version string
}

// New returns a new Celeritas application
func (c *Celeritas) New(rootPath string) error {
	c.AppName = "celeritas"
	c.Debug = true
	c.Version = Version

	return nil
}

// Init takes a initPaths and returns an error
func (c *Celeritas) Init(initPaths ...string) error {
	return nil
}
