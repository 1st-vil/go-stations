package handler

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"reflect"
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
				w.WriteHeader(http.StatusInternalServerError)
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
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Size = size
		}
		todos, err := h.svc.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := &model.ReadTODOResponse{TODOs: todos}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		req := &model.CreateTODORequest{}
		if err := decoder.Decode(req); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := &model.CreateTODOResponse{TODO: *todo}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		req := &model.UpdateTODORequest{}
		if err := decoder.Decode(req); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if req.ID == 0 || req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
		if err != nil {
			log.Println(err)
			if reflect.TypeOf(err) == reflect.TypeOf(&model.ErrNotFound{}) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		response := &model.UpdateTODOResponse{TODO: *todo}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		decoder := json.NewDecoder(r.Body)
		req := &model.DeleteTODORequest{}
		if err := decoder.Decode(req); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.svc.DeleteTODO(r.Context(), req.IDs); err != nil {
			log.Println(err)
			if reflect.TypeOf(err) == reflect.TypeOf(&model.ErrNotFound{}) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		response := &model.DeleteTODOResponse{}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
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
