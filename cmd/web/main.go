package main

import (
	"encoding/gob"
	"fmt"
	"github.com/FilipeParreiras/Bookings/internal/config"
	"github.com/FilipeParreiras/Bookings/internal/handlers"
	"github.com/FilipeParreiras/Bookings/internal/models"
	"github.com/FilipeParreiras/Bookings/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// What am I going to put in the session
	gob.Register(models.Reservation{})

	// Change this to true when is production
	app.InProduction = false

	// Sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Printf("Starting application on port %s", portNumber))

	// Routing
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
