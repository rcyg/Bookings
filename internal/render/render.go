package render

import (
	"Bookings/internal/config"
	"Bookings/internal/models"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

// var NewTemplates(a *config.AppConfig)
var app *config.AppConfig

// NewTemplate receives appconfig from main.go
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData receives templateData from handlers and return a pointer to that data
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplates renders templates using html/templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UserCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
	}
}

// CreateTemplateCache creates the template cache
// which create html template for page.tmpl and layout.tmpl
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// get all of the files named *.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl") //using filepath.Glob() method to find all file names that matches page.tmpl
	if err != nil {
		fmt.Println("Failed to find template for page")
		return myCache, err
	}
	// range through all files ending with *.page.tmpl
	for _, page := range pages { //page is the full path name
		name := filepath.Base(page) //using Base() method to get the file name
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("Failed to create template for page")
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl") //find the matches of layout.tmpl
		if err != nil {
			fmt.Println("Failed to find template for layout.tmpl")
			return myCache, err
		}
		if len(matches) > 0 { // if there actually exist matches for layout.tmpl
			//substitute it
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				fmt.Println("Failed to create template for layout")
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
