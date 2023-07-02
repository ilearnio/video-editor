<script setup lang="ts">
import { NInput } from 'naive-ui'
import { nextTick, onBeforeMount, reactive, ref } from 'vue'

import { ContentCopyFilled as IconCopy } from '@vicons/material'

import { VideoQuote } from '@/models/videoQuotes'
import { useUserStore } from '@/stores/user'
import { useVideoQuoteAudiosStore } from '@/stores/videoQuoteAudios'
import { useVideoQuotesStore } from '@/stores/videoQuotes'

import VideoQuoteAudios from './video-quote-audios.vue'
import VideoQuoteCard from './video-quote-card.vue'

const videoQuotesStore = useVideoQuotesStore()
const videoQuoteAudiosStore = useVideoQuoteAudiosStore()
const userStore = useUserStore()

const containerRef = ref<HTMLElement>()
const nInputRefs = ref<(typeof NInput)[]>()

let isDragging = false
let draggingInitialQuotePositions: Record<string, number> = {}

const props = defineProps<{
  videoId: string
}>()

const data = reactive({
  hasQuoteChanged: false,
  draggingIndex: -1,
  dropdownVisible: {} as Record<string, boolean>,
})

const methods = {
  async selectRow(quoteId: string, focus = true) {
    videoQuotesStore.actions.setActiveVideoQuoteId(quoteId)
    videoQuoteAudiosStore.actions.setActiveVideoQuoteId(quoteId)
    if (focus) {
      await nextTick()
      nInputRefs.value?.find((component) => component.$attrs['data-id'] === quoteId)?.focus()
    }
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
      <video-quote-card
        v-for="(quote, index) in videoQuotesStore.state.videoQuotes"
        :key="quote.id"
        :quote="quote"
        :index="index"
        :video-id="props.videoId"
        :dragging="index === data.draggingIndex"
        @select="(quoteId, focus) => methods.selectRow(quoteId, focus)"
        @drag-handle-mousedown="methods.dragStart(index, $event)"
      />
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

.quote-audios-area {
  height: 500px;
  width: 420px;
  position: sticky;
  top: 55px;
  z-index: 5;
}
</style>
