package entity

type ScaffoldRequest struct {
	Template    string
	ProjectID   string
	Environment string
	Variables   map[string]string
}

type ScaffoldRequests []ScaffoldRequest
