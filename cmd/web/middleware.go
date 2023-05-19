package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf is the csrf protection middleware
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next) //create a csrf handler

	csrfHandler.SetBaseCookie(http.Cookie{ //
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode, //set same site attribute which prevents cross-site leakage
	})
	return csrfHandler
}

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
	//LoadAndSave method is actually a middleware provided by scs package
	//which enables automatic loads and saves for session every single request
}
