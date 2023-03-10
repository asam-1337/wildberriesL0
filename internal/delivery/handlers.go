package delivery

import (
	"context"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"log"
	"net/http"
)

func (h *Handler) RootPage(w http.ResponseWriter, r *http.Request) {
	order, ok := r.Context().Value(ctxOrderKey).(entity.Order)
	log.Printf("context: %v", ok)
	err := h.rnd.RenderRootPage(w, order, ok)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) FindOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue(formOrderKey)
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := h.svc.GetOrder(r.Context(), uid)
	if err != nil {
		return
	}

	ctx := context.WithValue(r.Context(), ctxOrderKey, order)
	r = r.WithContext(ctx)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
