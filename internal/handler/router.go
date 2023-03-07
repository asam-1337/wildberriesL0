package handler

import (
	"github.com/asam-1337/wildberriesL0/internal/render"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	orderUidFormKey = "order_uid"
)

type Handler struct {
	rnd *render.Render
}

func NewHandler(rnd *render.Render) *Handler {
	return &Handler{
		rnd: rnd,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.RootPage)
	r.HandleFunc("/find", h.FindOrder).Methods("POST")
	r.HandleFunc("/find", h.FindPage).Methods("GET")
	r.HandleFunc("/delivery/{id}", h.Delivery).Methods("GET")
	return r
}

func (h *Handler) RootPage(w http.ResponseWriter, r *http.Request) {
	err := h.rnd.RenderRootPage(w, nil)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
	}
}

func (h *Handler) Delivery(w http.ResponseWriter, r *http.Request) {
	err := h.rnd.RenderFindPage(w)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
	}
}

func (h *Handler) FindPage(w http.ResponseWriter, r *http.Request) {
	err := h.rnd.RenderFindPage(w)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) FindOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue(orderUidFormKey)
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
