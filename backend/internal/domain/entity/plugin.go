package entity

import (
	"time"

	"github.com/google/uuid"
)

type PluginType string

const (
	PluginScaffolder PluginType = "scaffolder"
	PluginRunner     PluginType = "runner"
)

type Plugin struct {
	ID          uuid.UUID
	Name        string
	Version     string
	Description string
	Type        PluginType
	InstalledAt time.Time
}

type Plugins []Plugin
