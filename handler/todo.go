package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{Todos: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) renderError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

// ServeHTTP implements http.Handler to accept HTTP requests for TODO endpoints.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		var req model.ReadTODORequest
		prevIDStr := r.URL.Query().Get("prev_id")
		if prevIDStr != "" {
			prevID, err := strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				h.renderError(w, "invalid prev_id", http.StatusBadRequest)
				return
			}
			req.PrevID = prevID
		}
		sizeStr := r.URL.Query().Get("size")
		if sizeStr != "" {
			size, err := strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				h.renderError(w, "invalid size", http.StatusBadRequest)
				return
			}
			req.Size = size
		}

		resp, err := h.Read(ctx, &req)
		if err != nil {
			h.renderError(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)

	case http.MethodPost:
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.renderError(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			h.renderError(w, "subject is required", http.StatusBadRequest)
			return
		}

		resp, err := h.Create(ctx, &req)
		if err != nil {
			h.renderError(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)

	case http.MethodPut:
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.renderError(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.ID == 0 {
			h.renderError(w, "id is required", http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			h.renderError(w, "subject is required", http.StatusBadRequest)
			return
		}

		resp, err := h.Update(ctx, &req)
		if err != nil {
			var errNotFound *model.ErrNotFound
			if errors.As(err, &errNotFound) {
				h.renderError(w, "not found", http.StatusNotFound)
				return
			}
			h.renderError(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)

	case http.MethodDelete:
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.renderError(w, "bad request", http.StatusBadRequest)
			return
		}

		if len(req.IDs) == 0 {
			h.renderError(w, "ids must not be empty", http.StatusBadRequest)
			return
		}

		resp, err := h.Delete(ctx, &req)
		if err != nil {
			var errNotFound *model.ErrNotFound
			if errors.As(err, &errNotFound) {
				h.renderError(w, "not found", http.StatusNotFound)
				return
			}
			h.renderError(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)

	default:
		h.renderError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
