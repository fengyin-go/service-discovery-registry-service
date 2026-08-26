package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerInstanceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/instances", s.registerInstance)
	mux.HandleFunc("GET /api/instances", s.listInstances)
	mux.HandleFunc("GET /api/instances/{id}", s.getInstance)
	mux.HandleFunc("PATCH /api/instances/{id}/status", s.changeInstanceStatus)
	mux.HandleFunc("DELETE /api/instances/{id}", s.deregisterInstance)
	mux.HandleFunc("GET /api/instances/{id}/health", s.instanceHealth)
}

type registerInstanceRequest struct {
	ServiceID string            `json:"service_id"`
	Host      string            `json:"host"`
	Port      int               `json:"port"`
	Weight    int               `json:"weight"`
	Status    string            `json:"status"`
	Region    string            `json:"region"`
	Zone      string            `json:"zone"`
	Metadata  map[string]string `json:"metadata"`
}

func (s *Server) registerInstance(w http.ResponseWriter, r *http.Request) {
	var req registerInstanceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inst, err := s.svc.RegisterInstance(model.Instance{
		ServiceID: req.ServiceID,
		Host:      req.Host,
		Port:      req.Port,
		Weight:    req.Weight,
		Status:    req.Status,
		Region:    req.Region,
		Zone:      req.Zone,
		Metadata:  req.Metadata,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, inst)
}

func (s *Server) listInstances(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.InstanceFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Status:    r.URL.Query().Get("status"),
		Region:    r.URL.Query().Get("region"),
		Zone:      r.URL.Query().Get("zone"),
	}
	items, total, err := s.svc.ListInstances(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getInstance(w http.ResponseWriter, r *http.Request) {
	inst, err := s.svc.GetInstance(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, inst)
}

type changeInstanceStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) changeInstanceStatus(w http.ResponseWriter, r *http.Request) {
	var req changeInstanceStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inst, err := s.svc.ChangeInstanceStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, inst)
}

func (s *Server) deregisterInstance(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeregisterInstance(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) instanceHealth(w http.ResponseWriter, r *http.Request) {
	view, err := s.svc.InstanceHealth(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, view)
}
