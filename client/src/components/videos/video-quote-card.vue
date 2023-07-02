<script setup lang="ts">
import { NInput } from 'naive-ui'
import type { DropdownMixedOption } from 'naive-ui/es/dropdown/src/interface'
import { reactive, ref } from 'vue'
import { computed } from 'vue'

import { InformationCircleOutline as InfoIcon } from '@vicons/ionicons5'
import {
  ExpandLessOutlined as CollapseIcon,
  DragIndicatorSharp as DragIcon,
  ExpandMoreOutlined as ExpandIcon,
} from '@vicons/material'

import { upperFirst } from '@/helpers/strings'
import { VideoQuote } from '@/models/videoQuotes'
import { getInvalidReason } from '@/services/videoQuotes'
import { stripQuotes } from '@/services/videoQuotes'
import { useUserStore } from '@/stores/user'
import { useVideoQuoteAudiosStore } from '@/stores/videoQuoteAudios'
import { useVideoQuotesStore } from '@/stores/videoQuotes'

const videoQuotesStore = useVideoQuotesStore()
const videoQuoteAudiosStore = useVideoQuoteAudiosStore()
const userStore = useUserStore()

const dragHandleRefs = ref<HTMLElement[]>()
const nInputRefs = ref<(typeof NInput)[]>()

const dropdownOptions: DropdownMixedOption[] = [
  {
    label: 'Create a quote above',
    key: 'createAbove',
  },
  {
    label: 'Create a quote below',
    key: 'createBelow',
  },
  {
    label: 'Delete',
    key: 'delete',
  },
]

const props = defineProps<{
  videoId: string
  quote: VideoQuote
  index: number
  dragging: boolean
}>()

const emit = defineEmits<{
  (name: 'select', quoteId: string, focus: boolean): void
  (name: 'drag-handle-mousedown', event: MouseEvent): void
}>()

const data = reactive({
  hasQuoteChanged: false,
  dropdownVisible: {} as Record<string, boolean>,
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
  async toggleHTML(quote: VideoQuote) {
    videoQuotesStore.actions.setVideoQuoteProperty(quote.id!, 'isHtmlEnabled', !quote.isHtmlEnabled)
    await videoQuotesStore.actions.updateVideoQuote(quote.id!)
  },

  handleInputFocus(quoteId: string) {
    emit('select', quoteId, false)
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

  // when pasting, check if the clipboard contains text that is multiple paragraphs long, if so,
  // split it into multiple quotes, otherwise paste it into the current quote
  async handleInputPaste(e: ClipboardEvent) {
    const textarea = e.target as HTMLTextAreaElement
    const selectedText = window.getSelection()?.toString()

    if (textarea.value.length > 0 && selectedText !== textarea.value) return

    const clipboardData = e.clipboardData
    if (!clipboardData) return

    const pastedText = clipboardData.getData('text')
    if (!pastedText) return

    const paragraphs = pastedText.split('\n').filter((p) => p.trim().length > 0)
    if (paragraphs.length === 1) return

    e.preventDefault()

    const preparedParagraphs = paragraphs.map((p) => p.trim()).map(stripQuotes)
    const userId = userStore.state.user!.id
    const activeQuote = videoQuotesStore.getters.activeVideoQuote!
    const activeQuoteIndex = videoQuotesStore.getters.activeVideoQuoteIndex
    const isLastQuote = activeQuoteIndex === videoQuotesStore.state.videoQuotes.length - 1

    const newQuotes = preparedParagraphs.map((paragraph, i) => ({
      ...new VideoQuote(),
      videoId: props.videoId,
      userId,
      content: paragraph,
      position: activeQuote.position + i + 1,
    }))
    let targetIndex = activeQuoteIndex + 1

    if (isLastQuote) {
      targetIndex = activeQuoteIndex
    } else {
      videoQuotesStore.actions.setVideoQuoteProperty(
        activeQuote.id!,
        'content',
        preparedParagraphs[0],
      )

      newQuotes.splice(0, 1)
    }

    const createdQuotes = await videoQuotesStore.actions.insertVideoQuotesAtIndex(
      newQuotes,
      targetIndex,
    )

    // Select last inserted quote
    if (!isLastQuote) {
      const lastCreatedQuote = createdQuotes[createdQuotes.length - 1]
      if (lastCreatedQuote) {
        emit('select', lastCreatedQuote.id!, true)
      }
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

  async deleteQuote(quoteId: string) {
    await videoQuotesStore.actions.deleteVideoQuote(quoteId)
    videoQuoteAudiosStore.actions.visuallyRemoveVideoQuoteAudios(quoteId)
  },

  async handleSelect(key: string, index: number) {
    const quote = videoQuotesStore.state.videoQuotes[index]
    const userId = userStore.state.user!.id
    if (key === 'createAbove') {
      await videoQuotesStore.actions.createNewVideoAtIndex(index, userId, props.videoId)
    } else if (key === 'createBelow') {
      await videoQuotesStore.actions.createNewVideoAtIndex(index + 1, userId, props.videoId)
    } else if (key === 'delete') {
      await methods.deleteQuote(quote.id!)
    }
  },
}
</script>

<template>
  <div
    class="quote-row"
    :data-id="quote.id"
    :class="{
      active: quote.id === videoQuotesStore.state.activeVideoQuoteId,
      dragging: props.dragging,
      'quote-new': index === videoQuotesStore.state.videoQuotes.length - 1,
      invalid: getters.getInvalidReason.value(quote),
    }"
  >
    <div class="text-content-wrapper" @click="emit('select', quote.id!, false)">
      <n-input
        v-if="
          quote.id === videoQuotesStore.state.activeVideoQuoteId ||
          index === videoQuotesStore.state.videoQuotes.length - 1
        "
        ref="nInputRefs"
        :value="quote.content"
        class="text-content-input"
        type="textarea"
        :data-id="quote.id"
        :autosize="{ minRows: 1, maxRows: 3 }"
        @focus="methods.handleInputFocus(quote.id!)"
        @keydown="methods.handleInputKeydown"
        @paste="methods.handleInputPaste"
        @input="methods.handleInputEvent"
        @blur="methods.handleInputBlur(quote.id!)"
        @update:value="videoQuotesStore.actions.setVideoQuoteProperty(quote.id!, 'content', $event)"
      />
      <div v-else class="text-content">
        {{ quote.content || '&nbsp;' }}
      </div>

      <template v-if="index !== videoQuotesStore.state.videoQuotes.length - 1">
        <n-tooltip v-if="getters.getInvalidReason.value(quote)" trigger="hover" placement="bottom">
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

        <n-dropdown
          v-model:show="data.dropdownVisible[index]"
          :options="dropdownOptions"
          placement="bottom-end"
          trigger="click"
          :animated="false"
          @select="(key) => methods.handleSelect(key, index)"
        >
          <n-icon size="16" class="dropdown-icon" @click.stop>
            <CollapseIcon v-if="data.dropdownVisible[index]" />
            <ExpandIcon v-else />
          </n-icon>
        </n-dropdown>

        <div
          v-if="videoQuotesStore.state.videoQuotes.length > 1"
          ref="dragHandleRefs"
          class="drag-handle"
          @mousedown="emit('drag-handle-mousedown', $event)"
          @click.stop
        >
          <drag-icon />
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped lang="scss">
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

.dropdown-icon {
  margin: 0 0 0 5px;
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
</style>
