package handler

import (
	"net/http"
	"strconv"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerHeartbeatRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/heartbeats", s.sendHeartbeat)
	mux.HandleFunc("GET /api/heartbeats", s.listHeartbeats)
	mux.HandleFunc("POST /api/maintenance/run", s.runMaintenance)
	mux.HandleFunc("GET /api/maintenance/stale", s.listStale)
}

type sendHeartbeatRequest struct {
	InstanceID string `json:"instance_id"`
	LatencyMs  int64  `json:"latency_ms"`
	Result     string `json:"result"`
}

func (s *Server) sendHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req sendHeartbeatRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.SendHeartbeat(req.InstanceID, req.LatencyMs, req.Result)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, h)
}

func (s *Server) listHeartbeats(w http.ResponseWriter, r *http.Request) {
	filter := model.HeartbeatFilter{
		InstanceID: r.URL.Query().Get("instance_id"),
		Result:     r.URL.Query().Get("result"),
	}
	httpx.OK(w, s.svc.ListHeartbeats(filter))
}

func (s *Server) runMaintenance(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.RunMaintenance())
}

func (s *Server) listStale(w http.ResponseWriter, r *http.Request) {
	timeout, _ := strconv.Atoi(r.URL.Query().Get("timeout"))
	if timeout <= 0 {
		timeout = 30
	}
	httpx.OK(w, s.svc.StaleInstances(timeout))
}
