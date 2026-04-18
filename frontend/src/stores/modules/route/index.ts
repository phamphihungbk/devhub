import { defineStore } from 'pinia'

import { SetupStoreId } from '@/enum'

export interface CurrentRouteState {
  name: string
  path: string
  fullPath: string
  title: string
}

export interface RouteState {
  currentRoute: CurrentRouteState
}

const defaultRouteState = (): CurrentRouteState => ({
  name: '',
  path: '',
  fullPath: '',
  title: '',
})

export const useRouteStore = defineStore(SetupStoreId.Route, {
  state: (): RouteState => ({
    currentRoute: defaultRouteState(),
  }),
  getters: {
    routeTitle(state) {
      return state.currentRoute.title
    },
  },
  actions: {
    setCurrentRoute(route: Partial<CurrentRouteState>) {
      this.currentRoute = {
        ...this.currentRoute,
        ...route,
      }
    },
    resetCurrentRoute() {
      this.currentRoute = defaultRouteState()
    },
  },
})
