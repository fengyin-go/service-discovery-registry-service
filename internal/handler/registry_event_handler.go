package handler

import (
	"net/http"

	"serviceregistry/internal/model"
	"serviceregistry/pkg/httpx"
)

func (s *Server) registerEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/events", s.listEvents)
}

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EventFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		EventType: r.URL.Query().Get("event_type"),
	}
	items, total, err := s.svc.ListEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
