package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3 string) error {

	switch arg2 {
	case "migration":
		dbType := cel.DB.DataType
		if arg3 == "" {
			exitGracefully(errors.New("you must give the migration a name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFileFromTemplate("templates/migrations/migrations."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFileFromTemplate("templates/migrations/migrations."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the handler a name"))
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handlerName := strcase.ToCamel(strings.ReplaceAll(arg3, "handler", "Handler"))
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", handlerName)
		handler = strings.ReplaceAll(handler, "$FIRSTLETTER$", strings.ToLower(handlerName[:1]))
		err = copyDataToFile([]byte(handler), fileName)
		if err != nil {
			exitGracefully(err)
		}
	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the model a name"))
		}

		fileName := cel.RootPath + "/data/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		// Create pluralize client
		p := pluralize.NewClient()

		// Get singular form for the model name
		modelName := strcase.ToCamel(p.Singular(arg3))

		// Get plural form for table name
		tableName := strings.ToLower(p.Plural(arg3))

		model := string(data)
		model = strings.ReplaceAll(model, "$MODELNAME$", modelName)
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)
		model = strings.ReplaceAll(model, "$FIRSTLETTER$", strings.ToLower(modelName[:1]))

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}
