package approvaltarget

import (
	"net/http"
	"testing"

	"devhub-backend/internal/domain/entity"
)

func TestApprovalTargetFromRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		route  string
		want   entity.ApprovalTarget
		ok     bool
	}{
		{
			name:   "create scaffold request",
			method: http.MethodPost,
			route:  "/projects/:project/scaffold-requests",
			want: entity.ApprovalTarget{
				Resource: entity.ApprovalResourceScaffoldRequest,
				Action:   entity.ApprovalActionCreate,
			},
			ok: true,
		},
		{
			name:   "create deployment",
			method: http.MethodPost,
			route:  "/services/:service/deployments",
			want: entity.ApprovalTarget{
				Resource: entity.ApprovalResourceDeployment,
				Action:   entity.ApprovalActionCreate,
			},
			ok: true,
		},
		{
			name:   "create release",
			method: http.MethodPost,
			route:  "/services/:service/releases",
			want: entity.ApprovalTarget{
				Resource: entity.ApprovalResourceRelease,
				Action:   entity.ApprovalActionCreate,
			},
			ok: true,
		},
		{
			name:   "unknown route",
			method: http.MethodGet,
			route:  "/projects/:project",
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := ApprovalTargetFromRoute(tt.method, tt.route)
			if ok != tt.ok {
				t.Fatalf("ApprovalTargetFromRoute() ok = %v, want %v", ok, tt.ok)
			}
			if !tt.ok {
				return
			}
			if got != tt.want {
				t.Fatalf("ApprovalTargetFromRoute() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
