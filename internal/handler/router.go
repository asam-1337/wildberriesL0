package handler

import "net/http"

type Handler struct {
}

func (h *Handler) InitRoutes() *http.ServeMux {
	s := &http.ServeMux{}
	return s
}
