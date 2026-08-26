package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerLoadBalanceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/load-balances", s.createLoadBalance)
	mux.HandleFunc("GET /api/load-balances", s.listLoadBalances)
	mux.HandleFunc("GET /api/load-balances/by-service/{serviceID}", s.getLoadBalanceByService)
	mux.HandleFunc("GET /api/load-balances/{id}", s.getLoadBalance)
	mux.HandleFunc("PUT /api/load-balances/{id}", s.updateLoadBalance)
	mux.HandleFunc("DELETE /api/load-balances/{id}", s.deleteLoadBalance)
}

type createLoadBalanceRequest struct {
	ServiceID string `json:"service_id"`
	Algorithm string `json:"algorithm"`
}

func (s *Server) createLoadBalance(w http.ResponseWriter, r *http.Request) {
	var req createLoadBalanceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.CreateLoadBalanceStrategy(model.LoadBalanceStrategy{
		ServiceID: req.ServiceID,
		Algorithm: req.Algorithm,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, l)
}

func (s *Server) listLoadBalances(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListLoadBalanceStrategies())
}

func (s *Server) getLoadBalanceByService(w http.ResponseWriter, r *http.Request) {
	l, err := s.svc.GetLoadBalanceByService(r.PathValue("serviceID"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) getLoadBalance(w http.ResponseWriter, r *http.Request) {
	l, err := s.svc.GetLoadBalanceStrategy(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

type updateLoadBalanceRequest struct {
	Algorithm string `json:"algorithm"`
}

func (s *Server) updateLoadBalance(w http.ResponseWriter, r *http.Request) {
	var req updateLoadBalanceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.UpdateLoadBalanceStrategy(r.PathValue("id"), req.Algorithm)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) deleteLoadBalance(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteLoadBalanceStrategy(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
