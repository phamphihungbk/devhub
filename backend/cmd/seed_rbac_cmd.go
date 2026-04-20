package cmd

import (
	"context"
	"fmt"

	"devhub-backend/internal/config"
	infraDB "devhub-backend/internal/infra/db"

	"github.com/spf13/cobra"
)

var seedRBACCmd = &cobra.Command{
	Use:   "seed-rbac",
	Short: "Seed RBAC roles, permissions, and role-permission mappings",
	Long: `Seed RBAC bootstrap data into the database.

This command creates or refreshes the canonical:
- roles
- permissions
- role_permissions

It is safe to rerun after migrations because inserts are idempotent.`,
	RunE: runSeedRBACCmd,
}

type seedRole struct {
	Name        string
	Description string
}

type seedPermission struct {
	Name        string
	Description string
}

var seedRoles = []seedRole{
	{Name: "platform_admin", Description: "Full platform administration access"},
	{Name: "org_admin", Description: "Organization-wide administration access"},
	{Name: "team_lead", Description: "Team-level control plane management"},
	{Name: "developer", Description: "Developer access for scaffold, release, and deploy actions"},
	{Name: "viewer", Description: "Read-only access"},
}

var seedPermissions = []seedPermission{
	{Name: "user.read", Description: "Read user records"},
	{Name: "user.write", Description: "Create, update, and delete users"},
	{Name: "project.write", Description: "Create, update, and delete projects"},
	{Name: "scaffold_request.write", Description: "Create and delete scaffold requests"},
	{Name: "release.write", Description: "Create releases"},
	{Name: "deployment.write", Description: "Create, update, and delete deployments"},
	{Name: "plugin.write", Description: "Create, update, and delete plugins"},
}

var seedRolePermissions = map[string][]string{
	"platform_admin": {
		"user.read",
		"user.write",
		"project.write",
		"scaffold_request.write",
		"release.write",
		"deployment.write",
		"plugin.write",
	},
	"org_admin": {
		"user.read",
		"user.write",
		"project.write",
		"scaffold_request.write",
		"release.write",
		"deployment.write",
	},
	"team_lead": {
		"project.write",
		"scaffold_request.write",
		"release.write",
		"deployment.write",
	},
	"developer": {
		"scaffold_request.write",
		"release.write",
		"deployment.write",
	},
	"viewer": {},
}

func runSeedRBACCmd(cmd *cobra.Command, args []string) error {
	cfg := config.MustConfigure()
	dbConn := infraDB.MustConnect(cfg)
	defer dbConn.Close()

	ctx := context.Background()
	tx, err := dbConn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	for _, role := range seedRoles {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO roles (name, description)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description
		`, role.Name, role.Description); err != nil {
			return fmt.Errorf("seed role %q: %w", role.Name, err)
		}
	}

	for _, permission := range seedPermissions {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO permissions (name, description)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description
		`, permission.Name, permission.Description); err != nil {
			return fmt.Errorf("seed permission %q: %w", permission.Name, err)
		}
	}

	for roleName, permissions := range seedRolePermissions {
		for _, permissionName := range permissions {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO role_permissions (role_id, permission_id)
				SELECT r.id, p.id
				FROM roles r
				JOIN permissions p ON p.name = $2
				WHERE r.name = $1
				ON CONFLICT (role_id, permission_id) DO NOTHING
			`, roleName, permissionName); err != nil {
				return fmt.Errorf("seed role permission %q -> %q: %w", roleName, permissionName, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	tx = nil

	fmt.Fprintln(cmd.OutOrStdout(), "rbac seed complete")
	return nil
}
