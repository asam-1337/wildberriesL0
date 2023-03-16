package delivery

import (
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
