package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerDependencyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/dependencies", s.createDependency)
	mux.HandleFunc("GET /api/dependencies", s.listDependencies)
	mux.HandleFunc("GET /api/dependencies/of/{serviceID}", s.dependenciesOf)
	mux.HandleFunc("GET /api/dependencies/on/{serviceID}", s.dependentsOf)
	mux.HandleFunc("GET /api/dependencies/{id}", s.getDependency)
	mux.HandleFunc("POST /api/dependencies/{id}/call", s.incrementCallCount)
	mux.HandleFunc("DELETE /api/dependencies/{id}", s.deleteDependency)
}

type createDependencyRequest struct {
	ServiceID   string `json:"service_id"`
	DependsOnID string `json:"depends_on_id"`
	Protocol    string `json:"protocol"`
}

func (s *Server) createDependency(w http.ResponseWriter, r *http.Request) {
	var req createDependencyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDependency(model.ServiceDependency{
		ServiceID:   req.ServiceID,
		DependsOnID: req.DependsOnID,
		Protocol:    req.Protocol,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDependencies(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListDependencies())
}

func (s *Server) dependenciesOf(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.DependenciesOf(r.PathValue("serviceID")))
}

func (s *Server) dependentsOf(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.DependentsOf(r.PathValue("serviceID")))
}

func (s *Server) getDependency(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDependency(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) incrementCallCount(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.IncrementCallCount(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDependency(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteDependency(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
