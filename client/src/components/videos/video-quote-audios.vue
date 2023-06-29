<script setup lang="ts">
import { onBeforeMount, reactive } from 'vue'
import { computed } from 'vue'
import { nextTick } from 'vue'
import PulseLoader from 'vue-spinner/src/PulseLoader.vue'

import GenerateAudioDropdown from '@/components/videos/generate-audio-dropdown.vue'
import { VideoQuoteAudio } from '@/models/videoQuoteAudios'
import * as videoQuoteAudiosApi from '@/services/api/videoQuoteAudios'
import { extractQuotePlainText } from '@/services/videoQuotes'
import { getInvalidReason } from '@/services/videoQuotes'
import { useVideoQuoteAudiosStore } from '@/stores/videoQuoteAudios'
import { useVideoQuotesStore } from '@/stores/videoQuotes'

import VideoQuoteAudioCard from './video-quote-audio-card.vue'

const videoQuotesStore = useVideoQuotesStore()
const videoQuoteAudiosStore = useVideoQuoteAudiosStore()

const props = defineProps<{
  videoId: string
}>()

const data = reactive({
  loadingIds: [] as string[],
  generatorPopoverVisibility: false,
})

const getters = {
  isLoading: computed(() => {
    return data.loadingIds.includes(videoQuotesStore.state.activeVideoQuoteId)
  }),
}

const methods = {
  extractAudioFileName(audioFile: string | File) {
    return typeof audioFile === 'string'
      ? audioFile.replace(/_[a-zA-Z0-9]+(\..+)$/, '$1') // remove pocketbase extension
      : audioFile.name
  },
  getInvalidReason(audio: VideoQuoteAudio) {
    if (!audio) return
    return getInvalidReason(
      videoQuotesStore.state.videoQuotes,
      videoQuoteAudiosStore.state.videoQuoteAudios,
      audio,
    )
  },
  async handleChange(id: string) {
    await videoQuoteAudiosStore.actions.selectAudioId(id)
  },
  async removeAudio(audio: VideoQuoteAudio) {
    await videoQuoteAudiosStore.actions.removeVideoQuoteAudio(audio.id!)
  },
  async handleAudioGeneratorSubmit(voice: string, speed: number, seed: number) {
    data.generatorPopoverVisibility = false
    await methods.generateAudio(voice, speed, seed)
  },
  async generateAudio(voice: string, speed: number, seed: number) {
    await nextTick()

    const { activeVideoQuoteId } = videoQuotesStore.state
    if (!activeVideoQuoteId) {
      throw new Error('Video quote id is required')
    }

    if (!videoQuotesStore.getters.activeVideoQuote?.content) {
      throw new Error('Video quote content is required')
    }

    data.loadingIds.push(activeVideoQuoteId)

    const quote = videoQuotesStore.getters.activeVideoQuote

    try {
      const createdVideoQuoteAudio = await videoQuoteAudiosStore.actions.generateAudio(
        extractQuotePlainText(quote),
        voice,
        speed,
        seed,
      )
      videoQuotesStore.actions.setVideoQuoteProperty(
        activeVideoQuoteId,
        'selectedAudioId',
        createdVideoQuoteAudio.id!,
      )
    } finally {
      const index = data.loadingIds.findIndex((id) => id === activeVideoQuoteId)
      data.loadingIds.splice(index, 1)
    }
  },
}

onBeforeMount(async () => {
  const videoQuoteAudios = await videoQuoteAudiosApi.getVideoQuoteAudios(props.videoId)
  videoQuoteAudiosStore.actions.setVideoQuoteAudios(videoQuoteAudios)
})
</script>

<template>
  <div class="quote-audios">
    <div v-if="videoQuotesStore.state.activeVideoQuoteId" class="heading-area">
      <h4>Pick your preferred sample</h4>
      <n-popover v-model:show="data.generatorPopoverVisibility" trigger="click" style="padding: 0">
        <template #trigger>
          <n-button
            class="generate-button"
            :disabled="
              getters.isLoading.value || !videoQuotesStore.getters.activeVideoQuote?.content
            "
            type="success"
            size="small"
          >
            <PulseLoader v-if="getters.isLoading.value" color="#fff" class="loader" size="8px" />
            <span :style="{ visibility: getters.isLoading.value ? 'hidden' : 'unset' }">
              Generate Audio
            </span>
          </n-button>
        </template>
        <generate-audio-dropdown @submit="methods.handleAudioGeneratorSubmit" />
      </n-popover>
    </div>
    <div class="content">
      <template v-if="videoQuotesStore.getters.activeVideoQuote">
        <template v-if="videoQuoteAudiosStore.getters.activeVideoQuoteAudios.length">
          <div
            v-for="audio in videoQuoteAudiosStore.getters.activeVideoQuoteAudios"
            :key="audio.id"
            class="audio"
          >
            <VideoQuoteAudioCard
              :audio="audio"
              :checked="audio.id === videoQuotesStore.getters.activeVideoQuote.selectedAudioId"
              :invalid-reason="methods.getInvalidReason(audio)"
              @remove="methods.removeAudio(audio)"
              @checked="methods.handleChange(audio.id!)"
            />
          </div>
        </template>
        <div v-else class="empty-area">
          <n-empty description="No audios have been generated yet." />
        </div>
      </template>
      <div v-else class="empty-area">
        <n-empty description="Select a quote on the left to view or generate audios." />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.quote-audios {
  background: #fdfdfd;
  width: 100%;
  height: 100%;
  overflow: hidden;
  border: 1px solid #ccc;
  display: flex;
  flex-direction: column;
}

.heading-area {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 1rem;
  height: 3rem;
  background: #fff;
  border-bottom: 1px solid #eaeaea;
}

.content {
  padding: 1rem;
  overflow-y: auto;
  flex: 1;
}

.audio {
  display: flex;
  align-items: center;

  + .audio {
    margin-top: 5px;
  }
}

.generate-button {
  position: relative;
}

.loader {
  position: absolute;
  inset-inline: -1in;
  scale: 1;
}

.empty-area {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 110%;
  color: #444;
  height: 100%;

  :deep(*) {
    color: #999;
  }
}
</style>
