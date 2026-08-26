package handler

import (
	"net/http"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerTopologyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/topology", s.topology)
	mux.HandleFunc("GET /api/topology/cycles", s.detectCycles)
}

func (s *Server) topology(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.Topology())
}

func (s *Server) detectCycles(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.DetectCycles())
}
