package main

import (
	"github.com/FilipeParreiras/Bookings/internal/helpers"
	"github.com/justinas/nosurf"
	"net/http"
)

// SessionLoad uses a function called LoadAndSave which provides middleware which
// automatically loads and saves session data for the current request, and
// communicates the session token to and from the client in a cookie.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// NoSurf deals with CSRF
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// Auth checks if the user is already authenticated
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !helpers.IsAuthenticated(request) {
			session.Put(request.Context(), "error", "Log in first.")
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
		}
		next.ServeHTTP(writer, request)
	})
}
