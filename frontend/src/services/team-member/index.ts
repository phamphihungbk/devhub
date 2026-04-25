import { h, onMounted, ref } from 'vue'
import { NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { fetchUsers } from '@/api'
import { getRoleTagColor } from '@/theme/role'
import type { UserRecord } from '@/api'
import { ApiError } from '@/api/request'

export function useTeamMemberService() {
  const message = useMessage()
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
          NTag,
          { bordered: false, color: getRoleTagColor(row.role) },
          { default: () => row.role },
        ),
    },
  ]

  const loadUsers = async() => {
    loading.value = true
    try {
      rows.value = await fetchUsers()
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
  }
}
