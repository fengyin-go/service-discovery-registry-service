package handler

import (
	"net/http"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerSearchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/search", s.globalSearch)
}

func (s *Server) globalSearch(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.GlobalSearch(r.URL.Query().Get("q")))
}
