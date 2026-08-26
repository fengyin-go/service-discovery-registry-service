package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerHealthCheckRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/health-checks", s.createHealthCheck)
	mux.HandleFunc("GET /api/health-checks", s.listHealthChecks)
	mux.HandleFunc("GET /api/health-checks/by-instance/{instanceID}", s.getHealthCheckByInstance)
	mux.HandleFunc("GET /api/health-checks/{id}", s.getHealthCheck)
	mux.HandleFunc("PUT /api/health-checks/{id}", s.updateHealthCheck)
	mux.HandleFunc("DELETE /api/health-checks/{id}", s.deleteHealthCheck)
	mux.HandleFunc("POST /api/health-checks/{id}/report", s.reportCheckResult)
}

type createHealthCheckRequest struct {
	InstanceID  string `json:"instance_id"`
	Type        string `json:"type"`
	Target      string `json:"target"`
	IntervalSec int    `json:"interval_sec"`
	TimeoutSec  int    `json:"timeout_sec"`
	Threshold   int    `json:"threshold"`
}

func (s *Server) createHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req createHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.CreateHealthCheck(model.HealthCheck{
		InstanceID:  req.InstanceID,
		Type:        req.Type,
		Target:      req.Target,
		IntervalSec: req.IntervalSec,
		TimeoutSec:  req.TimeoutSec,
		Threshold:   req.Threshold,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, h)
}

func (s *Server) listHealthChecks(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListHealthChecks())
}

func (s *Server) getHealthCheckByInstance(w http.ResponseWriter, r *http.Request) {
	h, err := s.svc.GetHealthCheckByInstance(r.PathValue("instanceID"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

func (s *Server) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	h, err := s.svc.GetHealthCheck(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

func (s *Server) updateHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req createHealthCheckRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.UpdateHealthCheck(r.PathValue("id"), model.HealthCheck{
		Type:        req.Type,
		Target:      req.Target,
		IntervalSec: req.IntervalSec,
		TimeoutSec:  req.TimeoutSec,
		Threshold:   req.Threshold,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

func (s *Server) deleteHealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteHealthCheck(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type reportCheckResultRequest struct {
	Status string `json:"status"`
}

func (s *Server) reportCheckResult(w http.ResponseWriter, r *http.Request) {
	var req reportCheckResultRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.GetHealthCheck(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	updated, err := s.svc.ReportCheckResult(h.InstanceID, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, updated)
}
