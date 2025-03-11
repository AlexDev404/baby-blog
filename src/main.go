// Filename: main.go

package main

import (
	"bytes"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"baby-blog/types"
	"html/template"
	"sync"
)

// Application is a wrapper for types.Application
type Application struct {
	*types.Application
	templates  *template.Template
	bufferPool sync.Pool
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, t *template.Template) {
	// Get the URL path
	path := r.URL.Path[1:]
	disallowedRoutes := []string{"context", "head", "header", "footer", "current_ctx", "index"}
	// Remove any trailing slashes
	if path == "" {
		path = "index"
	} else {
		path = strings.TrimSuffix(path, "/")
		// Check if path is in disallowedRoutes
		for _, route := range disallowedRoutes {
			if path == route {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
	}

	// TemplateData is a struct that holds the title, body, and data for the template
	// First try to get the template by path if it's not root
	var templateContent string
	var tmpl *template.Template = t.Lookup(path)
	if path != "/" {
		if tmpl != nil {
			var buf bytes.Buffer
			tmpl.Execute(&buf, nil)
			templateContent = buf.String()
		} else {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "./static/errors/404.html")
			return
		}
	}

	data := &types.TemplateData{
		Data: map[string]interface{}{
			"Path": path,
			"HTML": template.HTML(templateContent),
		},
	}
	// Section: Render Layouts
	// First: Let's check if there's a layout for the path
	// Remove the leading text following the last / in the string
	path = strings.TrimSuffix(path, "/"+path[strings.LastIndex(path, "/")+1:])
	layout := t.Lookup("layout/" + path)

	// Page buffer
	pageBuf := app.bufferPool.Get().(*bytes.Buffer)
	pageBuf.Reset()
	defer app.bufferPool.Put(pageBuf)

	// Apply the layout
	var err error
	if layout == nil {
		// Render the template directly
		err = t.ExecuteTemplate(pageBuf, "layout/app", data)
	} else {
		// Render the template with the layout
		err = layout.Execute(pageBuf, data)
	}
	if err != nil {
		app.Logger.Error("Error rendering page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = pageBuf.WriteTo(w)
	if err != nil {
		app.Logger.Error("Failed to write template to response: ", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "./static/errors/500.html")
	}

}

func getTemplates() (*template.Template, error) {
	// Parse the initial templates
	log.Println("Parsing 'initial' templates...")
	templates, err := template.ParseGlob("templates/*.mustache")
	if err != nil {
		return nil, err
	}
	// Add the partials to the templates
	log.Println("Parsing 'partial' templates...")
	templates, err = templates.ParseGlob("templates/partials/*.mustache")
	if err != nil {
		return nil, err
	}
	// Add the routes to the templates
	log.Println("Parsing 'route' templates...")
	err = filepath.Walk("templates/routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".mustache") {
			var tmpl *template.Template
			tmpl, err = templates.ParseFiles(path)
			if err != nil {
				return err
			}
			templates = tmpl
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// Starts the web server and listens for requests
func main() {
	addr := flag.String("addr", "4000", "HTTP network address")
	flag.Parse()
	mux := http.NewServeMux()
	// Serve the static files
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	templates, t_err := getTemplates()
	if t_err != nil {
		log.Fatalf("Error parsing templates: %v", t_err)
	}
	typesApp := &types.Application{
		Logger: logger,
	}

	// Initialize the application
	app := &Application{
		Application: typesApp,
		templates:   templates,
	}

	// Initialize the buffer pool
	app.bufferPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}

	log.Printf("Templates loaded: %v", templates.DefinedTemplates())
	log.Printf("App templates: %v", app.templates.DefinedTemplates())

	// Register the handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, app.templates)
	})

	app.Logger.Info("Now listening on port " + *addr)

	// Start listening for requests (start the web server)
	err := http.ListenAndServe((":" + *addr), app.Middleware.LoggingMiddleware(mux))

	// Log error message if server quits unexpectedly
	if err != nil {
		panic(err.Error())
	}
}
