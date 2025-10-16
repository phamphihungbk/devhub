package handler

import (
	"context"
	"encoding/json"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/phamphihungbk/devhub-backend/internal/app/models"
	"github.com/phamphihungbk/devhub-backend/internal/app/repository"
	"github.com/phamphihungbk/devhub-backend/internal/http/handler"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}
type UserHandler struct {
	userRepo  repository.UserRepositoryInterface
	validator *validator.Validate
}

func NewUserHandler(repo repository.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepo:  repo,
		validator: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		handler.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validator.Struct(req); err != nil {
		handler.ErrorResponse(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  models.UserRole(req.Role),
	}
	ctx := h.context(r)
	if err := h.userRepo.Create(ctx, &user); err != nil {
		handler.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.JSONResponse(w, http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		handler.ErrorResponse(w, http.StatusBadRequest, "missing id")
		return
	}
	var req UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		handler.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validator.Struct(req); err != nil {
		handler.ErrorResponse(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}
	ctx := h.context(r)
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil {
		handler.ErrorResponse(w, http.StatusNotFound, "user not found")
		return
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = models.UserRole(req.Role)
	}
	if err := h.userRepo.Update(ctx, user); err != nil {
		handler.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.JSONResponse(w, http.StatusOK, user)
}
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := h.context(r)
	users, err := h.userRepo.List(ctx)
	if err != nil {
		handler.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.JSONResponse(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		handler.ErrorResponse(w, http.StatusBadRequest, "missing id")
		return
	}
	ctx := h.context(r)
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil {
		handler.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	handler.JSONResponse(w, http.StatusOK, user)
}

// decodeJSON is a helper to decode JSON body
func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// context returns a context for the request (can be extended for tracing, auth, etc)
func (h *UserHandler) context(r *http.Request) context.Context {
	return context.Background()
}
