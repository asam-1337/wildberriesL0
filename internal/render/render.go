package render

import (
	"html/template"
	"io"
)

type Render struct {
	tmpl *template.Template
}

func NewRenderService() (*Render, error) {
	r := &Render{}
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		return nil, err
	}

	r.tmpl = tmpl
	return r, nil
}

func (r *Render) RenderRootPage(w io.Writer, order any) error {
	return r.tmpl.ExecuteTemplate(w, "index.html", struct {
		Order any
	}{
		Order: order,
	})
}
