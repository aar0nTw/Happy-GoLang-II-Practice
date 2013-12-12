package main

import ( 
  "bytes"
  "html/template"
  "path/filepath"
  "io"
  "io/ioutil"
  "net/http"
)

const ArticlePath string = "articles"
const ViewPath string = "views"
const LayoutPath string = "views/layout.html"

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

func (layout *Layout) render(w io.Writer, c template.HTML) {
  l, _ := template.ParseFiles(layout.Tmpl)
  layout.Content = c
  l.Execute(w, layout)
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
  t, _ := template.ParseFiles(filepath.Join(ViewPath, "view.html"))
  buf := new(bytes.Buffer)
  t.Execute(buf, p)
  layout, _ := loadLayout(LayoutPath)
  layout.render(w, template.HTML(buf.String()))
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }

  t, _ := template.ParseFiles(filepath.Join(ViewPath,"edit.html"))
  buf := new(bytes.Buffer)
  t.Execute(buf, p)
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

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/index", http.StatusFound)
  })
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}
