package handler

import (
	"net/http"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/regions", s.regionReport)
	mux.HandleFunc("GET /api/reports/algorithms", s.algorithmReport)
	mux.HandleFunc("GET /api/reports/events", s.eventTypeReport)
}

func (s *Server) regionReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RegionReport())
}

func (s *Server) algorithmReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.AlgorithmReport())
}

func (s *Server) eventTypeReport(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.EventTypeReport())
}
