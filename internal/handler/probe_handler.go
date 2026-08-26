package handler

import (
	"net/http"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerProbeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/probe", s.probeAll)
	mux.HandleFunc("GET /api/probe/{id}", s.probeInstance)
}

func (s *Server) probeAll(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ProbeAll())
}

func (s *Server) probeInstance(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.ProbeInstance(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}
