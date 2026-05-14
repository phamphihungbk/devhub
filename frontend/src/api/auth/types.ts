export interface Credentials {
  email: string
  password: string
}

export interface LoginTokens {
  accessToken: string
  refreshToken: string
}

export interface UserProfile {
  id: string
  name: string
  email: string
  role?: string
  roles?: string[]
  team_id: string
  permissions?: string[]
}
