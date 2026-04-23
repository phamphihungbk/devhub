import type { Component } from 'vue'
import {
  Dashboard,
  Document,
  Plug,
  UserAvatar,
} from '@vicons/carbon'
import { h } from 'vue'
import {
  canAccessMeta,
  permission,
} from '@/access/rbac'
import {
  createRouter,
  createWebHistory,
  type NavigationGuardNext,
  type RouteLocationNormalized,
  type RouteRecordRaw,
} from 'vue-router'

import { TOKEN_STORAGE_KEY } from '@/stores/modules/auth'
import { useAuthStore } from '@/stores/modules/auth'

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
        path: 'approvals',
        name: 'approvals',
        component: () => import('@/views/approvals/index.vue'),
        meta: {
          title: 'Approvals',
          icon: Document,
          permissions: [permission.scaffoldRequestWrite],
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
        path: 'services',
        name: 'services',
        component: () => import('@/views/services/index.vue'),
        meta: {
          title: 'Services',
          icon: Document,
        },
      },
      {
        path: 'releases',
        name: 'releases',
        component: () => import('@/views/releases/index.vue'),
        meta: {
          title: 'Releases',
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
          permissions: [permission.projectWrite],
        },
      },
      {
        path: 'projects/:projectId',
        name: 'project-details',
        component: () => import('@/views/projects/detail.vue'),
        meta: {
          title: 'Project Details',
          icon: Document,
          hideInMenu: true,
          activeMenu: 'projects',
        },
      },
      {
        path: 'projects/:projectId/services/:serviceId',
        name: 'service-details',
        component: () => import('@/views/services/detail.vue'),
        meta: {
          title: 'Service Details',
          icon: Document,
          hideInMenu: true,
          activeMenu: 'projects',
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
          permissions: [permission.userRead],
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

router.beforeEach(async (to: RouteLocationNormalized, _from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const authStore = useAuthStore()
  const hasToken = Boolean(authStore.accessToken || localStorage.getItem(TOKEN_STORAGE_KEY))

  if (to.meta.requiresAuth && !hasToken) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.guestOnly && hasToken) {
    next({ name: 'dashboard' })
    return
  }

  if (to.meta.requiresAuth && hasToken && !authStore.profile) {
    try {
      await authStore.loadProfile()
    } catch {
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }
  }

  if (to.meta.requiresAuth && !canAccessMeta(authStore.profile, to.meta as Record<string, unknown>)) {
    next({ name: 'dashboard' })
    return
  }

  next()
})

export function renderRouteIcon(icon?: Component) {
  if (!icon) return undefined
  return () => h(icon, { size: 18 })
}
