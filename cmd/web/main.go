package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dh-n/bookings/internal/config"
	"github.com/dh-n/bookings/internal/handlers"
	"github.com/dh-n/bookings/internal/helpers"
	"github.com/dh-n/bookings/internal/render"
)

const portNumber = ":3000"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s", portNumber)
	// err = http.ListenAndServe(portNumber, nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	//
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
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
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	// initialize logger
	infoLog := log.New(os.Stdout, "INFO", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)

	app.ErrorLog = errorLog
	app.InfoLog = infoLog

	// passing the app data to handler repo package
	repo := handlers.NewRepo(&app)

	// An alternate code for the above
	// repo := &Repository{
	// 	App: &app,
	// }

	handlers.SetNewRepoForHandlers(repo)

	// passing the app data to render package
	render.NewTemplate(&app)

	// passing the app data to helpers package
	helpers.NewHelper(&app)
	// http.HandleFunc("/", Repo.Home)
	// http.HandleFunc("/about", Repo.About)
	return nil
}
