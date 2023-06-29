<script setup lang="ts">
import { darkTheme } from 'naive-ui'
import { computed, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

import { useUserStore } from '@/stores/user'

import './assets/styles/main.scss'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const themeOverrides = {
  common: {
    primaryColor: '#60a5fa',
    primaryColorHover: '#3b82f6',
    primaryColorPressed: '#4c1d95',
  },
  Form: {
    feedbackPadding: '4px 0 10px 2px',
  },
  Dropdown: {
    borderRadius: '6px',
  },
}

const theme = computed(() => (userStore.state.isDark ? darkTheme : null))

const default_layout = 'default'
const layout = computed(() => {
  return `${route.meta.layout || default_layout}-layout`
})
watch(
  () => userStore.getters.isLoggedIn,
  (val) => {
    if (!val) router.push({ name: 'login' })
  },
)
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-global-style />
    <n-notification-provider>
      <n-message-provider>
        <Component :is="layout">
          <RouterView />
        </Component>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>
