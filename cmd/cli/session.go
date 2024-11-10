package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func doSessionTable() error {
	dbType := cel.DB.DataType

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_session_table", time.Now().UnixMicro())

	upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("/templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("/templates/migrations/session.down.sql.txt", downFile)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("Don't forget to run:")
	color.Yellow("migrate up")
	color.Green("Session table migration created successfully")

	return nil
}
