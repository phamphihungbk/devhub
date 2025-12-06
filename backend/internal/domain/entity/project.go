package entity

import "gorm.io/gorm"

type ProjectEnvironment string

const (
	EnvDev  ProjectEnvironment = "dev"
	EnvProd ProjectEnvironment = "prod"
)

type Project struct {
	ID           string               `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description,omitempty"`
	Environments []ProjectEnvironment `gorm:"type:text[]" json:"environments"`
	CreatedBy    string               `json:"createdBy"`
	DeletedAt    gorm.DeletedAt       `gorm:"index" json:"-"`
}

type Projects []Project
