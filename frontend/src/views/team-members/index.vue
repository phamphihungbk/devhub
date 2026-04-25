<script setup lang="ts">
import { NButton, NCard, NDataTable } from 'naive-ui'

import PageHeader from '@/components/page-header.vue'
import { useTeamMemberService } from '@/services/team-member'

const { columns, loadUsers, loading, rows } = useTeamMemberService()
</script>

<template>
  <div>
    <PageHeader
      eyebrow="Access"
      title="Team members"
      description="The operator directory for the control plane, ready for role governance, invite flows, and access reviews."
    >
      <NButton @click="loadUsers">
        Refresh
      </NButton>
    </PageHeader>

    <NCard class="rounded-3xl border border-[var(--app-border)] shadow-[var(--app-shadow)]">
      <NDataTable
        :columns="columns"
        :data="rows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
      />
    </NCard>
  </div>
</template>
