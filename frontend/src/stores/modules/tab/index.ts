import { defineStore } from 'pinia'

import { SetupStoreId } from '@/enum'

export interface TabItem {
  key: string
  label: string
  fullPath: string
  closable: boolean
}

export interface TabState {
  activeTabKey: string
  tabs: TabItem[]
}

const defaultTabs = (): TabItem[] => []

export const useTabStore = defineStore(SetupStoreId.Tab, {
  state: (): TabState => ({
    activeTabKey: '',
    tabs: defaultTabs(),
  }),
  getters: {
    hasTabs(state) {
      return state.tabs.length > 0
    },
  },
  actions: {
    setActiveTab(key: string) {
      this.activeTabKey = key
    },
    upsertTab(tab: TabItem) {
      const existingTabIndex = this.tabs.findIndex(item => item.key === tab.key)

      if (existingTabIndex >= 0) {
        this.tabs[existingTabIndex] = tab
      } else {
        this.tabs.push(tab)
      }

      this.activeTabKey = tab.key
    },
    removeTab(key: string) {
      this.tabs = this.tabs.filter(tab => tab.key !== key)

      if (this.activeTabKey === key) {
        this.activeTabKey = this.tabs.at(-1)?.key || ''
      }
    },
    resetTabs() {
      this.activeTabKey = ''
      this.tabs = defaultTabs()
    },
  },
})
