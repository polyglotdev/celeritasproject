package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// doAuth orchestrates the setup of authentication infrastructure in a new project.
// It generates timestamped migrations to ensure consistent database state across
// environments and deployment, choosing this approach over runtime table creation
// to maintain better version control and deployment predictability.
//
// The function copies predefined templates rather than generating code dynamically
// to ensure consistent, tested authentication patterns and reduce the likelihood
// of security anti-patterns being introduced during authentication setup.
//
// The choice of separate up/down migrations supports both forward development
// and rollback scenarios, while the use of cascading drops in down migrations
// ensures clean removal of the auth infrastructure when needed.
//
// Template files for models and middleware are copied to their respective
// directories to maintain separation of concerns and follow standard Go project
// layout conventions.
//
// The function provides immediate feedback through colored console output to guide
// developers through necessary post-setup steps, reducing the likelihood of
// incomplete authentication implementation.
func doAuth() error {
	dbType := cel.DB.DataType
	filename := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := cel.RootPath + "/migrations/" + filename + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + filename + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte(`
drop table if exists remember_tokens cascade;
drop table if exists tokens cascade;
drop table if exists users cascade;
`), downFile)
	if err != nil {
		exitGracefully(err)
	}
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/user.go.txt", cel.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", cel.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate(
		"templates/middleware/auth.go.txt",
		cel.RootPath+"/middleware/auth.go",
	)
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate(
		"templates/middleware/auth-token.go.txt",
		cel.RootPath+"/middleware/auth-token.go",
	)
	if err != nil {
		exitGracefully(err)
	}

	color.Green("dbType: %s", dbType)
	color.Green("auth middleware created successfully")
	color.Yellow("Don't forget to:")
	color.Yellow("1. Add User and Token models to data/models")
	color.Yellow("2. Add auth middleware to routes.go")
	return nil
}
