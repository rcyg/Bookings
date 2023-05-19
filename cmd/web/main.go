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

var app config.AppConfig        //declare a config
var session *scs.SessionManager //declare a session
var infoLog *log.Logger         //information logger
var errorLog *log.Logger        //error logger

// main is the main application function
func main() {
	err := run() //using run function for testing
	if err != nil {
		log.Fatal(err) //if run occurs some err stop the application
	}

	fmt.Println("Starting at port", portNumber) //print the targeted port

	srv := &http.Server{ //assign the server including Address and Handler
		Addr:    portNumber,
		Handler: routes(&app),
	}

	if err = srv.ListenAndServe(); err != nil { //if failed stop the process
		log.Fatal(err)
	}
}

func run() error {
	// !IMPORTANT register the reservation data before using it
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false
	//infoLog
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//errorLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) //using Lshortfile to remind the error position
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()                            //Initialize
	session.Lifetime = 24 * time.Hour              //Give lifetime
	session.Cookie.Persist = true                  //The session will not be destroyed even after the user close the browser
	session.Cookie.SameSite = http.SameSiteLaxMode //
	session.Cookie.Secure = app.InProduction       //using app config to control this varaible

	app.Session = session // assign session to app configuratoin

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
