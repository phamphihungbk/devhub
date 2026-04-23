export type RoleTagColor = {
  color: string
  textColor: string
}

export function getRoleTagColor(role: string): RoleTagColor {
  switch (role) {
    case 'platform_admin':
      return { color: '#fee2e2', textColor: '#b91c1c' }
    case 'org_admin':
      return { color: '#ffedd5', textColor: '#c2410c' }
    case 'team_lead':
      return { color: '#fef3c7', textColor: '#b45309' }
    case 'developer':
      return { color: '#dcfce7', textColor: '#15803d' }
    case 'viewer':
      return { color: '#e0f2fe', textColor: '#0369a1' }
    default:
      return { color: '#e2e8f0', textColor: '#334155' }
  }
}
