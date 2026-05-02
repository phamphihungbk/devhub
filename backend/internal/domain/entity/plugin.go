package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPluginType         = fmt.Errorf("invalid plugin type")
	ErrInvalidPluginRuntime      = fmt.Errorf("invalid plugin runtime")
	ErrInvalidPluginConfigSchema = fmt.Errorf("invalid plugin config schema")
)

const (
	PluginScaffoldRequest PluginType = "scaffold_request"
	PluginDeployment      PluginType = "deployment"
)

const (
	PluginRuntimePython PluginRuntime = "python"
	PluginRuntimeGo     PluginRuntime = "go"
	PluginRuntimeNode   PluginRuntime = "node"
)

type PluginType string

type PluginRuntime string

type PluginSchemaProperty struct {
	Type        string                          `json:"type" yaml:"type"` // string, number, boolean, object
	Description string                          `json:"description,omitempty" yaml:"description,omitempty"`
	Default     interface{}                     `json:"default,omitempty" yaml:"default,omitempty"`
	Enum        []string                        `json:"enum,omitempty" yaml:"enum,omitempty"`             // allowed values
	Properties  map[string]PluginSchemaProperty `json:"properties,omitempty" yaml:"properties,omitempty"` // nested
}

type PluginConfigSchema struct {
	Type       string                          `json:"type" yaml:"type"` // usually "object"
	Required   []string                        `json:"required,omitempty" yaml:"required,omitempty"`
	Properties map[string]PluginSchemaProperty `json:"properties" yaml:"properties"`
}

func (p PluginConfigSchema) Parse(input string) (PluginConfigSchema, error) {

	var configSchema PluginConfigSchema
	err := json.Unmarshal([]byte(input), &configSchema)

	if err != nil {
		return PluginConfigSchema{}, fmt.Errorf("%w: %s", ErrInvalidPluginConfigSchema, configSchema)
	}

	return configSchema, nil
}

func (p PluginConfigSchema) String() string {
	bytes, err := json.Marshal(p)

	if err != nil {
		return ""
	}

	return string(bytes)
}

func (p PluginConfigSchema) IsZero() bool {
	return p.Type == "" && len(p.Required) == 0 && len(p.Properties) == 0
}

var pluginTypeStringMapper = map[PluginType]string{
	PluginScaffoldRequest: "scaffold_request",
	PluginDeployment:      "deployment",
}

var pluginRuntimeStringMapper = map[PluginRuntime]string{
	PluginRuntimePython: "python",
	PluginRuntimeGo:     "go",
	PluginRuntimeNode:   "node",
}

func (pt PluginType) String() string {
	return pluginTypeStringMapper[pt]
}

func (pt PluginType) IsValid() bool {
	switch pt {
	case PluginScaffoldRequest, PluginDeployment:
		return true
	default:
		return false
	}
}

// Parse parses a string into a PluginType. It returns an error if the string is not a valid PluginType.
func (pt PluginType) Parse(input string) (PluginType, error) {
	pluginType := PluginType(input)

	if !pluginType.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidPluginType, pluginType)
	}

	return pluginType, nil
}

func (s PluginRuntime) String() string {
	return pluginRuntimeStringMapper[s]
}

func (s PluginRuntime) IsValid() bool {
	switch s {
	case PluginRuntimePython, PluginRuntimeGo, PluginRuntimeNode:
		return true
	default:
		return false
	}
}

func (s PluginRuntime) Parse(runtime string) (PluginRuntime, error) {
	pluginRuntime := PluginRuntime(runtime)

	if !pluginRuntime.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidPluginRuntime, runtime)
	}

	return pluginRuntime, nil
}

type Plugin struct {
	ID           uuid.UUID
	Name         string
	Version      string
	Type         PluginType
	Runtime      PluginRuntime
	Entrypoint   string
	ConfigSchema PluginConfigSchema
	Description  string
	Enabled      bool
	CreatedAt    time.Time
}

type Plugins []Plugin
