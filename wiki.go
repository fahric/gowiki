package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("./tmpl/edit.html","./tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
type Page struct {
	Path   string
	Title  string
	Body   string
	Footer string
}

func (p *Page) save() error {
	filename := "./data/"+ p.Path + ".json"
	pageJSON, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, pageJSON, 0600)
	if err != nil {
		return err
	}
	return nil
}

func loadPage(path string) (*Page, error) {
	filename := "./data/"+path + ".json"
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	pageJSON := Page{}
	err = json.Unmarshal(fileContent, &pageJSON)
	if err != nil {
		return nil, err
	}
	return &pageJSON, nil
}

func editHandler(w http.ResponseWriter, r *http.Request,path string) {
	p, err := loadPage(path)
	if err != nil {
		p = &Page{Title: path, Path: path}
	}

	renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request,path string) {
	p, err := loadPage(path)
	if err != nil {
		http.Redirect(w, r, "/edit/"+path, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w,tmpl + ".html",p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request,path string) {
	bodyText := r.FormValue("body")
	titleText := r.FormValue("title")
	footerText := r.FormValue("footer")
	p := &Page{Path: path, Title: titleText, Body: bodyText, Footer: footerText}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+path, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter,*http.Request, string)) http.HandlerFunc  {
	return func(w http.ResponseWriter,r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil{
			http.NotFound(w,r)
			return
		}
		fn(w,r,m[2])
	}
}

func main() {

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":9090", nil)
}
