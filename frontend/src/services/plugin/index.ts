import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NPopconfirm, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

import { permission } from '@/services/access/rbac'
import { createPlugin, deletePlugin, fetchPlugins, updatePlugin } from '@/api'
import { getPluginTypeTagColor, pluginTypeOptions } from '@/theme/plugin'
import { useAuthStore } from '@/stores/modules/auth'
import type { PluginPayload, PluginRecord } from '@/api'
import { ApiError } from '@/api/request'

const emptyPluginForm = (): PluginPayload => ({
  name: '',
  version: '',
  type: 'scaffolder',
  runtime: 'python',
  entrypoint: '',
  scope: 'global',
  description: '',
  enabled: true,
})

export function usePluginService() {
  const message = useMessage()
  const authStore = useAuthStore()
  const loading = ref(false)
  const saving = ref(false)
  const deletingPluginId = ref('')
  const togglingPluginId = ref('')
  const rows = ref<PluginRecord[]>([])
  const formModalOpen = ref(false)
  const detailModalOpen = ref(false)
  const editingPluginId = ref('')
  const selectedPlugin = ref<PluginRecord | null>(null)
  const form = reactive<PluginPayload>(emptyPluginForm())
  const filters = reactive({
    keyword: '',
    type: null as string | null,
    runtime: null as string | null,
    scope: null as string | null,
  })

  const canManagePlugins = computed(() =>
    authStore.canAccess({ permissions: [permission.pluginWrite] }),
  )

  const formTitle = computed(() =>
    editingPluginId.value ? 'Edit plugin' : 'New plugin',
  )

  const runtimeSelectOptions = [
    { label: 'Python', value: 'python' },
    { label: 'Go', value: 'go' },
    { label: 'Node', value: 'node' },
  ]

  const scopeSelectOptions = [
    { label: 'Global', value: 'global' },
    { label: 'Project', value: 'project' },
    { label: 'Environment', value: 'environment' },
  ]

  const fillForm = (plugin?: PluginRecord | null) => {
    const next = plugin
      ? {
          name: plugin.name,
          version: plugin.version,
          type: plugin.type,
          runtime: plugin.runtime,
          entrypoint: plugin.entrypoint,
          scope: plugin.scope,
          description: plugin.description,
          enabled: plugin.enabled ?? true,
        }
      : emptyPluginForm()

    Object.assign(form, next)
  }

  const openCreatePlugin = () => {
    editingPluginId.value = ''
    fillForm()
    formModalOpen.value = true
  }

  const openEditPlugin = (plugin: PluginRecord) => {
    editingPluginId.value = plugin.id
    fillForm(plugin)
    formModalOpen.value = true
  }

  const openPluginDetails = (plugin: PluginRecord) => {
    selectedPlugin.value = plugin
    detailModalOpen.value = true
  }

  const requirePluginForm = () => {
    if (!form.name.trim()) return 'Plugin name is required.'
    if (!form.version.trim()) return 'Version is required.'
    if (!form.type.trim()) return 'Type is required.'
    if (!form.runtime.trim()) return 'Runtime is required.'
    if (!form.entrypoint.trim()) return 'Entrypoint is required.'
    if (!form.scope.trim()) return 'Scope is required.'
    if (!form.description.trim()) return 'Description is required.'
    return null
  }

  const toPayload = (): PluginPayload => ({
    name: form.name.trim(),
    version: form.version.trim(),
    type: form.type,
    runtime: form.runtime,
    entrypoint: form.entrypoint.trim(),
    scope: form.scope,
    description: form.description.trim(),
    enabled: form.enabled ?? true,
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
    {
      title: 'State',
      key: 'enabled',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: row.enabled === false
              ? { color: '#e5e7eb', textColor: '#4b5563' }
              : { color: '#dcfce7', textColor: '#15803d' },
          },
          { default: () => row.enabled === false ? 'Disabled' : 'Enabled' },
        ),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: row =>
        h(
          'div',
          { class: 'flex flex-wrap gap-2' },
          [
            h(
              NButton,
              {
                size: 'small',
                onClick: () => openPluginDetails(row),
              },
              { default: () => 'Details' },
            ),
            canManagePlugins.value
              ? h(
                  NButton,
                  {
                    size: 'small',
                    onClick: () => openEditPlugin(row),
                  },
                  { default: () => 'Edit' },
                )
              : null,
            canManagePlugins.value
              ? h(
                  NButton,
                  {
                    size: 'small',
                    loading: togglingPluginId.value === row.id,
                    onClick: () => togglePluginEnabled(row),
                  },
                  { default: () => row.enabled === false ? 'Enable' : 'Disable' },
                )
              : null,
            canManagePlugins.value
              ? h(
                  NPopconfirm,
                  {
                    onPositiveClick: () => removePlugin(row),
                  },
                  {
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: 'small',
                          type: 'error',
                          ghost: true,
                          loading: deletingPluginId.value === row.id,
                        },
                        { default: () => 'Delete' },
                      ),
                    default: () => `Delete ${row.name}?`,
                  },
                )
              : null,
          ],
        ),
    },
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

  const submitPlugin = async() => {
    const validationError = requirePluginForm()
    if (validationError) {
      message.warning(validationError)
      return
    }

    saving.value = true
    try {
      if (editingPluginId.value) {
        const updated = await updatePlugin(editingPluginId.value, toPayload())
        rows.value = rows.value.map(row => row.id === updated.id ? updated : row)
        if (selectedPlugin.value?.id === updated.id) {
          selectedPlugin.value = updated
        }
        message.success('Plugin updated successfully.')
      } else {
        const created = await createPlugin(toPayload())
        rows.value = [created, ...rows.value]
        message.success('Plugin created successfully.')
      }
      formModalOpen.value = false
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to save plugin.')
    } finally {
      saving.value = false
    }
  }

  const togglePluginEnabled = async(row: PluginRecord) => {
    togglingPluginId.value = row.id
    try {
      const updated = await updatePlugin(row.id, { enabled: !(row.enabled ?? true) })
      rows.value = rows.value.map(item => item.id === updated.id ? updated : item)
      if (selectedPlugin.value?.id === updated.id) {
        selectedPlugin.value = updated
      }
      message.success(updated.enabled ? 'Plugin enabled.' : 'Plugin disabled.')
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to update plugin state.')
    } finally {
      togglingPluginId.value = ''
    }
  }

  const removePlugin = async(row: PluginRecord) => {
    deletingPluginId.value = row.id
    try {
      await deletePlugin(row.id)
      rows.value = rows.value.filter(item => item.id !== row.id)
      if (selectedPlugin.value?.id === row.id) {
        detailModalOpen.value = false
        selectedPlugin.value = null
      }
      message.success('Plugin deleted successfully.')
    } catch (error) {
      message.error(error instanceof ApiError ? error.message : 'Unable to delete plugin.')
    } finally {
      deletingPluginId.value = ''
    }
  }

  onMounted(loadPlugins)

  return {
    canManagePlugins,
    columns,
    detailModalOpen,
    filteredRows,
    filters,
    form,
    formModalOpen,
    formTitle,
    loadPlugins,
    loading,
    openCreatePlugin,
    resetFilters,
    runtimeOptions,
    runtimeSelectOptions,
    saving,
    scopeOptions,
    scopeSelectOptions,
    selectedPlugin,
    submitPlugin,
    typeOptions,
  }
}
