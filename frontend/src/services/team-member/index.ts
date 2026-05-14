import { h, onMounted, ref } from 'vue'
import { NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { fetchUsers } from '@/api'
import { useAuthStore } from '@/stores/modules/auth'
import { getRoleTagColor } from '@/theme/role'
import type { UserRecord } from '@/api'
import { ApiError } from '@/api/request'

function userRoles(row: UserRecord) {
  return row.roles?.length ? row.roles : row.role ? [row.role] : []
}

export function useTeamMemberService() {
  const message = useMessage()
  const authStore = useAuthStore()
  const loading = ref(false)
  const rows = ref<UserRecord[]>([])

  const columns: DataTableColumns<UserRecord> = [
    { title: 'Name', key: 'name' },
    { title: 'Email', key: 'email' },
    {
      title: 'Role',
      key: 'role',
      render: row =>
        h(
          'div',
          { class: 'flex flex-wrap gap-2' },
          userRoles(row).map(role =>
            h(
              NTag,
              { bordered: false, color: getRoleTagColor(role) },
              { default: () => role },
            ),
          ),
        ),
    },
  ]

  const loadUsers = async() => {
    loading.value = true
    try {
      rows.value = await fetchUsers({
        team_id: authStore.profile?.team_id || undefined,
      })
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to load users.')
    } finally {
      loading.value = false
    }
  }

  onMounted(loadUsers)

  return {
    columns,
    loadUsers,
    loading,
    rows,
    userRoles,
  }
}
