package render

import (
	"html/template"
	"io"
	"log"
)

type Render struct {
	tmpl *template.Template
}

func NewRenderService() *Render {
	r := &Render{}
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatalf(err.Error())
	}

	r.tmpl = tmpl
	return r
}

func (r *Render) RenderRootPage(w io.Writer, order any) error {
	return r.tmpl.ExecuteTemplate(w, "index.html", struct {
		Order any
	}{
		Order: order,
	})
}
