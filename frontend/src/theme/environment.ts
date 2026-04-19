export type EnvironmentTagColor = {
  color: string
  textColor: string
}

export const environmentOptions = [
  { label: 'Development', value: 'dev' },
  { label: 'Staging', value: 'staging' },
  { label: 'Production', value: 'prod' },
] as const

export function getEnvironmentTagColor(environment: string): EnvironmentTagColor {
  switch (environment) {
    case 'dev':
      return { color: '#dbeafe', textColor: '#1d4ed8' }
    case 'staging':
      return { color: '#fef3c7', textColor: '#b45309' }
    case 'prod':
      return { color: '#dcfce7', textColor: '#15803d' }
    default:
      return { color: '#e2e8f0', textColor: '#334155' }
  }
}
