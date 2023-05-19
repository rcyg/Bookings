package main

import (
	"Bookings/internal/config"
	"Bookings/internal/handlers"
	"Bookings/internal/helpers"
	"Bookings/internal/models"
	"Bookings/internal/render"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig //declare a config
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err) //if run occurs some err stop the application
	}

	fmt.Println("Starting at port", portNumber) //port

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//register the reservation data
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false
	//infoLog
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//errorLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
		return err
	}

	app.TemplateCache = tc //set the template cache
	app.UseCache = false   //set UserCache to false

	repo := handlers.NewRepo(&app) //deliver the configuration to handler
	render.NewTemplates(&app)      //deliver the configuration to render

	handlers.NewHandlers(repo) //create a new handler
	helpers.NewHelpers(&app)   //deliver app configuration to helpers

	return nil
}
