package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/services", s.createService)
	mux.HandleFunc("GET /api/services", s.listServices)
	mux.HandleFunc("GET /api/services/{id}", s.getService)
	mux.HandleFunc("PUT /api/services/{id}", s.updateService)
	mux.HandleFunc("PATCH /api/services/{id}/status", s.changeServiceStatus)
	mux.HandleFunc("DELETE /api/services/{id}", s.deleteService)
	mux.HandleFunc("GET /api/services/{id}/versions", s.listServiceVersions)
	mux.HandleFunc("GET /api/services/{id}/latest", s.latestService)
	mux.HandleFunc("GET /api/services/{id}/stats", s.serviceStats)
	mux.HandleFunc("GET /api/services/{id}/select", s.selectInstance)
}

type createServiceRequest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
}

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sv, err := s.svc.CreateService(model.Service{
		Name:        req.Name,
		Version:     req.Version,
		Description: req.Description,
		Owner:       req.Owner,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sv)
}

func (s *Server) listServices(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ServiceFilter{
		Keyword: r.URL.Query().Get("keyword"),
		Status:  r.URL.Query().Get("status"),
		Owner:   r.URL.Query().Get("owner"),
	}
	items, total, err := s.svc.ListServices(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	sv, err := s.svc.GetService(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sv)
}

func (s *Server) updateService(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sv, err := s.svc.UpdateService(r.PathValue("id"), model.Service{
		Name:        req.Name,
		Version:     req.Version,
		Description: req.Description,
		Owner:       req.Owner,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sv)
}

type changeServiceStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) changeServiceStatus(w http.ResponseWriter, r *http.Request) {
	var req changeServiceStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sv, err := s.svc.ChangeServiceStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sv)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteService(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listServiceVersions(w http.ResponseWriter, r *http.Request) {
	sv, err := s.svc.GetService(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, s.svc.ListServiceVersions(sv.Name))
}

func (s *Server) latestService(w http.ResponseWriter, r *http.Request) {
	sv, err := s.svc.GetService(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	latest, err := s.svc.LatestService(sv.Name)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, latest)
}

func (s *Server) serviceStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.StatsService(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) selectInstance(w http.ResponseWriter, r *http.Request) {
	inst, err := s.svc.SelectInstance(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, inst)
}
