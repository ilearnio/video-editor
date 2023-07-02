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

const quotesStore = useVideoQuotesStore()
const audiosStore = useVideoQuoteAudiosStore()

const props = defineProps<{
  videoId: string
}>()

const data = reactive({
  loadingIds: [] as string[],
  generatorPopoverVisibility: false,
})

const getters = {
  isLoading: computed(() => {
    return data.loadingIds.includes(quotesStore.state.activeVideoQuoteId)
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
      quotesStore.state.videoQuotes,
      audiosStore.state.videoQuoteAudios,
      audio,
    )
  },
  async handleChange(audioId: string) {
    const quoteId = quotesStore.state.activeVideoQuoteId
    await audiosStore.actions.selectAudioId(audioId)
    quotesStore.actions.setVideoQuoteProperty(quoteId, 'selectedAudioId', audioId)
  },
  async deleteAudio(audio: VideoQuoteAudio) {
    const quote = quotesStore.getters.activeVideoQuote!
    const isSelected = quote.selectedAudioId === audio.id!
    const selectedAudio = await audiosStore.actions.deleteVideoQuoteAudio(audio.id!, isSelected)
    if (selectedAudio) {
      quotesStore.actions.setVideoQuoteProperty(quote.id!, 'selectedAudioId', selectedAudio?.id)
    }
  },
  async handleAudioGeneratorSubmit(voice: string, speed: number, seed: number) {
    data.generatorPopoverVisibility = false
    await methods.generateAudio(voice, speed, seed)
  },
  async generateAudio(voice: string, speed: number, seed: number) {
    await nextTick()

    const { activeVideoQuoteId } = quotesStore.state
    if (!activeVideoQuoteId) {
      throw new Error('Video quote id is required')
    }

    if (!quotesStore.getters.activeVideoQuote?.content) {
      throw new Error('Video quote content is required')
    }

    data.loadingIds.push(activeVideoQuoteId)

    const quote = quotesStore.getters.activeVideoQuote

    try {
      const createdVideoQuoteAudio = await audiosStore.actions.generateAudio(
        extractQuotePlainText(quote),
        voice,
        speed,
        seed,
      )
      quotesStore.actions.setVideoQuoteProperty(
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
  audiosStore.actions.setVideoQuoteAudios(videoQuoteAudios)
})
</script>

<template>
  <div class="quote-audios">
    <div v-if="quotesStore.state.activeVideoQuoteId" class="heading-area">
      <h4>Pick your preferred sample</h4>
      <n-popover v-model:show="data.generatorPopoverVisibility" trigger="click" style="padding: 0">
        <template #trigger>
          <n-button
            class="generate-button"
            :disabled="getters.isLoading.value || !quotesStore.getters.activeVideoQuote?.content"
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
      <template v-if="quotesStore.getters.activeVideoQuote">
        <template v-if="audiosStore.getters.activeVideoQuoteAudios.length">
          <div
            v-for="audio in audiosStore.getters.activeVideoQuoteAudios"
            :key="audio.id"
            class="audio"
          >
            <video-quote-audio-card
              :audio="audio"
              :checked="audio.id === quotesStore.getters.activeVideoQuote.selectedAudioId"
              :invalid-reason="methods.getInvalidReason(audio)"
              @checked="methods.handleChange(audio.id!)"
              @delete="methods.deleteAudio(audio)"
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
