import { defineStore } from 'pinia'
import { reactive } from 'vue';

import { canAccess, resolveProfilePermissions, type AccessRequirement } from '@/services/access/rbac'
import { SetupStoreId } from '@/enum'
import { fetchCurrentUser, loginWithPassword } from '@/services/api'
import type { Credentials, LoginTokens, UserProfile } from '@/services/api'

export const TOKEN_STORAGE_KEY = 'devhub.access_token'
export const REFRESH_STORAGE_KEY = 'devhub.refresh_token'

const initialState = () => reactive({
  accessToken: '',
  refreshToken: '',
  profile: null as UserProfile | null,
  ready: false,
})

function readStoredSession(): Pick<AuthState, 'accessToken' | 'refreshToken'> {
  return {
    accessToken: localStorage.getItem(TOKEN_STORAGE_KEY) || '',
    refreshToken: localStorage.getItem(REFRESH_STORAGE_KEY) || '',
  }
}

function persistSessionTokens(tokens: LoginTokens) {
  localStorage.setItem(TOKEN_STORAGE_KEY, tokens.accessToken)
  localStorage.setItem(REFRESH_STORAGE_KEY, tokens.refreshToken)
}

function clearStoredSession() {
  localStorage.removeItem(TOKEN_STORAGE_KEY)
  localStorage.removeItem(REFRESH_STORAGE_KEY)
}

export interface AuthState {
  accessToken: string
  refreshToken: string
  profile: UserProfile | null
  ready: boolean
}

export const useAuthStore = defineStore(SetupStoreId.Auth, {
  state: (): AuthState => initialState(),
  getters: {
    isAuthenticated(state: AuthState) {
      return Boolean(state.accessToken)
    },
    permissions(state: AuthState) {
      return resolveProfilePermissions(state.profile)
    },
    canAccess(state: AuthState) {
      return (requirement: AccessRequirement) => canAccess(state.profile, requirement)
    },
  },
  actions: {
    applySession(tokens: LoginTokens) {
      this.accessToken = tokens.accessToken
      this.refreshToken = tokens.refreshToken
      persistSessionTokens(tokens)
    },
    resetSessionState() {
      this.accessToken = ''
      this.refreshToken = ''
      this.profile = null
    },
    restoreSession() {
      const session = readStoredSession()
      this.accessToken = session.accessToken
      this.refreshToken = session.refreshToken
      this.ready = true
    },
    persistTokens(tokens: LoginTokens) {
      this.applySession(tokens)
    },
    clearSession() {
      this.resetSessionState()
      clearStoredSession()
    },
    async login(credentials: Credentials) {
      const tokens = await loginWithPassword(credentials)
      this.applySession(tokens)
      await this.loadProfile()
    },
    async loadProfile() {
      if (!this.accessToken) return null
      try {
        const profile = await fetchCurrentUser()
        this.profile = profile
        return profile
      } catch (error) {
        this.clearSession()
        throw error
      }
    },
  },
})
