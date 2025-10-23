package models

import "time"

type PluginType string

const (
	PluginScaffolder PluginType = "scaffolder"
	PluginRunner     PluginType = "runner"
)

type Plugin struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Type        PluginType `gorm:"type:varchar(16)" json:"type"`
	Description string     `json:"description,omitempty"`
	InstalledAt time.Time  `json:"installedAt"`
}
