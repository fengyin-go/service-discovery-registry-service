package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerLabelRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/labels", s.createLabel)
	mux.HandleFunc("GET /api/labels", s.listLabels)
	mux.HandleFunc("GET /api/labels/match", s.listServicesByLabel)
	mux.HandleFunc("PUT /api/labels/{id}", s.updateLabel)
	mux.HandleFunc("DELETE /api/labels/{id}", s.deleteLabel)
	mux.HandleFunc("GET /api/services/{id}/labels", s.listServiceLabels)
}

type createLabelRequest struct {
	ServiceID string `json:"service_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (s *Server) createLabel(w http.ResponseWriter, r *http.Request) {
	var req createLabelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.CreateLabel(req.ServiceID, req.Key, req.Value)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, l)
}

func (s *Server) listLabels(w http.ResponseWriter, r *http.Request) {
	filter := model.LabelFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Key:       r.URL.Query().Get("key"),
		Value:     r.URL.Query().Get("value"),
	}
	httpx.OK(w, s.svc.ListLabels(filter))
}

func (s *Server) listServicesByLabel(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	httpx.OK(w, s.svc.ListServicesByLabel(key, value))
}

func (s *Server) listServiceLabels(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListLabelsByService(r.PathValue("id")))
}

type updateLabelRequest struct {
	Value string `json:"value"`
}

func (s *Server) updateLabel(w http.ResponseWriter, r *http.Request) {
	var req updateLabelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.UpdateLabel(r.PathValue("id"), req.Value)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) deleteLabel(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteLabel(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
