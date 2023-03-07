package render

import (
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
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

func (r *Render) RenderRootPage(w io.Writer, orders []entity.Order) error {
	return r.tmpl.ExecuteTemplate(w, "index.html", struct {
		Orders []entity.Order
	}{
		Orders: orders,
	})
}

func (r *Render) RenderFindPage(w io.Writer) error {
	return r.tmpl.ExecuteTemplate(w, "find.html", nil)
}
