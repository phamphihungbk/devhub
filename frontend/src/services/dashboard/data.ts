import { h } from 'vue'
import { getEnvironmentTagColor } from '@/theme/environment'
import { getPluginTypeTagColor } from '@/theme/plugin'
import { NTag } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import type { PluginRecord, Project } from '@/api'

function environmentName(value: NonNullable<Project['environments']>[number]) {
  return typeof value === 'string' ? value : value.name
}

const projectColumns: DataTableColumns<Project> = [
    {
      title: 'Project',
      key: 'name',
    },
    {
      title: 'Environments',
      key: 'environments',
      render: row =>
        h(
          'div',
          { class: 'flex flex-wrap gap-2' },
          (row.environments || []).map((value) =>
            h(
              NTag,
              {
                bordered: false,
                color: getEnvironmentTagColor(environmentName(value)),
              },
              { default: () => environmentName(value) },
            ),
          ).concat((row.environments || []).length === 0 ? [h('span', 'Not configured')] : []),
        ),
    },
    {
      title: 'Description',
      key: 'description',
      render: row => row.description || 'No description yet.',
    },
  ];

  const pluginColumns: DataTableColumns<PluginRecord> = [
    { title: 'Plugin', key: 'name' },
    {
      title: 'Type',
      key: 'type',
      render: row =>
        h(
          NTag,
          {
            bordered: false,
            color: getPluginTypeTagColor(row.type),
          },
          { default: () => row.type },
        ),
    },
    // { title: 'Runtime', key: 'runtime' },
    { title: 'Scope', key: 'scope' },
  ];


export function useDashboardData() {
    return {
        projectColumns,
        pluginColumns
    }
}
