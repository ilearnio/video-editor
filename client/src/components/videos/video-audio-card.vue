<script setup lang="ts">
import { computed } from 'vue'

import { RemoveCircleOutline as DeleteIcon } from '@vicons/ionicons5'

import { getFileUrl } from '@/services/pocketbase'

import AudioPlayer from './audio-player.vue'

const props = defineProps<{
  videoId: string
  file: string | File
  duration: number // in seconds
}>()

const emit = defineEmits<{
  (name: 'delete'): void
}>()

const getters = {
  audioFileUrl: computed(() => {
    if (typeof props.file == 'string') {
      return getFileUrl('videos', props.videoId, props.file)
    } else {
      return URL.createObjectURL(props.file)
    }
  }),
  readableFileName: computed(() => {
    return typeof props.file === 'string'
      ? props.file.replace(/_[a-zA-Z0-9]+(\..+)$/, '$1') // remove pocketbase extension
      : props.file.name
  }),
}
</script>

<template>
  <audio-player
    class="audio-player"
    :src="getters.audioFileUrl.value"
    :name="getters.readableFileName.value"
    :duration="props.duration"
  >
    <template #right>
      <a class="delete-link" title="Remove" @click="emit('delete')">
        <delete-icon class="delete-icon" />
      </a>
    </template>
  </audio-player>
</template>

<style scoped lang="scss">
.audio-player {
  width: 100%;
}

.delete-link {
  cursor: pointer;
  display: block;
  padding: 2px;
  margin-left: 5px;

  &:hover .delete-icon {
    color: #111;
  }
}

.delete-icon {
  width: 20px;
  display: block;
  margin-top: -1px;
  color: #555;
}
</style>
