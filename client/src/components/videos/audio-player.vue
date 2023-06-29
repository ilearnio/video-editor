<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'

import { CloudDownloadOutline as CloudDownload, PauseCircle, PlayCircle } from '@vicons/ionicons5'

import { formatDuration } from '@/helpers/duration'

const progressBarRef = ref<HTMLDivElement>()
const scrubberRef = ref<HTMLDivElement>()

const props = defineProps<{
  name: string
  src: string
  subtext?: string
  duration: number // in seconds
}>()

const data = reactive({
  audio: new Audio(),
  isPlaying: ref(false),
  currentTime: ref(0),
  isScrubbing: false,
})

const methods = {
  togglePlayback() {
    if (data.isPlaying) {
      data.audio.pause()
    } else {
      data.audio.currentTime = data.currentTime
      data.audio.play()
    }
    data.isPlaying = !data.isPlaying
  },

  seek(event: MouseEvent) {
    if (!data.isScrubbing) {
      const rect = progressBarRef.value!.getBoundingClientRect()
      const seekTime = ((event.clientX - rect.left) / rect.width) * data.audio.duration
      data.audio.currentTime = seekTime
    }
  },

  startScrubbing() {
    data.isScrubbing = true
    window.addEventListener('mousemove', methods.scrub)
    window.addEventListener('mouseup', methods.stopScrubbing)
  },

  stopScrubbing() {
    data.isScrubbing = false
    window.removeEventListener('mousemove', methods.scrub)
    window.removeEventListener('mouseup', methods.stopScrubbing)
  },

  scrub(event: MouseEvent) {
    if (!data.isScrubbing) return
    const rect = progressBarRef.value!.getBoundingClientRect()
    const scrubTime = ((event.clientX - rect.left) / rect.width) * data.audio.duration
    data.audio.currentTime = scrubTime
  },

  handleTimeUpdate() {
    data.currentTime = data.audio.currentTime
  },
}

const getters = {
  progressWidth: computed(() => {
    if (data.currentTime >= data.audio.duration) {
      return '100%'
    }
    const percent = (data.currentTime / data.audio.duration) * 100
    return `${percent.toFixed(2)}%`
  }),

  scrubberPosition: computed(() => {
    const progressBarWidth = progressBarRef.value?.clientWidth || 0
    const percent = (data.currentTime / data.audio.duration) * 100
    const scrubberWidth = scrubberRef.value?.offsetWidth || 0
    const scrubberMaxPosition = (100 * progressBarWidth) / 100 - scrubberWidth
    const scrubberPosition = Math.min((percent * progressBarWidth) / 100, scrubberMaxPosition)
    return `${scrubberPosition}px`
  }),

  playbackIcon: computed(() => {
    return data.isPlaying ? PauseCircle : PlayCircle
  }),
}

onMounted(() => {
  data.audio.addEventListener('timeupdate', methods.handleTimeUpdate)

  data.audio.addEventListener('ended', () => {
    data.isPlaying = false
    setTimeout(() => {
      data.currentTime = 0
    }, 200)
  })

  data.audio.src = props.src
})

onUnmounted(() => {
  data.audio.removeEventListener('timeupdate', methods.handleTimeUpdate)
})
</script>

<template>
  <div class="audio-player">
    <div class="audio-controls">
      <div class="audio-column-left">
        <slot name="left" />
      </div>
      <button class="play-button" @click="methods.togglePlayback">
        <n-icon size="40" :component="getters.playbackIcon.value" />
      </button>
      <div class="audio-column-middle">
        <span class="audio-name" :title="props.name">{{ props.name }}</span>
        <div ref="progressBarRef" class="progress-bar" @click="methods.seek">
          <div class="progress" :style="{ width: getters.progressWidth.value }"></div>
          <div
            ref="scrubberRef"
            class="scrubber"
            :style="{ left: getters.scrubberPosition.value }"
            @mousedown="methods.startScrubbing"
          ></div>
        </div>
        <div class="audio-info">
          <div class="audio-subtext">{{ props.subtext }}</div>
          <span class="audio-duration">{{ formatDuration(props.duration) }}</span>
        </div>
      </div>
      <div class="audio-column-right">
        <a
          :href="props.src"
          :download="props.name"
          class="download-link"
          title="Download"
          target="_blank"
        >
          <CloudDownload class="download-icon" />
        </a>
        <slot name="right" />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.audio-player {
  background: #fff;
  padding: 5px 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  min-width: 0;
}

.audio-controls {
  display: flex;
  align-items: center;
  margin-right: 5px;
}

.play-button {
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  padding: 0;
  margin-right: 8px;
  color: #18a058;
}

.audio-column-middle {
  flex: 1;
  min-width: 0;
}

.audio-name {
  line-height: 1;
  margin-bottom: 6px;
  min-width: 0;
  max-width: 90%;
  white-space: nowrap;
  text-overflow: ellipsis;
  display: block;
  overflow: hidden;
}

.progress-bar {
  height: 5px;
  background-color: #ddd;
  cursor: pointer;
  position: relative;
  user-select: none;

  .progress {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    background-color: #18a058;
    width: 0;
  }
}

.scrubber {
  position: absolute;
  top: -3px;
  left: 0;
  width: 11px;
  height: 11px;
  background-color: #18a058;
  border-radius: 50%;
  cursor: pointer;
  transition: transform 100ms ease-in-out;

  &:hover {
    transform: scale(1.2);
  }
}

.audio-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 2px;
}

.audio-subtext {
  font-size: 80%;
  color: #777;
  min-width: 0;
  text-overflow: ellipsis;
  overflow: hidden;
  max-width: 100px;
  white-space: nowrap;
}

.audio-duration {
  font-size: 80%;
  color: #777;
}

.audio-column-left {
  display: flex;
  align-items: center;
}

.audio-column-right {
  display: flex;
  align-items: center;
}

.download-link {
  display: block;
  padding: 3px;
  margin-left: 12px;
  color: #888;

  &:hover .download-icon {
    color: #111;
  }
}

.download-icon {
  width: 20px;
  display: block;
  color: #888;
}
</style>
