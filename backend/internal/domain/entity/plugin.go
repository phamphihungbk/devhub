package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPluginType    = fmt.Errorf("invalid plugin type")
	ErrInvalidPluginScope   = fmt.Errorf("invalid plugin scope")
	ErrInvalidPluginRuntime = fmt.Errorf("invalid plugin runtime")
)

type PluginType string
type PluginScope string
type PluginRuntime string

const (
	PluginScaffolder PluginType = "scaffolder"
	PluginDeployer   PluginType = "deployer"
	PluginReleaser   PluginType = "releaser"
)

const (
	PluginScopeGlobal      PluginScope = "global"
	PluginScopeProject     PluginScope = "project"
	PluginScopeEnvironment PluginScope = "environment"
)

const (
	PluginRuntimePython PluginRuntime = "python"
	PluginRuntimeGo     PluginRuntime = "go"
	PluginRuntimeNode   PluginRuntime = "node"
)

var pluginTypeStringMapper = map[PluginType]string{
	PluginScaffolder: "scaffolder",
	PluginDeployer:   "deployer",
	PluginReleaser:   "releaser",
}

var pluginScopeStringMapper = map[PluginScope]string{
	PluginScopeGlobal:      "global",
	PluginScopeProject:     "project",
	PluginScopeEnvironment: "environment",
}

var pluginRuntimeStringMapper = map[PluginRuntime]string{
	PluginRuntimePython: "python",
	PluginRuntimeGo:     "go",
	PluginRuntimeNode:   "node",
}

func (s PluginType) String() string {
	return pluginTypeStringMapper[s]
}

func (s PluginType) IsValid() bool {
	switch s {
	case PluginScaffolder, PluginDeployer, PluginReleaser:
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

func (s PluginScope) String() string {
	return pluginScopeStringMapper[s]
}

func (s PluginScope) IsValid() bool {
	switch s {
	case PluginScopeGlobal, PluginScopeProject, PluginScopeEnvironment:
		return true
	default:
		return false
	}
}

func (s PluginScope) Parse(scope string) (PluginScope, error) {
	pluginScope := PluginScope(scope)

	if !pluginScope.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidPluginScope, scope)
	}

	return pluginScope, nil
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
	ID          uuid.UUID
	Name        string
	Version     string
	Type        PluginType
	Runtime     PluginRuntime
	Entrypoint  string
	Enabled     bool
	Scope       PluginScope
	Description string
	InstalledAt time.Time
}

type Plugins []Plugin
