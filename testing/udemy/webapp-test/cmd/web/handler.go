package main

import (
	"html/template"
	"net/http"
	"path"
)

var pathToTemplate = "./template/"

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {

	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplate, t))

	data.IP = app.ipFromConext(r.Context())
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	err = parsedTemplate.Execute(w, data)

	if err != nil {
		return err
	}
	return nil
}
