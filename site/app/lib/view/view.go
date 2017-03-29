package view

import (
	"github.com/itzaname/anime-streaming/site/app/lib/session"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// View attributes
type View struct {
	BaseURI   string
	Extension string
	Folder    string
	Name      string
	Caching   bool
	Vars      map[string]interface{}
	request   *http.Request
	funcs     template.FuncMap
}

// Template root and children
type Template struct {
	Root     string   `json:"Root"`
	Children []string `json:"Children"`
}

var (
	mutex              sync.RWMutex
	viewInfo           View
	childTemplates     []string
	rootTemplate       string
	templateCollection = make(map[string]*template.Template)
)

// Configure sets the view information
func Configure(vi View) {
	viewInfo = vi
}

// New returns a new view
func New(req *http.Request, fms ...template.FuncMap) *View {
	v := &View{}
	v.Vars = make(map[string]interface{})
	v.Vars["AuthLevel"] = "anon"

	v.BaseURI = viewInfo.BaseURI
	v.Extension = viewInfo.Extension
	v.Folder = viewInfo.Folder
	v.Name = viewInfo.Name

	// Make sure BaseURI is available in the templates
	v.Vars["BaseURI"] = v.BaseURI

	// This is required for the view to access the request
	v.request = req

	// Get session
	sess := session.Instance(v.request)

	// Set the AuthLevel to auth if the user is logged in
	if sess.Values["maluser"] != nil {
		v.Vars["username"] = sess.Values["maluser"]
	}

	/// Final FuncMap
	fm := make(template.FuncMap)

	// Loop through the maps
	for _, m := range fms {
		// Loop through each key and value
		for k, v := range m {
			fm[k] = v
		}
	}

	v.funcs = fm

	return v
}

// RenderSingle renders a template to the writer
func (v *View) RenderSingle(w http.ResponseWriter) {
	templateList := []string{v.Name}

	// Loop through each template and test the full path
	for i, name := range templateList {
		// Get the absolute path of the root template
		path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
		if err != nil {
			http.Error(w, "Template Path Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		templateList[i] = path
	}

	// Determine if there is an error in the template syntax
	templates, err := template.New(v.Name).Funcs(v.funcs).ParseFiles(templateList...)

	if err != nil {
		http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the template collection
	tc := templates

	// Display the content to the screen
	err = tc.Funcs(v.funcs).ExecuteTemplate(w, v.Name+"."+v.Extension, v.Vars)

	if err != nil {
		http.Error(w, "Template File Error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Render renders a template to the writer
func (v *View) Render(w http.ResponseWriter) {

	// Get the template collection from cache
	mutex.RLock()
	tc, ok := templateCollection[v.Name]
	mutex.RUnlock()

	// If the template collection is not cached or caching is disabled
	if !ok || !viewInfo.Caching {

		// List of template names
		var templateList []string
		templateList = append(templateList, rootTemplate)
		templateList = append(templateList, v.Name)
		templateList = append(templateList, childTemplates...)

		// Loop through each template and test the full path
		for i, name := range templateList {
			// Get the absolute path of the root template
			path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
			if err != nil {
				http.Error(w, "Template Path Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			templateList[i] = path
		}

		// Determine if there is an error in the template syntax
		templates, err := template.New(v.Name).Funcs(v.funcs).ParseFiles(templateList...)

		if err != nil {
			http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Cache the template collection
		mutex.Lock()
		templateCollection[v.Name] = templates
		mutex.Unlock()

		// Save the template collection
		tc = templates
	}

	// Display the content to the screen
	err := tc.Funcs(v.funcs).ExecuteTemplate(w, rootTemplate+"."+v.Extension, v.Vars)

	if err != nil {
		http.Error(w, "Template File Error: "+err.Error(), http.StatusInternalServerError)
	}
}

// LoadTemplates will set the root and child templates
func LoadTemplates(rootTemp string, childTemps []string) {
	rootTemplate = rootTemp
	childTemplates = childTemps
}
