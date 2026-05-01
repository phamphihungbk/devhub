package cmd

import (
	"context"
	"fmt"

	"devhub-backend/internal/config"
	"devhub-backend/internal/domain/entity"
	infraDB "devhub-backend/internal/infra/db"
	"devhub-backend/internal/util/misc"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed-data",
	Short: "Seed bootstrap teams, RBAC data, and approval policies",
	Long: `Seed RBAC bootstrap data into the database.

This command creates or refreshes the canonical:
- teams
- roles
- permissions
- role_permissions
- approval_policies

It is safe to rerun after migrations because inserts are idempotent.`,
	RunE: runSeedCmd,
}

type seedRole struct {
	Name        string
	Description string
}

type seedPermission struct {
	Name        string
	Description string
}

type seedApprovalPolicy struct {
	Resource          string
	Action            string
	Environment       *string
	RequiredApprovals int
	Enabled           bool
}

type seedTeam struct {
	Name         string
	OwnerContact string
}

type seedUser struct {
	Name     string
	Email    string
	Password string
	Role     string
	TeamName string
}

var seedRoles = []seedRole{
	{Name: entity.RolePlatformAdmin.String(), Description: "Full platform administration access"},
	{Name: entity.RoleOrgAdmin.String(), Description: "Organization-wide administration access"},
	{Name: entity.RoleTeamLead.String(), Description: "Team-level control plane management"},
	{Name: entity.RoleDeveloper.String(), Description: "Developer access for scaffold, release, and deploy actions"},
	{Name: entity.RoleViewer.String(), Description: "Read-only access"},
}

var seedPermissions = []seedPermission{
	{Name: entity.PermissionUserRead.String(), Description: "Read user records"},
	{Name: entity.PermissionUserWrite.String(), Description: "Create, update, and delete users"},
	{Name: entity.PermissionProjectWrite.String(), Description: "Create, update, and delete projects"},
	{Name: entity.PermissionScaffoldRequestWrite.String(), Description: "Create and delete scaffold requests"},
	{Name: entity.PermissionReleaseWrite.String(), Description: "Create releases"},
	{Name: entity.PermissionDeploymentWrite.String(), Description: "Create, update, and delete deployments"},
	{Name: entity.PermissionPluginWrite.String(), Description: "Create, update, and delete plugins"},
}

var seedRolePermissions = map[string][]string{
	entity.RolePlatformAdmin.String(): {
		entity.PermissionUserRead.String(),
		entity.PermissionUserWrite.String(),
		entity.PermissionProjectWrite.String(),
		entity.PermissionScaffoldRequestWrite.String(),
		entity.PermissionReleaseWrite.String(),
		entity.PermissionDeploymentWrite.String(),
		entity.PermissionPluginWrite.String(),
	},
	entity.RoleOrgAdmin.String(): {
		entity.PermissionUserRead.String(),
		entity.PermissionUserWrite.String(),
		entity.PermissionProjectWrite.String(),
		entity.PermissionScaffoldRequestWrite.String(),
		entity.PermissionReleaseWrite.String(),
		entity.PermissionDeploymentWrite.String(),
	},
	entity.RoleTeamLead.String(): {
		entity.PermissionProjectWrite.String(),
		entity.PermissionScaffoldRequestWrite.String(),
		entity.PermissionReleaseWrite.String(),
		entity.PermissionDeploymentWrite.String(),
	},
	entity.RoleDeveloper.String(): {
		entity.PermissionScaffoldRequestWrite.String(),
		entity.PermissionReleaseWrite.String(),
		entity.PermissionDeploymentWrite.String(),
	},
	entity.RoleViewer.String(): {},
}

var seedApprovalPolicies = []seedApprovalPolicy{
	{
		Resource:          entity.ApprovalResourceScaffoldRequest.String(),
		Action:            entity.ApprovalActionCreate.String(),
		RequiredApprovals: 1,
		Enabled:           true,
	},
	{
		Resource:          entity.ApprovalResourceRelease.String(),
		Action:            entity.ApprovalActionCreate.String(),
		Environment:       misc.ToPointer("development"),
		RequiredApprovals: 1,
		Enabled:           true,
	},
	{
		Resource:          entity.ApprovalResourceDeployment.String(),
		Action:            entity.ApprovalActionCreate.String(),
		Environment:       misc.ToPointer("development"),
		RequiredApprovals: 1,
		Enabled:           true,
	},
}

var seedTeams = []seedTeam{
	{
		Name:         "phamphihungbk",
		OwnerContact: "phamphihungbk",
	},
	{
		Name:         "devhub",
		OwnerContact: "devhub",
	},
}

var seedUsers = []seedUser{
	{
		Name:     "DevHub Admin",
		Email:    "admin@devhub.local",
		Password: "admindevhub123",
		Role:     entity.RolePlatformAdmin.String(),
		TeamName: "phamphihungbk",
	},
	{
		Name:     "DevHub Org Admin",
		Email:    "org-admin@devhub.local",
		Password: "orgadmindevhub123",
		Role:     entity.RoleOrgAdmin.String(),
		TeamName: "phamphihungbk",
	},
	{
		Name:     "DevHub Team Lead",
		Email:    "team-lead@devhub.local",
		Password: "teamleaddevhub123",
		Role:     entity.RoleTeamLead.String(),
		TeamName: "phamphihungbk",
	},
	{
		Name:     "DevHub Developer",
		Email:    "developer@devhub.local",
		Password: "developerdevhub123",
		Role:     entity.RoleDeveloper.String(),
		TeamName: "phamphihungbk",
	},
	{
		Name:     "DevHub Viewer",
		Email:    "viewer@devhub.local",
		Password: "viewerdevhub123",
		Role:     entity.RoleViewer.String(),
		TeamName: "phamphihungbk",
	},
}

func runSeedCmd(cmd *cobra.Command, args []string) error {
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

	for _, team := range seedTeams {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO teams (name, owner_contact)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE
			SET owner_contact = EXCLUDED.owner_contact,
				updated_at = now()
		`, team.Name, team.OwnerContact); err != nil {
			return fmt.Errorf("seed team %q: %w", team.Name, err)
		}
	}

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

	for _, user := range seedUsers {
		passwordHash, err := misc.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("hash password for %q: %w", user.Email, err)
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO users (name, email, password_hash, team_id)
			SELECT $1, $2, $3, t.id
			FROM teams t
			WHERE t.name = $4
			ON CONFLICT (email) DO UPDATE
			SET name = EXCLUDED.name,
				password_hash = EXCLUDED.password_hash,
				team_id = EXCLUDED.team_id,
				updated_at = now()
		`, user.Name, user.Email, passwordHash, user.TeamName); err != nil {
			return fmt.Errorf("seed user %q: %w", user.Email, err)
		}

		if _, err := tx.ExecContext(ctx, `
			DELETE FROM user_roles ur
			USING users u
			WHERE ur.user_id = u.id
				AND u.email = $1
		`, user.Email); err != nil {
			return fmt.Errorf("clear user roles for %q: %w", user.Email, err)
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_roles (user_id, role_id)
			SELECT u.id, r.id
			FROM users u
			JOIN roles r ON r.name = $2
			WHERE u.email = $1
			ON CONFLICT (user_id, role_id) DO NOTHING
		`, user.Email, user.Role); err != nil {
			return fmt.Errorf("seed user role %q -> %q: %w", user.Email, user.Role, err)
		}
	}

	for _, policy := range seedApprovalPolicies {
		if policy.Environment == nil {
			continue
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE approval_policies
			SET required_approvals = $4::int,
				enabled = $5::boolean,
				updated_at = now()
			WHERE resource = $1::varchar
				AND action = $2::varchar
				AND project_id IS NULL
				AND service_id IS NULL
				AND environment_id IN (
					SELECT id
					FROM environments
					WHERE name = $3::varchar
				)
		`, policy.Resource, policy.Action, policy.Environment, policy.RequiredApprovals, policy.Enabled); err != nil {
			return fmt.Errorf("update approval policy %q/%q: %w", policy.Resource, policy.Action, err)
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO approval_policies (
				resource,
				action,
				project_id,
				service_id,
				environment_id,
				required_approvals,
				enabled
			)
			SELECT $1::varchar, $2::varchar, NULL, NULL, e.id, $4::int, $5::boolean
			FROM environments e
			WHERE e.name = $3::varchar
				AND NOT EXISTS (
					SELECT 1
					FROM approval_policies
					WHERE resource = $1::varchar
						AND action = $2::varchar
						AND project_id IS NULL
						AND service_id IS NULL
						AND environment_id = e.id
				)
		`, policy.Resource, policy.Action, policy.Environment, policy.RequiredApprovals, policy.Enabled); err != nil {
			return fmt.Errorf("seed approval policy %q/%q: %w", policy.Resource, policy.Action, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	tx = nil

	fmt.Fprintln(cmd.OutOrStdout(), "seed complete")
	return nil
}
