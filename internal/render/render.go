package render

import (
	"bytes"
	"github.com/FilipeParreiras/Bookings/internal/config"
	"github.com/FilipeParreiras/Bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

/*
// RenderTemplate with tmpl(template) renders html pages

	func RenderTemplate(w http.ResponseWriter, tmpl string) {
		// ParseFiles takes one or more arguments -> variadic functions
		parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")

		err := parsedTemplate.Execute(w, nil)
		if err != nil {
			fmt.Println("Error parsing template:", err)
		}
	}

// template cache
var tc = make(map[string]*template.Template)

// Simpler version

	func RenderTemplateTestS(w http.ResponseWriter, t string) {
		var tmpl *template.Template
		var err error

		// check to see if we already have the template in our cache
		_, inMap := tc[t]
		if !inMap {
			// need to create the template
			err = createTemplateCache(t)
			if err != nil {
				log.Println(err)
			}
		} else {
			// we have the template in the cache
			log.Println("Using cached template")
		}

		tmpl := tc[t]
		err = tmpl.Execute(w, nil)
	}

// createTemplateCache used on RenderTemplate(simpler version)

	func createTemplateCache(t string) error {
		templates := []string{
			fmt.Sprintf("./templates/%s", t),
			"./templates/base.layout.tmpl",
		}

		// parse the template
		tmpl, err := template.ParseFiles(templates...)
		if err != nil {
			return err
		}

		// add template to cache
		tc[t] = tmpl

		return nil
	}
*/

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the template cache from the app config
	tc = app.TemplateCache

	// create a template cache
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(err)
	}
	// buffer
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.html
	for _, page := range pages {
		// name of the file without full path
		name := filepath.Base(page)
		// templateSet
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		//adds the final resulting template to myCache
		myCache[name] = ts
	}

	return myCache, nil
}
