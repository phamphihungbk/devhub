import type { UserProfile } from '@/services/api'

export type AppRole =
  | 'platform_admin'
  | 'org_admin'
  | 'team_lead'
  | 'developer'
  | 'viewer'

export type AppPermission =
  | 'user.read'
  | 'user.write'
  | 'project.write'
  | 'scaffold_request.write'
  | 'release.write'
  | 'deployment.write'
  | 'plugin.write'

export type AccessRequirement = {
  roles?: string[]
  permissions?: string[]
}

export const permission = {
  userRead: 'user.read',
  userWrite: 'user.write',
  projectWrite: 'project.write',
  scaffoldRequestWrite: 'scaffold_request.write',
  releaseWrite: 'release.write',
  deploymentWrite: 'deployment.write',
  pluginWrite: 'plugin.write',
} as const satisfies Record<string, AppPermission>

const rolePermissions: Record<AppRole, AppPermission[]> = {
  platform_admin: [
    permission.userRead,
    permission.userWrite,
    permission.projectWrite,
    permission.scaffoldRequestWrite,
    permission.releaseWrite,
    permission.deploymentWrite,
    permission.pluginWrite,
  ],
  org_admin: [
    permission.userRead,
    permission.userWrite,
    permission.projectWrite,
    permission.scaffoldRequestWrite,
    permission.releaseWrite,
    permission.deploymentWrite,
  ],
  team_lead: [
    permission.projectWrite,
    permission.scaffoldRequestWrite,
    permission.releaseWrite,
    permission.deploymentWrite,
  ],
  developer: [
    permission.scaffoldRequestWrite,
    permission.releaseWrite,
    permission.deploymentWrite,
  ],
  viewer: [],
}

function normalizeStringList(value: unknown): string[] {
  if (!Array.isArray(value)) {
    return []
  }

  return value.filter((item): item is string => typeof item === 'string')
}

export function getPermissionsForRole(role?: string | null): AppPermission[] {
  if (!role || !(role in rolePermissions)) {
    return []
  }

  return rolePermissions[role as AppRole]
}

export function resolveProfilePermissions(profile?: UserProfile | null): string[] {
  const explicitPermissions = normalizeStringList(profile?.permissions)

  if (explicitPermissions.length > 0) {
    return explicitPermissions
  }

  return getPermissionsForRole(profile?.role)
}

export function hasRole(profile: UserProfile | null | undefined, roles: string[] = []): boolean {
  if (roles.length === 0) {
    return true
  }

  return Boolean(profile?.role && roles.includes(profile.role))
}

export function hasPermissions(profile: UserProfile | null | undefined, permissions: string[] = []): boolean {
  if (permissions.length === 0) {
    return true
  }

  const resolvedPermissions = new Set(resolveProfilePermissions(profile))

  return permissions.every(item => resolvedPermissions.has(item))
}

export function canAccess(profile: UserProfile | null | undefined, requirement: AccessRequirement = {}): boolean {
  return hasRole(profile, requirement.roles) && hasPermissions(profile, requirement.permissions)
}

export function accessRequirementFromMeta(meta: Record<string, unknown> | undefined): AccessRequirement {
  return {
    roles: normalizeStringList(meta?.roles),
    permissions: normalizeStringList(meta?.permissions),
  }
}

export function canAccessMeta(profile: UserProfile | null | undefined, meta: Record<string, unknown> | undefined): boolean {
  return canAccess(profile, accessRequirementFromMeta(meta))
}
