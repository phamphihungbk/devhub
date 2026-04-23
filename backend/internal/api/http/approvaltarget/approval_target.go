package approvaltarget

import (
	"fmt"
	"net/http"

	"devhub-backend/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

var routeApprovalTargets = map[string]entity.ApprovalTarget{
	approvalRouteKey(http.MethodPost, "/projects/:project/scaffold-requests"): {
		Resource: entity.ApprovalResourceScaffoldRequest,
		Action:   entity.ApprovalActionCreate,
	},
	approvalRouteKey(http.MethodDelete, "/scaffold-requests/:scaffold-request"): {
		Resource: entity.ApprovalResourceScaffoldRequest,
		Action:   entity.ApprovalActionDelete,
	},
	approvalRouteKey(http.MethodPost, "/services/:service/deployments"): {
		Resource: entity.ApprovalResourceDeployment,
		Action:   entity.ApprovalActionCreate,
	},
	approvalRouteKey(http.MethodDelete, "/deployments/:deployment"): {
		Resource: entity.ApprovalResourceDeployment,
		Action:   entity.ApprovalActionDelete,
	},
	approvalRouteKey(http.MethodPatch, "/deployments/:deployment"): {
		Resource: entity.ApprovalResourceDeployment,
		Action:   entity.ApprovalActionUpdate,
	},
	approvalRouteKey(http.MethodPost, "/services/:service/releases"): {
		Resource: entity.ApprovalResourceRelease,
		Action:   entity.ApprovalActionCreate,
	},
}

func approvalRouteKey(method, route string) string {
	return fmt.Sprintf("%s %s", method, route)
}

func ApprovalTargetFromRoute(method, route string) (entity.ApprovalTarget, bool) {
	target, ok := routeApprovalTargets[approvalRouteKey(method, route)]
	return target, ok
}

func ApprovalTargetFromContext(c *gin.Context) (entity.ApprovalTarget, bool) {
	if c == nil {
		return entity.ApprovalTarget{}, false
	}

	fullPath := c.FullPath()
	if fullPath == "" {
		return entity.ApprovalTarget{}, false
	}

	return ApprovalTargetFromRoute(c.Request.Method, fullPath)
}

func StringsFromContext(c *gin.Context) (resource, action string) {
	target, ok := ApprovalTargetFromContext(c)
	if !ok {
		return "", ""
	}

	return target.Resource.String(), target.Action.String()
}
