package handler

import (
	"net/http"
	"strconv"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/registry", s.statsRegistry)
	mux.HandleFunc("GET /api/stats/top-dependencies", s.topDependencies)
}

func (s *Server) statsRegistry(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsRegistry())
}

func (s *Server) topDependencies(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	httpx.OK(w, s.svc.TopDependencies(n))
}
