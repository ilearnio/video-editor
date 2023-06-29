import { acceptHMRUpdate, defineStore } from 'pinia'
import { computed, reactive } from 'vue'

import { usePreferredDark, useStorage } from '@vueuse/core'

import { type User } from '@/models/users'
import { pocketbaseParseAuth } from '@/services/pocketbaseAuth'

interface State {
  token: string | undefined
  user: User | undefined
  isLoading: boolean
  isDark: boolean
  sidebarCollapsed: boolean
}

export const useUserStore = defineStore('userStore', () => {
  const pocketbaseAuthRef = useStorage<string>('pocketbase_auth', null)
  const isDarkRef = usePreferredDark()
  const sidebarCollapsedRef = useStorage<boolean>('sidebarCollapsed', false)

  const state: State = reactive<State>({
    get token() {
      const { token } = pocketbaseParseAuth(pocketbaseAuthRef.value) || {}
      return token
    },
    get user() {
      const { model } = pocketbaseParseAuth(pocketbaseAuthRef.value) || {}
      return model
    },
    isLoading: false,
    get isDark() {
      return isDarkRef.value
    },
    set isDark(value: boolean) {
      isDarkRef.value = value
    },
    get sidebarCollapsed() {
      return sidebarCollapsedRef.value
    },
    set sidebarCollapsed(value: boolean) {
      sidebarCollapsedRef.value = value
    },
  })

  const getters = {
    isLoggedIn: computed(() => {
      return state.token !== undefined && state.token !== '' && state.token != null
    }),
  }

  const actions = {
    async toggleSidebarCollapsed() {
      state.sidebarCollapsed = !state.sidebarCollapsed
    },
    async setUser(userData: User) {
      state.user = userData
    },
    async setIsDark(value: boolean) {
      state.isDark = value
    },
    async setIsLoading(value: boolean) {
      state.isLoading = value
    },
  }

  return {
    state,
    getters,
    actions,
  } as const
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
}
