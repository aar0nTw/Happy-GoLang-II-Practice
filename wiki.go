package main

import ( 
  "bytes"
  "html/template"
  "path/filepath"
  "io"
  "io/ioutil"
  "net/http"
  "log"
)

const ArticlePath string = "articles"
const ViewPath string = "views"
const LayoutPath string = "layout.html"
var TemplateDict *template.Template

type Page struct {
 Title string
 Body template.HTML
}

func (p *Page) save() error {
  filename := filepath.Join(ArticlePath, p.Title + ".txt")
  return ioutil.WriteFile(filename, []byte(p.Body), 0600)
}

type Layout struct {
  Tmpl string
  Content template.HTML
}

// Render function
//
func (layout *Layout) render(w io.Writer, c template.HTML) {
  layout.Content = c
  TemplateDict.ExecuteTemplate(w, layout.Tmpl, layout)
}

func Log(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
        handler.ServeHTTP(w, r)
    })
}

func loadPage(title string) (*Page, error) {
  filename := filepath.Join(ArticlePath, title + ".txt")
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: template.HTML(string(body))}, nil
}

func loadLayout(tmpl string) (*Layout, error) {
  return &Layout{Tmpl: tmpl}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  buf := new(bytes.Buffer)
  TemplateDict.ExecuteTemplate(buf, "view.html", p)
  layout, _ := loadLayout(LayoutPath)
  layout.render(w, template.HTML(buf.String()))
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }

  buf := new(bytes.Buffer)
  TemplateDict.ExecuteTemplate(buf, "edit.html", p)
  layout, _ := loadLayout(LayoutPath)
  layout.render(w, template.HTML(buf.String()))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p := &Page{Title: title, Body: template.HTML(body)}
  p.save()
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Cache the Templates file
func init() {
  TemplateDict, _ = template.ParseGlob("views/*.html")
}

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/index", http.StatusFound)
  })
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", Log(http.DefaultServeMux))
}
