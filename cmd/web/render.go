package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// templateData holds dynamic data that will be passed to templates.
// The map allows storing any kind of value (string, int, struct, etc.)
// which makes templates flexible.
type templateData struct {
	Data map[string]any
}

// render is responsible for:
// 1. Finding the requested template
// 2. Loading it from cache or disk
// 3. Executing it and sending HTML to the browser
func (app *application) render(w http.ResponseWriter, t string, td *templateData) {
	var tmpl *template.Template

	// If template caching is enabled, try to fetch the template
	// from the in-memory map instead of reading from disk.
	// This improves performance in production.
	if app.config.useCache {
		// Check if the template exists in the map
		if templateFromMap, ok := app.templateMap[t]; ok {
			tmpl = templateFromMap
		}
	}

	// If tmpl is still nil, it means:
	// - caching is disabled, OR
	// - template was not found in cache
	// So we build (parse) the template from disk.
	if tmpl == nil {
		newTemplate, err := app.buildTemplateFromDisk(t)
		if err != nil {
			log.Println("Error building template:", err)
			return
		}
		log.Println("building template from disk")
		tmpl = newTemplate
	}

	// If no template data was provided,
	// initialize an empty templateData struct
	// to avoid nil pointer errors in templates.
	if td == nil {
		td = &templateData{}
	}

	// Execute the template:
	// - `w` is the HTTP response writer
	// - `t` is the template name to execute
	// - `td` is the dynamic data passed to the template
	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("Error executing template:", err)

		// Send a 500 Internal Server Error response to the client
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// buildTemplateFromDisk parses templates from files and returns a compiled template.
// This is usually used when caching is disabled or template is not found in cache.
func (app *application) buildTemplateFromDisk(t string) (*template.Template, error) {

	// List of templates to be parsed together.
	// Order matters:
	// - base layout first
	// - shared partials (header/footer)
	// - page-specific template last
	templateSlice := []string{
		"./templates/base.layout.gohtml",
		"./templates/partials/header.partial.gohtml",
		"./templates/partials/footer.partial.gohtml",
		fmt.Sprintf("./templates/%s", t),
	}

	// Parse all template files into a single template object
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		return nil, err
	}

	// Store the compiled template in the map
	// so it can be reused later without re-parsing
	app.templateMap[t] = tmpl

	return tmpl, nil
}
