package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/phamphihungbk/devhub-backend/internal/app/repository"
)

var projectRepo repository.ProjectRepositoryInterface

func SetProjectRepository(repo repository.ProjectRepositoryInterface) {
	projectRepo = repo
}

func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func ListServiceTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	if projectRepo == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("repository not set"))
		return
	}
	ctx := context.Background()
	projects, err := projectRepo.List(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
func GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
