// Package render
package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dh-n/bookings/internal/config"
	"github.com/dh-n/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// var templateCache = map[string]*template.Template{}
var (
	app             *config.AppConfig
	pathToTemplates = "./templates"
)

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	var err error
	if app.UseCache {
		// create template templateCache
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	// get the required template from tc
	t, ok := tc[tmpl]
	if !ok {
		log.Println("No cache found")
		return errors.New("can't get tempate from cache")
	}

	// render template
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)

	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}

	// if there are no errors write the buffer to respose writer
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	// parsedtemplate := template.Must(template.ParseFiles("%s/"+tmpl, "%s/base.layout.tmpl"))
	// err := parsedtemplate.Execute(w, nil)
	// if err != nil {
	// fmt.Println("Error parsing the file ", err)
	// }
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the file named *.page.tmpl from %s
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// see if there are any layout files
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// see if the current page.tmpl needs a layout file parsed
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
