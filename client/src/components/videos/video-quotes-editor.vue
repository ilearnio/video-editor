<script setup lang="ts">
import { NInput } from 'naive-ui'
import { nextTick, onBeforeMount, reactive, ref } from 'vue'
import { computed } from 'vue'

import {
  InformationCircleOutline as InfoIcon,
  RemoveCircleOutline as RemoveIcon,
} from '@vicons/ionicons5'
import { DragIndicatorSharp as DragIcon } from '@vicons/material'

import { upperFirst } from '@/helpers/strings'
import { VideoQuote } from '@/models/videoQuotes'
import { getInvalidReason } from '@/services/videoQuotes'
import { useUserStore } from '@/stores/user'
import { useVideoQuoteAudiosStore } from '@/stores/videoQuoteAudios'
import { useVideoQuotesStore } from '@/stores/videoQuotes'

import VideoQuoteAudios from './video-quote-audios.vue'

const videoQuotesStore = useVideoQuotesStore()
const videoQuoteAudiosStore = useVideoQuoteAudiosStore()
const userStore = useUserStore()

const containerRef = ref<HTMLElement>()
const dragHandleRefs = ref<HTMLElement[]>()
const nInputRefs = ref<(typeof NInput)[]>()

let isDragging = false
let draggingInitialQuotePositions: Record<string, number> = {}

const props = defineProps<{
  videoId: string
}>()

const data = reactive({
  hasQuoteChanged: false,
  draggingIndex: -1,
})

const getters = {
  maxPosition: computed(() => {
    const positions = videoQuotesStore.state.videoQuotes.map((quote) => quote.position)
    return positions.sort((a, b) => b - a)[0] ?? -1
  }),
  quoteVoiceDetails: computed(() => (quote: VideoQuote) => {
    if (!quote.selectedAudioId) return ''
    const selectedAudio = videoQuoteAudiosStore.state.videoQuoteAudios.find(
      (a) => a.id === quote.selectedAudioId,
    )
    if (!selectedAudio) return ''

    const voice = upperFirst(selectedAudio.voice)
    const speed = Number(selectedAudio.speed.toFixed(2))

    return `${voice} - ${speed}x`
  }),
  getInvalidReason: computed(() => (quote: VideoQuote) => {
    if (!quote.selectedAudioId) return ''
    const selectedAudio = videoQuoteAudiosStore.state.videoQuoteAudios.find(
      (a) => a.id === quote.selectedAudioId,
    )
    if (!selectedAudio) return 'Audio must be selected.'

    return getInvalidReason(
      videoQuotesStore.state.videoQuotes,
      videoQuoteAudiosStore.state.videoQuoteAudios,
      selectedAudio,
    ) as string | undefined
  }),
}

const methods = {
  async selectRow(quoteId: string, focus = true) {
    videoQuotesStore.actions.setActiveVideoQuoteId(quoteId)
    videoQuoteAudiosStore.actions.setActiveVideoQuoteId(quoteId)
    if (focus) {
      await nextTick()
      nInputRefs.value?.find((component) => component.$attrs['data-id'] === quoteId)?.focus()
    }
  },

  async toggleHTML(quote: VideoQuote) {
    videoQuotesStore.actions.setVideoQuoteProperty(quote.id!, 'isHtmlEnabled', !quote.isHtmlEnabled)
    await videoQuotesStore.actions.updateVideoQuote(quote.id!)
  },

  handleInputFocus(quoteId: string) {
    methods.selectRow(quoteId, false)
  },

  handleInputKeydown(e: KeyboardEvent) {
    if (e.key !== 'Tab') return

    const isCreatingNewQuote =
      videoQuotesStore.getters.activeVideoQuoteIndex ===
      videoQuotesStore.state.videoQuotes.length - 1
    if (e.ctrlKey || e.altKey || e.metaKey || (!e.shiftKey && isCreatingNewQuote)) {
      e.preventDefault()
      return
    }

    const { activeVideoQuoteIndex } = videoQuotesStore.getters

    if (activeVideoQuoteIndex === -1 || (e.shiftKey && activeVideoQuoteIndex === 0)) {
      e.preventDefault()
      return
    }

    if (e.shiftKey && videoQuotesStore.state.videoQuotes[activeVideoQuoteIndex - 1]) {
      videoQuotesStore.actions.setActiveVideoQuoteId(
        videoQuotesStore.state.videoQuotes[activeVideoQuoteIndex - 1].id!,
      )
    } else if (videoQuotesStore.state.videoQuotes[activeVideoQuoteIndex + 1]) {
      videoQuotesStore.actions.setActiveVideoQuoteId(
        videoQuotesStore.state.videoQuotes[activeVideoQuoteIndex + 1].id!,
      )
    }
  },

  async handleInputEvent() {
    data.hasQuoteChanged = true

    const activeQuoteIndex = videoQuotesStore.getters.activeVideoQuoteIndex
    const activeQuote = videoQuotesStore.getters.activeVideoQuote

    if (
      activeQuoteIndex === videoQuotesStore.state.videoQuotes.length - 1 &&
      activeQuote?.content !== ''
    ) {
      const userId = userStore.state.user!.id
      const position = getters.maxPosition.value + 1
      await videoQuotesStore.actions.createNewVideoQuote(userId, props.videoId, position)
    }
  },

  async handleInputBlur(quoteId: string) {
    if (data.hasQuoteChanged) {
      data.hasQuoteChanged = false
      await videoQuotesStore.actions.updateVideoQuote(quoteId)
    }
  },

  async removeQuote(quoteId: string) {
    await videoQuotesStore.actions.removeVideoQuote(quoteId)
    videoQuoteAudiosStore.actions.visuallyRemoveVideoQuoteAudios(quoteId)
  },

  dragStart(index: number, event: MouseEvent) {
    const dragHandle = event.target as HTMLElement
    if (!dragHandle.classList.contains('drag-handle') && !dragHandle.closest('.drag-handle')) return

    event.preventDefault()

    videoQuotesStore.state.activeVideoQuoteId = ''
    data.draggingIndex = index
    isDragging = true

    draggingInitialQuotePositions = {}
    videoQuotesStore.state.videoQuotes.reduce((acc, quote) => {
      acc[quote.id!] = quote.position
      return acc
    }, draggingInitialQuotePositions)

    document.addEventListener('mousemove', methods.dragHandle)
    document.addEventListener('mouseup', methods.dragEnd)
  },

  dragHandle(event: MouseEvent) {
    event.preventDefault()
    if (!isDragging) return

    const newIndex = methods.dragCalculateNewIndex(event.clientY)
    if (newIndex !== data.draggingIndex) {
      const [draggedRow] = videoQuotesStore.state.videoQuotes.splice(data.draggingIndex, 1)
      videoQuotesStore.state.videoQuotes.splice(newIndex, 0, draggedRow)
      data.draggingIndex = newIndex
    }
  },

  async dragEnd() {
    if (!isDragging) return
    data.draggingIndex = -1
    isDragging = false
    document.removeEventListener('mousemove', methods.dragHandle)
    document.removeEventListener('mouseup', methods.dragEnd)

    await nextTick()

    const repositionedQuotes: VideoQuote[] = []
    containerRef.value!.querySelectorAll('.quote-row').forEach((el, index) => {
      const quoteId = (el as HTMLElement).dataset.id as string
      const quote = videoQuotesStore.state.videoQuotes.find((q) => q.id === quoteId)
      if (!quote) return

      videoQuotesStore.actions.setVideoQuoteProperty(quoteId, 'position', index)

      const oldPosition = draggingInitialQuotePositions[quoteId]
      if (oldPosition !== undefined && oldPosition !== index) {
        repositionedQuotes.push({ ...quote, position: index })
      }
    })

    if (repositionedQuotes.length) {
      await videoQuotesStore.actions.updateVideoQuotes(repositionedQuotes)
    }
  },

  dragCalculateNewIndex(clientY: number) {
    const quoteRows = Array.from(containerRef.value!.querySelectorAll('.quote-row'))

    // Find the index of the row over which the cursor is pointing
    const newIndex = quoteRows.findIndex((row) => {
      const rowRect = row.getBoundingClientRect()
      return clientY >= rowRect.top && clientY <= rowRect.bottom
    })

    // If the cursor is not pointing over any row, return the current dragging index
    if (newIndex === -1) {
      return data.draggingIndex
    }

    // Return the new index
    return newIndex
  },
}

onBeforeMount(async () => {
  await videoQuotesStore.actions.getVideoQuotes(userStore.state.user!.id, props.videoId)
})
</script>

<template>
  <n-h3>Quotes ({{ videoQuotesStore.state.videoQuotes.length - 1 }})</n-h3>
  <div ref="containerRef" class="container">
    <div class="column-left">
      <div
        v-for="(quote, index) in videoQuotesStore.state.videoQuotes"
        :key="quote.id"
        :data-id="quote.id"
        class="quote-row"
        :class="{
          active: quote.id === videoQuotesStore.state.activeVideoQuoteId,
          dragging: index === data.draggingIndex,
          'quote-new': index === videoQuotesStore.state.videoQuotes.length - 1,
          invalid: getters.getInvalidReason.value(quote),
        }"
      >
        <div class="text-content-wrapper" @click="methods.selectRow(quote.id!)">
          <n-input
            v-if="
              quote.id === videoQuotesStore.state.activeVideoQuoteId ||
              index === videoQuotesStore.state.videoQuotes.length - 1
            "
            ref="nInputRefs"
            v-model:value="quote.content"
            class="text-content-input"
            type="textarea"
            :data-id="quote.id"
            :autosize="{ minRows: 1, maxRows: 3 }"
            @focus="methods.handleInputFocus(quote.id!)"
            @keydown="methods.handleInputKeydown"
            @input="methods.handleInputEvent"
            @blur="methods.handleInputBlur(quote.id!)"
          />
          <div v-else class="text-content">
            {{ quote.content || '&nbsp;' }}
          </div>

          <template v-if="index !== videoQuotesStore.state.videoQuotes.length - 1">
            <n-tooltip
              v-if="getters.getInvalidReason.value(quote)"
              trigger="hover"
              placement="bottom"
            >
              <template #trigger>
                <div class="info-icon-wrapper" @click.stop>
                  <InfoIcon class="info-icon" />
                </div>
              </template>
              {{ getters.getInvalidReason.value(quote) }}
            </n-tooltip>

            <n-tag
              v-if="quote.selectedAudioId"
              class="voice-details"
              type="default"
              size="small"
              @click.stop
            >
              {{ getters.quoteVoiceDetails.value(quote) }}
            </n-tag>

            <n-button
              size="tiny"
              :type="quote.isHtmlEnabled ? 'info' : 'tertiary'"
              class="toggle-html-button"
              tabindex="-1"
              @click.stop="methods.toggleHTML(quote)"
            >
              HTML
            </n-button>

            <n-popconfirm
              :negative-text="null"
              placement="bottom"
              @positive-click="methods.removeQuote(quote.id!)"
            >
              <template #trigger>
                <a class="remove-link" @click.stop>
                  <RemoveIcon />
                </a>
              </template>
              Delete the quote including it's audio files?
            </n-popconfirm>

            <div
              v-if="videoQuotesStore.state.videoQuotes.length > 1"
              ref="dragHandleRefs"
              class="drag-handle"
              @mousedown="methods.dragStart(index, $event)"
              @click.stop
            >
              <drag-icon />
            </div>
          </template>
        </div>
      </div>
    </div>
    <div class="column-right">
      <div class="quote-audios-area">
        <video-quote-audios :video-id="props.videoId" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.container {
  width: 100%;
  position: relative;
  z-index: 0;
  display: flex;
  padding-bottom: 100px;
}

.column-left {
  flex: 1;
  margin-right: 20px;
}

.quote-row {
  display: flex;
  align-items: flex-start;
  position: relative;
  width: 100%;

  + .quote-row {
    margin-top: 10px;
  }
}

.drag-handle {
  display: flex;
  align-items: center;
  justify-content: start;
  align-self: stretch;
  cursor: move;
  color: #999;

  svg {
    display: block;
    width: 19px;
  }
}

.text-content-wrapper {
  display: flex;
  align-items: center;
  flex: 1;
  padding: 7px 7px 7px 10px;
  background-color: #fdfdfd;
  border: 1px solid #ddd;
  cursor: pointer;
  border-radius: 2px;
  position: relative;
  z-index: 1;
}

.text-content-input {
  flex: 1;
  margin: -7px;

  :deep(.n-input-wrapper) {
    padding-left: 7px;
  }

  :deep(.n-input__textarea textarea) {
    border: 0;
  }

  :deep(.n-input__state-border) {
    display: none;
  }

  :deep(.n-input__border) {
    display: none;
  }
}

.text-content {
  flex: 1;
}

.quote-row.active {
  .text-content-wrapper {
    background-color: #fefefe;
    border-color: #2aa0db;
    box-shadow: 0 0 0 0.5px #2aa0db;
  }
}

.dragging .text-content-wrapper {
  border-color: #48975b;
  box-shadow: 0 0 0 0.5px #48975b;
}

.voice-details {
  margin: 0 -5px 0 12px;
  opacity: 0.9;
}

.info-icon-wrapper {
  cursor: help;
  margin: 0 -5px 0 12px;
}

.info-icon {
  width: 20px;
  color: #f70;
  display: block;
}

.toggle-html-button {
  margin-left: 12px;
}

.remove-link {
  color: #999;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  margin-left: 5px;

  &:hover {
    color: #333;
  }

  svg {
    width: 16px;
  }
}

.quote-audios-area {
  height: 500px;
  width: 420px;
  position: sticky;
  top: 55px;
  z-index: 5;
}
</style>
