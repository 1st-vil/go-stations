package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"math"

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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var req model.ReadTODORequest
		if len(r.URL.Query().Get("prev_id")) == 0 {
			req.PrevID = 0
		} else {
			prev_id, err := strconv.ParseInt(r.URL.Query().Get("prev_id"), 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			req.PrevID = prev_id
		}
		
		if len(r.URL.Query().Get("size")) == 0 {
			req.Size = int64(math.MaxInt64)
		} else {
			size, err := strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			req.Size = size
		}
		todos, err := h.svc.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			log.Println(err)
			return
		}
		response := &model.ReadTODOResponse{todos}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Println(err)
			return
		}
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var todo model.CreateTODORequest
		if err := decoder.Decode(&todo); err != nil {
			log.Println(err)
			return
		}
		if todo.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			todo, err := h.svc.CreateTODO(r.Context(), todo.Subject, todo.Description)
			if err != nil {
				log.Println(err)
				return
			}
			response := &model.CreateTODOResponse{*todo}
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(response); err != nil {
				log.Println(err)
				return
			}
		}
	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var todo model.UpdateTODORequest
		if err := decoder.Decode(&todo); err != nil {
			log.Println(err)
			return
		}
		if todo.ID == 0 || todo.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			todo, err := h.svc.UpdateTODO(r.Context(), todo.ID, todo.Subject, todo.Description)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response := &model.UpdateTODOResponse{*todo}
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(response); err != nil {
				log.Println(err)
				return
			}
		}
	default:

	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
