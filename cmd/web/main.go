package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dh-n/bookings/pkg/config"
	"github.com/dh-n/bookings/pkg/handlers"
	"github.com/dh-n/bookings/pkg/render"
)

const portNumber = ":3000"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {
	// change this to true when in production
	app.InProduction = false

	// setting the session defaults for scs
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// create TemplateCache for the first time
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	// passing the app data to handler repo package
	repo := handlers.NewRepo(&app)

	// An alternate code for the above
	// repo := &handlers.Repository{
	// 	App: &app,
	// }

	handlers.SetNewRepoForHandlers(repo)

	// passing the app data to render package
	render.NewTemplate(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on port %s", portNumber)
	// err = http.ListenAndServe(portNumber, nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	//
	srv := &http.Server{
		Addr:    portNumber,
		Handler: router(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
