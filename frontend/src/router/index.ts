import type { Component } from 'vue'
import {
  Dashboard,
  Document,
  Plug,
  UserAvatar,
} from '@vicons/carbon'
import { h } from 'vue'
import {
  createRouter,
  createWebHistory,
  type NavigationGuardNext,
  type RouteLocationNormalized,
  type RouteRecordRaw,
} from 'vue-router'

import { TOKEN_STORAGE_KEY } from '@/stores/modules/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard',
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: 'Sign in',
      guestOnly: true,
    },
  },
  {
    path: '/',
    component: () => import('@/layouts/admin-shell.vue'),
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: 'Dashboard',
          icon: Dashboard,
        },
      },
      {
        path: 'projects',
        name: 'projects',
        component: () => import('@/views/projects/index.vue'),
        meta: {
          title: 'Project List',
          icon: Document,
        },
      },
      {
        path: 'projects/new',
        name: 'project-create',
        component: () => import('@/views/projects/create.vue'),
        meta: {
          title: 'Create Project',
          icon: Document,
        },
      },
      {
        path: 'plugins',
        name: 'plugins',
        component: () => import('@/views/plugins/index.vue'),
        meta: {
          title: 'Plugins',
          icon: Plug,
        },
      },
      {
        path: 'users',
        name: 'users',
        component: () => import('@/views/users/index.vue'),
        meta: {
          title: 'Users',
          icon: UserAvatar,
        },
      },
    ],
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const hasToken = Boolean(localStorage.getItem(TOKEN_STORAGE_KEY))

  if (to.meta.requiresAuth && !hasToken) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.guestOnly && hasToken) {
    next({ name: 'dashboard' })
    return
  }

  next()
})

export function renderRouteIcon(icon?: Component) {
  if (!icon) return undefined
  return () => h(icon, { size: 18 })
}
