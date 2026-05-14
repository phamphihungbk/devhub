package handler

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"devhub-backend/internal/domain/errs"
	releaseUsecase "devhub-backend/internal/usecase/release"
	"devhub-backend/internal/util/httpresponse"
	"devhub-backend/internal/util/misc"

	"github.com/gin-gonic/gin"
)

const ciWebhookTokenHeader = "X-Devhub-Webhook-Token"

type createReleaseFromCIRequest struct {
	RepoURL     string `json:"repo_url" example:"https://gitea.devhub.local/acme/service.git"`
	RepoOwner   string `json:"repo_owner" example:"acme"`
	RepoName    string `json:"repo_name" example:"service"`
	Tag         string `json:"tag" example:"v1.0.0" binding:"required"`
	Target      string `json:"target" example:"main"`
	Name        string `json:"name" example:"v1.0.0"`
	Notes       string `json:"notes" example:"Built by CI."`
	HTMLURL     string `json:"html_url" example:"https://gitea.devhub.local/acme/service/releases/tag/v1.0.0"`
	Image       string `json:"image" example:"registry.devhub.local/service:v1.0.0"`
	ExternalRef string `json:"external_ref" example:"registry.devhub.local/service:v1.0.0"`
}

func (h *releaseHandler) CreateReleaseFromCI(c *gin.Context) {
	if h.appConfig.WebhookToken != "" {
		token := strings.TrimSpace(c.GetHeader(ciWebhookTokenHeader))
		if subtle.ConstantTimeCompare([]byte(token), []byte(h.appConfig.WebhookToken)) != 1 {
			httpresponse.Error(c, errs.NewUnauthorizedError("invalid webhook token", nil))
			return
		}
	}

	var input createReleaseFromCIRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		err = misc.WrapError(err, errs.NewBadRequestError("unable to parse request", map[string]string{"details": err.Error()}))
		httpresponse.Error(c, err)
		return
	}

	release, err := h.releaseUsecase.CreateReleaseFromCI(c.Request.Context(), releaseUsecase.CreateReleaseFromCIInput{
		RepoURL:     input.RepoURL,
		RepoOwner:   input.RepoOwner,
		RepoName:    input.RepoName,
		Tag:         input.Tag,
		Target:      input.Target,
		Name:        input.Name,
		Notes:       input.Notes,
		HTMLURL:     input.HTMLURL,
		Image:       input.Image,
		ExternalRef: input.ExternalRef,
	})
	if err != nil {
		httpresponse.Error(c, err)
		return
	}

	httpresponse.SuccessWithStatus(c, http.StatusCreated, h.newCreateReleaseResponse(release))
}
