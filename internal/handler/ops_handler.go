package handler

import (
	"net/http"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerOpsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/reconcile", s.reconcile)
	mux.HandleFunc("GET /api/path", s.findPath)
	mux.HandleFunc("GET /api/closure/{id}", s.dependencyClosure)
	mux.HandleFunc("GET /api/versions/{name}", s.versionSummary)
	mux.HandleFunc("GET /api/versions/outdated", s.outdatedServices)
	mux.HandleFunc("GET /api/health-score", s.healthScore)
	mux.HandleFunc("GET /api/events/counts", s.eventTypeCounts)
}

func (s *Server) reconcile(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.Reconcile())
}

func (s *Server) findPath(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	httpx.OK(w, s.svc.FindDependencyPath(from, to))
}

func (s *Server) dependencyClosure(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.DependencyClosure(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) versionSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := s.svc.VersionSummary(r.PathValue("name"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, summary)
}

func (s *Server) outdatedServices(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.OutdatedServices())
}

func (s *Server) healthScore(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.HealthScore())
}

func (s *Server) eventTypeCounts(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.EventTypeCounts())
}
