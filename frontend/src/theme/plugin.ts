export type PluginTypeTagColor = {
  color: string
  textColor: string
}

export const pluginTypeOptions = [
  { label: 'Deployer', value: 'deployer' },
  { label: 'Releaser', value: 'releaser' },
  { label: 'Scaffolder', value: 'scaffolder' },
] as const

export function getPluginTypeTagColor(type: string): PluginTypeTagColor {
  switch (type) {
    case 'deployer':
      return { color: '#ede9fe', textColor: '#6d28d9' }
    case 'releaser':
      return { color: '#ffe4e6', textColor: '#be123c' }
    case 'scaffolder':
      return { color: '#ccfbf1', textColor: '#0f766e' }
    default:
      return { color: '#e2e8f0', textColor: '#334155' }
  }
}
