package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPluginType = fmt.Errorf("invalid plugin type")
)

type PluginType string

const (
	PluginScaffolder PluginType = "scaffolder"
	PluginRunner     PluginType = "runner"
)

var pluginTypeStringMapper = map[PluginType]string{
	PluginScaffolder: "scaffolder",
	PluginRunner:     "runner",
}

func (s PluginType) String() string {
	return pluginTypeStringMapper[s]
}

func (s PluginType) IsValid() bool {
	switch s {
	case PluginScaffolder, PluginRunner:
		return true
	default:
		return false
	}
}

// Parse parses a string into a PluginType. It returns an error if the string is not a valid PluginType.
func (s PluginType) Parse(role string) (PluginType, error) {
	pluginType := PluginType(role)

	if !pluginType.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidPluginType, role)
	}

	return pluginType, nil
}

type Plugin struct {
	ID          uuid.UUID
	Name        string
	Version     string
	Type        PluginType
	Entrypoint  string
	Enabled     bool
	Scope       string
	Description string
	InstalledAt time.Time
}

type Plugins []Plugin
