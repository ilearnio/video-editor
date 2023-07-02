<script setup lang="ts">
import { useRouter } from 'vue-router'

import { Pencil as EditIcon, Trash as TrashIcon } from '@vicons/tabler'

import { Video } from '@/models/videos'
import { useVideoStore } from '@/stores/video'

const videosStore = useVideoStore()

const props = defineProps<{
  video: Video
}>()

const emit = defineEmits<{
  (name: 'delete', videoId: string): void
}>()

const router = useRouter()

const methods = {
  editVideo() {
    router.push({ name: 'videos.edit', params: { id: props.video.id } })
  },

  async deleteVideo() {
    await videosStore.actions.deleteVideo(props.video.id!)
    emit('delete', props.video.id!)
  },
}
</script>

<template>
  <n-space>
    <n-button circle quaternary @click="methods.editVideo">
      <n-icon><edit-icon /></n-icon>
    </n-button>
    <n-popconfirm negative-text="" placement="bottom" @positive-click="methods.deleteVideo">
      <template #default>
        Are you sure you want to delete this video? This action cannot be undone.
      </template>
      <template #trigger>
        <n-button type="error" circle tertiary>
          <n-icon><trash-icon /></n-icon>
        </n-button>
      </template>
    </n-popconfirm>
  </n-space>
</template>
