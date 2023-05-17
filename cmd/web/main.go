package main

import (
	"Bookings/internal/config"
	"Bookings/internal/handlers"
	"Bookings/internal/render"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig //declare a config
var session *scs.SessionManager

// main is the main application function
func main() {
	//change this to true when in production
	app.InProduction = true

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache() //create template cache
	if err != nil {                         //if can't create template cache
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc //set the template cache
	app.UserCache = false  //set UserCache to false

	repo := handlers.NewRepo(&app) //deliver the configuration to handler
	render.NewTemplates(&app)      //deliver the configuration to render

	handlers.NewHandlers(repo) //create a new handler

	fmt.Println("Starting at port", portNumber) //port

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
