package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

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

	color.Green("dbType: %s", dbType)
	color.Green("auth tables created successfully")
	return nil
}
