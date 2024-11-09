package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	celeritas "github.com/polyglotdev/celeritasproject"
)

// Add type alias if needed
type Celeritas = celeritas.Celeritas

// Generator interface defines the contract for different types of generators
type Generator interface {
	Generate(name string) error
	ValidateName(name string) error
}

// BaseGenerator provides common functionality
type BaseGenerator struct {
	rootPath   string
	templateFS TemplateFS
	pluralize  *pluralize.Client
}

// TemplateFS interface to make testing easier
type TemplateFS interface {
	ReadFile(name string) ([]byte, error)
}

// FileSystem interface for testing file operations
type FileSystem interface {
	Exists(path string) bool
	WriteFile(path string, data []byte) error
	CopyFile(src, dst string) error
}

// MigrationGenerator handles migration generation
type MigrationGenerator struct {
	BaseGenerator
	dbType string
}

func (g *MigrationGenerator) Generate(name string) error {
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), name)

	upFile := fmt.Sprintf("%s/migrations/%s.%s.up.sql", g.rootPath, fileName, g.dbType)
	downFile := fmt.Sprintf("%s/migrations/%s.%s.down.sql", g.rootPath, fileName, g.dbType)

	if err := copyFileFromTemplate(
		fmt.Sprintf("templates/migrations/migrations.%s.up.sql", g.dbType),
		upFile,
	); err != nil {
		return err
	}

	return copyFileFromTemplate(
		fmt.Sprintf("templates/migrations/migrations.%s.down.sql", g.dbType),
		downFile,
	)
}

// AuthGenerator handles auth generation
type AuthGenerator struct {
	BaseGenerator
}

func (g *AuthGenerator) Generate(name string) error {
	return doAuth()
}

// HandlerGenerator handles handler generation
type HandlerGenerator struct {
	BaseGenerator
}

func (g *HandlerGenerator) Generate(name string) error {
	fileName := fmt.Sprintf("%s/handlers/%s.go", g.rootPath, strings.ToLower(name))
	if fileExists(fileName) {
		return fmt.Errorf("%s already exists", fileName)
	}

	data, err := g.templateFS.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		return err
	}

	handlerName := strcase.ToCamel(strings.ReplaceAll(name, "handler", "Handler"))
	content := string(data)
	content = strings.ReplaceAll(content, "$HANDLERNAME$", handlerName)
	content = strings.ReplaceAll(content, "$FIRSTLETTER$", strings.ToLower(handlerName[:1]))

	return copyDataToFile([]byte(content), fileName)
}

// ModelGenerator handles model generation
type ModelGenerator struct {
	BaseGenerator
}

func (g *ModelGenerator) Generate(name string) error {
	fileName := fmt.Sprintf("%s/data/%s.go", g.rootPath, strings.ToLower(name))
	if fileExists(fileName) {
		return fmt.Errorf("%s already exists", fileName)
	}

	data, err := g.templateFS.ReadFile("templates/data/model.go.txt")
	if err != nil {
		return err
	}

	modelName := strcase.ToCamel(g.pluralize.Singular(name))
	tableName := strings.ToLower(g.pluralize.Plural(name))

	content := string(data)
	content = strings.ReplaceAll(content, "$MODELNAME$", modelName)
	content = strings.ReplaceAll(content, "$TABLENAME$", tableName)
	content = strings.ReplaceAll(content, "$FIRSTLETTER$", strings.ToLower(modelName[:1]))

	return copyDataToFile([]byte(content), fileName)
}

// MiddlewareGenerator handles middleware generation
type MiddlewareGenerator struct {
	BaseGenerator
}

func (g *MiddlewareGenerator) Generate(name string) error {
	fileName := fmt.Sprintf("%s/middleware/%s.go", g.rootPath, strings.ToLower(name))
	if fileExists(fileName) {
		return fmt.Errorf("%s already exists", fileName)
	}

	data, err := g.templateFS.ReadFile("templates/middleware/middleware.go.txt")
	if err != nil {
		return err
	}

	middlewareName := strcase.ToCamel(name)
	content := string(data)
	content = strings.ReplaceAll(content, "$MIDDLEWARENAME$", middlewareName)

	return copyDataToFile([]byte(content), fileName)
}

// Factory to create appropriate generator
func newGenerator(genType string, cfg *Celeritas) (Generator, error) {
	base := BaseGenerator{
		rootPath:   cfg.RootPath,
		templateFS: templateFS,
		pluralize:  pluralize.NewClient(),
	}

	switch genType {
	case "migration":
		return &MigrationGenerator{BaseGenerator: base, dbType: cfg.DB.DataType}, nil
	case "auth":
		return &AuthGenerator{BaseGenerator: base}, nil
	case "handler":
		return &HandlerGenerator{BaseGenerator: base}, nil
	case "model":
		return &ModelGenerator{BaseGenerator: base}, nil
	case "middleware":
		return &MiddlewareGenerator{BaseGenerator: base}, nil
	default:
		return nil, fmt.Errorf("unknown generator type: %s", genType)
	}
}

// Main entry point
func doMake(arg2, arg3 string) error {
	gen, err := newGenerator(arg2, &cel)
	if err != nil {
		return err
	}

	if err := gen.ValidateName(arg3); err != nil {
		return err
	}

	return gen.Generate(arg3)
}

// Add validation methods for each generator
func (g *MigrationGenerator) ValidateName(name string) error {
	if name == "" {
		return errors.New("you must give the migration a name")
	}
	return nil
}

func (g *AuthGenerator) ValidateName(name string) error {
	return nil // Auth doesn't need a name
}

func (g *HandlerGenerator) ValidateName(name string) error {
	if name == "" {
		return errors.New("you must give the handler a name")
	}
	return nil
}

func (g *ModelGenerator) ValidateName(name string) error {
	if name == "" {
		return errors.New("you must give the model a name")
	}
	return nil
}

func (g *MiddlewareGenerator) ValidateName(name string) error {
	if name == "" {
		return errors.New("you must give the middleware a name")
	}
	return nil
}
