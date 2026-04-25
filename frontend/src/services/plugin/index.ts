import { computed, h, onMounted, reactive, ref } from 'vue'
import { NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { fetchPlugins } from '@/api'
import { getPluginTypeTagColor, pluginTypeOptions } from '@/theme/plugin'
import type { PluginRecord } from '@/api'
import { ApiError } from '@/api/request'

export function usePluginService() {
  const message = useMessage()
  const loading = ref(false)
  const rows = ref<PluginRecord[]>([])
  const filters = reactive({
    keyword: '',
    type: null as string | null,
    runtime: null as string | null,
    scope: null as string | null,
  })

  const columns: DataTableColumns<PluginRecord> = [
    { title: 'Name', key: 'name' },
    {
      title: 'Type',
      key: 'type',
      render: row =>
        h(
          NTag,
          { bordered: false, color: getPluginTypeTagColor(row.type) },
          { default: () => row.type },
        ),
    },
    {
      title: 'Runtime',
      key: 'runtime',
      render: row =>
        h(
          NTag,
          { bordered: false, color: { color: '#e2e8f0', textColor: '#334155' } },
          { default: () => row.runtime },
        ),
    },
    { title: 'Version', key: 'version' },
    { title: 'Scope', key: 'scope' },
  ]

  const runtimeOptions = computed(() =>
    [...new Set(rows.value.map(row => row.runtime))]
      .filter(Boolean)
      .map(value => ({ label: value, value })),
  )

  const scopeOptions = computed(() =>
    [...new Set(rows.value.map(row => row.scope))]
      .filter(Boolean)
      .map(value => ({ label: value, value })),
  )

  const typeOptions = pluginTypeOptions.map(option => ({ ...option }))

  const filteredRows = computed(() => {
    const keyword = filters.keyword.trim().toLowerCase()

    return rows.value.filter((row) => {
      const matchesKeyword = !keyword || [
        row.name,
        row.entrypoint,
        row.version,
      ].some(value => value?.toLowerCase().includes(keyword))

      const matchesType = !filters.type || row.type === filters.type
      const matchesRuntime = !filters.runtime || row.runtime === filters.runtime
      const matchesScope = !filters.scope || row.scope === filters.scope

      return matchesKeyword && matchesType && matchesRuntime && matchesScope
    })
  })

  const resetFilters = () => {
    filters.keyword = ''
    filters.type = null
    filters.runtime = null
    filters.scope = null
  }

  const loadPlugins = async() => {
    loading.value = true
    try {
      rows.value = await fetchPlugins()
    } catch (error) {
      message.error(
        error instanceof ApiError ? error.message : 'Unable to load plugins.',
      )
    } finally {
      loading.value = false
    }
  }

  onMounted(loadPlugins)

  return {
    resetFilters,
    filteredRows,
    scopeOptions,
    runtimeOptions,
    typeOptions,
    filters,
    columns,
    loadPlugins,
    loading,
  }
}
