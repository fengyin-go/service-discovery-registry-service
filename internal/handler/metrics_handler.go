package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/internal/service"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerMetricsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/metrics/instances", s.allInstanceMetrics)
	mux.HandleFunc("GET /api/health-matrix", s.serviceHealthMatrix)
	mux.HandleFunc("POST /api/instances/batch", s.batchRegister)
	mux.HandleFunc("POST /api/instances/batch-deregister", s.batchDeregister)
	mux.HandleFunc("POST /api/health-checks/bulk", s.bulkHealthReport)
}

func (s *Server) allInstanceMetrics(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.AllInstanceMetrics())
}

func (s *Server) serviceHealthMatrix(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ServiceHealthMatrix())
}

type batchRegisterRequest struct {
	ServiceID string           `json:"service_id"`
	Items     []model.Instance `json:"items"`
}

func (s *Server) batchRegister(w http.ResponseWriter, r *http.Request) {
	var req batchRegisterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.BatchRegisterInstances(req.ServiceID, req.Items)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"registered": n})
}

type batchDeregisterRequest struct {
	InstanceIDs []string `json:"instance_ids"`
}

func (s *Server) batchDeregister(w http.ResponseWriter, r *http.Request) {
	var req batchDeregisterRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	httpx.OK(w, map[string]int{"deregistered": s.svc.BatchDeregister(req.InstanceIDs)})
}

type bulkHealthReportRequest struct {
	Reports []service.BulkHealthResult `json:"reports"`
}

func (s *Server) bulkHealthReport(w http.ResponseWriter, r *http.Request) {
	var req bulkHealthReportRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	httpx.OK(w, map[string]int{"reported": s.svc.BulkHealthReport(req.Reports)})
}
