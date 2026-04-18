export const apiBaseURL = {
  auth: '/auth',
  users: '/users',
  projects: '/projects',
  plugins: '/plugins',
  deployments: '/deployments',
  scaffoldRequests: '/scaffold-requests',
} as const

export * from './auth'
export * from './deployments'
export * from './plugins'
export * from './projects'
export * from './releases'
export * from './scaffold-requests'
export * from './users'
