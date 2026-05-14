package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"devhub-backend/internal/domain/entity"
	"devhub-backend/internal/domain/repository"

	"gopkg.in/yaml.v3"
)

type SyncRegistryInput struct {
	PluginsDir string `json:"plugins_dir" validate:"required"`
}

type SyncRegistryOutput struct {
	Discovered int `json:"discovered"`
	Created    int `json:"created"`
	Updated    int `json:"updated"`
}

type pluginManifest struct {
	Name         string                     `yaml:"name"`
	Type         string                     `yaml:"type"`
	Version      string                     `yaml:"version"`
	Entrypoint   string                     `yaml:"entrypoint"`
	Description  string                     `yaml:"description"`
	Enabled      *bool                      `yaml:"enabled"`
	Runtime      string                     `yaml:"runtime"`
	ConfigSchema *entity.PluginConfigSchema `yaml:"config_schema"`
}

type discoveredPlugin struct {
	entity.Plugin
	hasConfigSchema bool
}

func (u *pluginUsecase) SyncRegistry(ctx context.Context, input SyncRegistryInput) (SyncRegistryOutput, error) {
	discovered, err := u.discoverPlugins(input.PluginsDir)
	if err != nil {
		return SyncRegistryOutput{}, err
	}

	existing, err := u.loadExistingPlugins(ctx)
	if err != nil {
		return SyncRegistryOutput{}, err
	}

	result := SyncRegistryOutput{Discovered: len(discovered)}
	for _, plugin := range discovered {
		matched := u.matchExistingPlugin(existing, plugin)
		if matched == nil {
			if _, err := u.pluginRepository.CreateOne(ctx, &plugin.Plugin); err != nil {
				return SyncRegistryOutput{}, fmt.Errorf("create plugin %q: %w", plugin.Name, err)
			}
			result.Created++
			continue
		}

		updateInput := u.buildRepositoryUpdatePluginInput(matched, plugin)
		if !hasPluginChanges(updateInput) {
			continue
		}

		if _, err := u.pluginRepository.UpdateOne(ctx, updateInput); err != nil {
			return SyncRegistryOutput{}, fmt.Errorf("update plugin %q: %w", plugin.Name, err)
		}
		result.Updated++
	}

	return result, nil
}

func (u *pluginUsecase) discoverPlugins(pluginsDir string) ([]discoveredPlugin, error) {
	manifestPattern := filepath.Join(pluginsDir, "*", "*", "plugin.yaml")
	manifestPaths, err := filepath.Glob(manifestPattern)
	if err != nil {
		return nil, fmt.Errorf("glob plugin manifests: %w", err)
	}

	discovered := make([]discoveredPlugin, 0, len(manifestPaths))
	for _, manifestPath := range manifestPaths {
		plugin, err := u.readPluginManifest(manifestPath)
		if err != nil {
			return nil, err
		}
		discovered = append(discovered, plugin)
	}

	return discovered, nil
}

func (u *pluginUsecase) readPluginManifest(manifestPath string) (discoveredPlugin, error) {
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		return discoveredPlugin{}, fmt.Errorf("read manifest %s: %w", manifestPath, err)
	}

	var manifest pluginManifest
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		return discoveredPlugin{}, fmt.Errorf("parse manifest %s: %w", manifestPath, err)
	}

	pluginType, err := u.parsePluginType(manifestPath, manifest.Type)
	if err != nil {
		return discoveredPlugin{}, err
	}
	pluginRuntime, err := u.parsePluginRuntime(manifestPath, manifest.Runtime)
	if err != nil {
		return discoveredPlugin{}, err
	}

	entrypoint := buildPluginEntrypoint(manifestPath, manifest.Entrypoint)
	if entrypoint == "" {
		return discoveredPlugin{}, fmt.Errorf("manifest %s missing entrypoint", manifestPath)
	}

	enabled := true
	if manifest.Enabled != nil {
		enabled = *manifest.Enabled
	}

	configSchema, hasConfigSchema, err := u.readPluginConfigSchema(manifestPath, manifest.ConfigSchema)
	if err != nil {
		return discoveredPlugin{}, err
	}

	plugin := discoveredPlugin{
		Plugin: entity.Plugin{
			Name:        strings.TrimSpace(manifest.Name),
			Version:     strings.TrimSpace(manifest.Version),
			Type:        pluginType,
			Runtime:     pluginRuntime,
			Entrypoint:  entrypoint,
			Enabled:     enabled,
			Description: strings.TrimSpace(manifest.Description),
		},
		hasConfigSchema: hasConfigSchema,
	}
	if hasConfigSchema {
		plugin.ConfigSchema = configSchema
	}

	return plugin, nil
}

func (u *pluginUsecase) readPluginConfigSchema(manifestPath string, inline *entity.PluginConfigSchema) (entity.PluginConfigSchema, bool, error) {
	if inline != nil {
		return *inline, true, nil
	}

	schemaPath := filepath.Join(filepath.Dir(manifestPath), "schema.json")
	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return entity.PluginConfigSchema{}, false, nil
		}
		return entity.PluginConfigSchema{}, false, fmt.Errorf("read config schema %s: %w", schemaPath, err)
	}

	configSchema, err := new(entity.PluginConfigSchema).Parse(string(raw))
	if err != nil {
		return entity.PluginConfigSchema{}, false, fmt.Errorf("parse config schema %s: %w", schemaPath, err)
	}

	return configSchema, true, nil
}

func (u *pluginUsecase) parsePluginType(manifestPath string, rawType string) (entity.PluginType, error) {
	fmt.Println(rawType, "hungdeptrai")
	pluginType, err := new(entity.PluginType).Parse(strings.TrimSpace(rawType))
	if err != nil {
		return "", fmt.Errorf("manifest %s has invalid type: %w", manifestPath, err)
	}

	return pluginType, nil
}

func (u *pluginUsecase) parsePluginRuntime(manifestPath string, rawRuntime string) (entity.PluginRuntime, error) {
	pluginRuntime, err := new(entity.PluginRuntime).Parse(strings.TrimSpace(rawRuntime))
	if err != nil {
		return "", fmt.Errorf("manifest %s has invalid runtime: %w", manifestPath, err)
	}

	return pluginRuntime, nil
}

func buildPluginEntrypoint(manifestPath string, entrypoint string) string {
	entrypoint = strings.TrimSpace(entrypoint)
	if entrypoint == "" {
		return ""
	}

	if strings.HasPrefix(entrypoint, "/") {
		return entrypoint
	}

	pluginDir := filepath.Dir(manifestPath)
	pluginDir = filepath.ToSlash(pluginDir)

	idx := strings.Index(pluginDir, "plugins/")
	if idx == -1 {
		return entrypoint
	}

	return "/app/" + pluginDir[idx:] + "/" + entrypoint
}

func (u *pluginUsecase) loadExistingPlugins(ctx context.Context) (entity.Plugins, error) {
	limit := int64(1000)
	offset := int64(0)
	existing, _, err := u.pluginRepository.FindAll(ctx, repository.FindAllPluginsFilter{
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		return nil, fmt.Errorf("load existing plugins: %w", err)
	}
	if existing == nil {
		return entity.Plugins{}, nil
	}
	return *existing, nil
}

func (u *pluginUsecase) matchExistingPlugin(existing entity.Plugins, plugin discoveredPlugin) *entity.Plugin {
	for i := range existing {
		if strings.TrimSpace(existing[i].Entrypoint) == strings.TrimSpace(plugin.Entrypoint) {
			return &existing[i]
		}
	}

	for i := range existing {
		if strings.TrimSpace(existing[i].Name) == strings.TrimSpace(plugin.Name) &&
			existing[i].Type == plugin.Type {
			return &existing[i]
		}
	}

	return nil
}

func (u *pluginUsecase) buildRepositoryUpdatePluginInput(existing *entity.Plugin, plugin discoveredPlugin) repository.UpdatePluginInput {
	input := repository.UpdatePluginInput{
		ID: existing.ID,
	}

	if existing.Name != plugin.Name {
		input.Name = &plugin.Name
	}
	if existing.Description != plugin.Description {
		input.Description = &plugin.Description
	}
	if existing.Type != plugin.Type {
		pluginType := plugin.Type
		input.Type = &pluginType
	}
	if existing.Version != plugin.Version {
		input.Version = &plugin.Version
	}
	if existing.Runtime != plugin.Runtime {
		pluginRuntime := plugin.Runtime
		input.Runtime = &pluginRuntime
	}
	if existing.Entrypoint != plugin.Entrypoint {
		input.Entrypoint = &plugin.Entrypoint
	}
	if plugin.hasConfigSchema && !reflect.DeepEqual(existing.ConfigSchema, plugin.ConfigSchema) {
		configSchema := plugin.ConfigSchema
		input.ConfigSchema = &configSchema
	}
	if existing.Enabled != plugin.Enabled {
		enabled := plugin.Enabled
		input.Enabled = &enabled
	}

	return input
}

func hasPluginChanges(input repository.UpdatePluginInput) bool {
	return input.Name != nil ||
		input.Description != nil ||
		input.Type != nil ||
		input.Version != nil ||
		input.Runtime != nil ||
		input.Entrypoint != nil ||
		input.ConfigSchema != nil ||
		input.Enabled != nil
}
