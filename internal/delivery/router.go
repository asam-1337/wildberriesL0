package delivery

import (
	"context"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/render"
	"github.com/gorilla/mux"
)

const (
	formOrderKey = "order_uid"
	ctxOrderKey  = "order"
)

type Service interface {
	GetOrder(ctx context.Context, uid string) (entity.Order, error)
}

type Handler struct {
	rnd *render.Render
	svc Service
}

func NewHandler(rnd *render.Render, svc Service) *Handler {
	return &Handler{
		rnd: rnd,
		svc: svc,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.RootPage)
	r.HandleFunc("/find", h.FindOrder).Methods("POST")
	return r
}
