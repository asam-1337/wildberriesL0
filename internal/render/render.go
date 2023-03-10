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

func (r *Render) RenderRootPage(w io.Writer, orders entity.Order, exist bool) error {
	return r.tmpl.ExecuteTemplate(w, "index.html", struct {
		Order entity.Order
		Exist bool
	}{
		Order: orders,
		Exist: exist,
	})
}
