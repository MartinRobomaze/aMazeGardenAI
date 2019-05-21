package main

import (
	"html/template"
	"net/http"
)

type FileTemplateLoader struct {
	Path      string
	DbHandler DatabaseHandler
}

func (fileTemplateLoader FileTemplateLoader) LoadFileTemplate(writer http.ResponseWriter, request *http.Request) {
	plants, err := fileTemplateLoader.DbHandler.GetAllPlantsNames()

	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles(fileTemplateLoader.Path)

	if err != nil {
		panic(err)
	}

	options := ""

	for i := 0; i < len(plants); i++ {
		options += "<option value=" + plants[i] + ">" + plants[i] + "</option>"
	}

	err = t.Execute(writer, template.HTML(options))

	if err != nil {
		panic(err)
	}
}
