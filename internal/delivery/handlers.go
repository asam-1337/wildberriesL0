package delivery

import (
	"context"
	"errors"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) RootPage(w http.ResponseWriter, r *http.Request) {
	var found bool
	uid := r.FormValue(formOrderKey)
	if uid == "" && r.Method == "POST" {
		log.Info("bad request from user")
		http.Error(w, "bad request from user", http.StatusBadRequest)
		return
	}

	order, err := h.svc.GetOrder(r.Context(), uid)
	if err != nil {
		if errors.Is(err, localErrors.ErrNotFound) {
			log.WithField("order_uid", uid).Info("not found")
		} else {
			log.WithField("err", err.Error()).Error("cant get order")
		}
	} else {
		found = true
	}

	if !found {
		err = h.rnd.RenderRootPage(w, nil)
	} else {
		err = h.rnd.RenderRootPage(w, order)
	}

	if err != nil {
		log.WithField("err", err.Error()).Error("cant render page")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) FindOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue(formOrderKey)
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Info("bad request from user")
		return
	}

	order, err := h.svc.GetOrder(r.Context(), uid)
	if err != nil {
		log.WithField("err", err.Error()).Info(err.Error())
		return
	}

	ctx := context.WithValue(r.Context(), ctxOrderKey, order)
	r = r.WithContext(ctx)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
