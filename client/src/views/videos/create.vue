<script setup lang="ts">
import { onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'

import { useUserStore } from '@/stores/user'
import { useVideoStore } from '@/stores/video'
import { useVideoQuoteAudiosStore } from '@/stores/videoQuoteAudios'
import { useVideoQuotesStore } from '@/stores/videoQuotes'

const router = useRouter()
const userStore = useUserStore()
const videoStore = useVideoStore()
const videoQuoteStore = useVideoStore()
const videoQuoteAudiosStore = useVideoQuoteAudiosStore()
const videoQuotesStore = useVideoQuotesStore()

const methods = {
  resetVideoStore() {
    videoStore.reset()
    videoQuoteStore.reset()
    videoQuoteAudiosStore.reset()
  },
  async createVideo() {
    if (!userStore.state.user!.id) throw new Error('User is not logged in')
    videoStore.actions.setVideoProperty('userId', userStore.state.user!.id)
    const createdVideo = await videoStore.actions.createVideo()

    const userId = userStore.state.user!.id
    const position = 0
    await videoQuotesStore.actions.createNewVideoQuote(userId, createdVideo.id!, position)

    router.replace({ name: 'videos.edit', params: { id: createdVideo.id } })
    return createdVideo
  },
}

onBeforeMount(async () => {
  methods.resetVideoStore()
  await methods.createVideo()
})
</script>

<template>&nbsp;</template>
