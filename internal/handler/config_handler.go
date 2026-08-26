package handler

import (
	"net/http"

	"serviceregistry/pkg/httpx"
)

func (s *Server) registerConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/configs", s.listConfigs)
	mux.HandleFunc("PUT /api/configs", s.setConfig)
	mux.HandleFunc("GET /api/configs/{key}", s.getConfig)
	mux.HandleFunc("DELETE /api/configs/{key}", s.deleteConfig)
}

type setConfigRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (s *Server) listConfigs(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListConfigs())
}

func (s *Server) setConfig(w http.ResponseWriter, r *http.Request) {
	var req setConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.SetConfig(req.Key, req.Value, req.Description)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetConfig(r.PathValue("key"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteConfig(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteConfig(r.PathValue("key")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
