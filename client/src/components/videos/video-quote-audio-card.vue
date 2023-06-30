<script setup lang="ts">
import { computed } from 'vue'

import {
  InformationCircleOutline as InfoIcon,
  RemoveCircleOutline as RemoveIcon,
} from '@vicons/ionicons5'

import type { VideoQuoteAudio } from '@/models/videoQuoteAudios'
import { getFileUrl } from '@/services/pocketbase'
import { InvalidVideoQuoteAudioReason, extractAudioPlainText } from '@/services/videoQuotes'

import AudioPlayer from './audio-player.vue'

const props = defineProps<{
  audio: VideoQuoteAudio
  checked: boolean
  invalidReason?: InvalidVideoQuoteAudioReason
}>()

const emit = defineEmits<{
  (name: 'remove'): void
  (name: 'checked', event: InputEvent): void
}>()

const getters = {
  audioFileUrl: computed(() => {
    if (typeof props.audio.audioFile == 'string') {
      return getFileUrl('videoQuoteAudios', props.audio.id!, props.audio.audioFile)
    } else {
      return URL.createObjectURL(props.audio.audioFile)
    }
  }),
  subtext: computed(() => {
    const voice = props.audio.voice.charAt(0).toUpperCase() + props.audio.voice.slice(1)
    return `${voice} - ${props.audio.speed}x`
  }),
}
</script>

<template>
  <audio-player
    class="audio-player"
    :src="getters.audioFileUrl.value"
    :name="extractAudioPlainText(audio)"
    :subtext="getters.subtext.value"
    :duration="audio.duration / 1000"
    :class="{
      'audio-player__invalid': invalidReason,
    }"
  >
    <template #left>
      <div class="left">
        <n-radio :checked="checked" :value="audio.id" @change="emit('checked', $event)" />
      </div>
    </template>
    <template #right>
      <n-tooltip v-if="invalidReason" trigger="hover" placement="bottom">
        <template #trigger>
          <div class="info-icon-wrapper">
            <info-icon class="info-icon" />
          </div>
        </template>
        {{ invalidReason }}
      </n-tooltip>
      <a class="remove-link" title="Remove" @click="emit('remove')">
        <remove-icon class="remove-icon" />
      </a>
    </template>
  </audio-player>
</template>

<style scoped lang="scss">
.audio-player {
  width: 100%;
}

.audio-player__invalid {
  border-color: #f70;
}

.left {
  margin: 0 8px 0 3px;
}

.remove-link {
  cursor: pointer;
  display: block;
  padding: 2px;
  margin-left: 5px;

  &:hover .remove-icon {
    color: #111;
  }
}

.info-icon-wrapper {
  order: -1;
  cursor: help;
  margin: 0 -5px 0 12px;
}

.info-icon {
  width: 20px;
  color: #f70;
  display: block;
}

.remove-icon {
  width: 20px;
  display: block;
  margin-top: -1px;
  color: #555;
}
</style>
