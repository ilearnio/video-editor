<script setup lang="ts">
import { useUserStore } from '@/stores/user'

import Content from './Content.vue'
import Footer from './Footer.vue'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'

const userStore = useUserStore()
</script>

<template>
  <PreLoader v-if="userStore.state.isLoading" />
  <n-layout position="absolute" has-sider>
    <Sidebar />
    <n-layout>
      <Header />
      <Content />
      <!-- <Footer /> -->
    </n-layout>
  </n-layout>
</template>
