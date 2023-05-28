package render

import (
	"Bookings/internal/config"
	"Bookings/internal/models"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRender sets the config for the template package
func NewRender(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// populate the message in session for better user remind animation
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)                  //set the csrf token
	if app.Session.Exists(r.Context(), "user_id") { // set authenticated status
		td.IsAuthenticated = 1
	}
	return td
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache { //if using UseCache, just populate it, else create a new one
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl] //check if the template exists
	if !ok {
		// log.Fatal("Could not get template from template cache")
		return errors.New("can't get template from template cache")
	}

	buf := new(bytes.Buffer) //assign data to buffer

	td = AddDefaultData(td, r) //invoke AddDefaultData to add data to template data

	_ = t.Execute(buf, td) //execute the template using

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}
	//using Glob method to find the acutal file names includes .page.tmpl
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	//handle errors if any
	if err != nil {
		return myCache, err
	}
	// tranverse through pages
	for _, page := range pages {
		name := filepath.Base(page)                                     //using base method to cut the prefix
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) //create template using the file name
		// handle errors if any
		if err != nil {
			return myCache, err
		}
		//search match for layout template
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		// handle errors if any
		if err != nil {
			return myCache, err
		}
		//if there exist layout files
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates)) //adding the layout to the template
			//handle errors if any
			if err != nil {
				return myCache, err
			}
		}
		//populate the cache using key string map
		myCache[name] = ts
	}

	return myCache, nil
}
