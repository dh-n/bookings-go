// Package render
package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dh-n/bookings/pkg/config"
	"github.com/dh-n/bookings/pkg/models"
)

// var templateCache = map[string]*template.Template{}
var app *config.AppConfig

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
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
		log.Fatal("No cache found")
	}

	// render template
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)

	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// if there are no errors write the buffer to respose writer
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	// parsedtemplate := template.Must(template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl"))
	// err := parsedtemplate.Execute(w, nil)
	// if err != nil {
	// fmt.Println("Error parsing the file ", err)
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the file named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// see if the current page.tmpl needs a layout file parsed
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
