<script setup lang="ts">
import debounce from 'just-debounce'
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import { useNotification } from 'naive-ui'
import type { FileInfo } from 'naive-ui/es/upload/src/interface'
import { type Ref, computed, onBeforeMount, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import VideoAudioCard from '@/components/videos/video-audio-card.vue'
import VideoQuotesEditor from '@/components/videos/video-quotes-editor.vue'
import { calculateAudioDuration } from '@/helpers/audios'
import type { Video } from '@/models/videos'
import { getVideo } from '@/services/api/videos'
import { handleApiError } from '@/services/apiErrors'
import { getFileUrl } from '@/services/pocketbase'
import { useVideoStore } from '@/stores/video'
import { useVideoQuoteAudiosStore } from '@/stores/videoQuoteAudios'
import { useVideoQuotesStore } from '@/stores/videoQuotes'

const route = useRoute()
const videoStore = useVideoStore()
const videoQuotesStore = useVideoQuotesStore()
const videoQuoteAudiosStore = useVideoQuoteAudiosStore()
const notification = useNotification()

const formRef = ref<FormInst>()
const formValueRef: Ref<Video> = ref<Video>(videoStore.state.video)

const apiUrl = import.meta.env.VITE_API_BASE_URL

const data = reactive({
  isLoading: false,
  showImagePreview: false,
  previewImageUrl: '',
})

const getters = {
  isEditMode: computed<boolean>(() => !!route.params.id),
  rules: computed<FormRules>(() => ({
    title: {
      required: true,
      trigger: ['input'],
      message: 'Title is required',
    },
  })),
  getFileListItems: computed(() => (property: keyof Video): FileInfo[] => {
    if (!videoStore.state.video[property]) return []

    const { video } = videoStore.state
    const isFile = video[property] instanceof File
    const id = '1'
    const name = isFile ? (video[property] as File).name : (video[property] as string)
    const status = 'finished'
    const url = isFile ? '' : getFileUrl('videos', video.id!, name)
    const thumbnailUrl = isFile ? '' : getFileUrl('videos', video.id!, name, '188x188')
    const file = isFile ? (video[property] as File) : null

    return [{ id, name, status, url, thumbnailUrl, file }]
  }),
  exportShotcutMltURL: computed(() => {
    return `${apiUrl}/api/videos/${videoStore.state.video.id}/exportShotcutProjectMlt`
  }),
  exportFullShotcutProjectURL: computed(() => {
    return `${apiUrl}/api/videos/${videoStore.state.video.id}/exportShotcutProjectFull`
  }),
}

const methods = {
  handlePreview(file: UploadFileInfo) {
    data.previewImageUrl = file.url as string
    data.showImagePreview = true
  },
  handleFormUpdate: debounce(async (video: Video) => {
    videoStore.actions.setVideo(video)
    await methods.save()
  }, 500),
  async fetchVideo() {
    const video = await getVideo(route.params.id as string)
    videoStore.actions.setVideo(video)
    formValueRef.value = video
  },
  async handleBackgroundAudioFileChange(file: string | File) {
    if (file instanceof File) {
      const durationMs = await calculateAudioDuration(file)
      videoStore.actions.setVideoProperty('backgroundAudioDuration', durationMs)
    }
  },
  async save() {
    data.isLoading = true
    try {
      await videoStore.actions.updateVideo()
    } catch (err) {
      handleApiError(err, notification)
    } finally {
      data.isLoading = false
    }
  },
  submit() {
    formRef.value?.validate(async (errors) => {
      if (errors) return
      await methods.save()
    })
  },
}

watch(formValueRef, methods.handleFormUpdate, { deep: true })
watch(() => formValueRef.value.backgroundAudioFile, methods.handleBackgroundAudioFileChange)

onBeforeMount(async () => {
  videoStore.reset()
  videoQuotesStore.reset()
  videoQuoteAudiosStore.reset()

  await methods.fetchVideo()
})
</script>

<template>
  <template v-if="videoStore.state.video.id">
    <n-form ref="formRef" :model="formValueRef" :rules="getters.rules.value">
      <n-space vertical>
        <n-form-item label="Title" path="title">
          <n-input
            v-model:value="formValueRef.title"
            placeholder=""
            :input-props="{ autocomplete: 'off' }"
          />
        </n-form-item>
        <n-form-item path="title" label-style="position: relative">
          <template #label>
            Heading
            <n-switch v-model:value="formValueRef.headingIsHTML" class="heading-html-switch">
              <template #checked>HTML</template>
              <template #unchecked>HTML</template>
            </n-switch>
          </template>
          <n-input
            v-model:value="formValueRef.heading"
            type="textarea"
            placeholder=""
            :autosize="{ minRows: 2, maxRows: 4 }"
            :input-props="{ autocomplete: 'off' }"
          />
        </n-form-item>
        <n-form-item label="Intro Image" path="uploadValue">
          <n-upload
            ref="upload"
            accept="image/jpeg,image/png"
            :list-type="formValueRef.introImageFile ? 'image-card' : 'text'"
            :file-list="getters.getFileListItems.value('introImageFile')"
            :show-trigger="!formValueRef.introImageFile"
            @change="formValueRef.introImageFile = $event.file.file!"
            @preview="methods.handlePreview"
          >
            <n-button v-if="!formValueRef.introImageFile">Upload file</n-button>
          </n-upload>
        </n-form-item>
        <n-form-item label="Outro Image" path="uploadValue">
          <n-upload
            ref="upload"
            accept="image/jpeg,image/png"
            :list-type="formValueRef.introImageFile ? 'image-card' : 'text'"
            :file-list="getters.getFileListItems.value('outroImageFile')"
            :show-trigger="!formValueRef.outroImageFile"
            @change="formValueRef.outroImageFile = $event.file.file!"
            @preview="methods.handlePreview"
          >
            <n-button>Upload file</n-button>
          </n-upload>
        </n-form-item>
        <n-form-item label="Outro Overlay Image" path="uploadValue">
          <n-upload
            ref="upload"
            accept="image/jpeg,image/png"
            :list-type="formValueRef.introImageFile ? 'image-card' : 'text'"
            :file-list="getters.getFileListItems.value('outroOverlayImageFile')"
            :show-trigger="!formValueRef.outroOverlayImageFile"
            @change="formValueRef.outroOverlayImageFile = $event.file.file!"
            @preview="methods.handlePreview"
          >
            <n-button>Upload file</n-button>
          </n-upload>
        </n-form-item>
        <n-form-item label="Background Image" path="uploadValue">
          <n-upload
            ref="upload"
            accept="image/jpeg,image/png"
            :list-type="formValueRef.introImageFile ? 'image-card' : 'text'"
            :file-list="getters.getFileListItems.value('backgroundImageFile')"
            :show-trigger="!formValueRef.backgroundImageFile"
            @change="formValueRef.backgroundImageFile = $event.file.file!"
            @preview="methods.handlePreview"
          >
            <n-button>Upload file</n-button>
          </n-upload>
        </n-form-item>
        <n-form-item label="Background Audio" path="uploadValue">
          <n-upload
            v-if="!formValueRef.backgroundAudioFile"
            ref="upload"
            accept="audio/wav,audio/mpeg,audio/mp3,audio/ogg,audio/flac,audio/aac"
            :file-list="getters.getFileListItems.value('backgroundAudioFile')"
            @change="formValueRef.backgroundAudioFile = $event.file.file!"
            @preview="methods.handlePreview"
          >
            <n-button>Upload file</n-button>
          </n-upload>
          <video-audio-card
            v-else
            :video-id="videoStore.state.video.id!"
            :file="formValueRef.backgroundAudioFile"
            :duration="formValueRef.backgroundAudioDuration / 1000"
            @delete="formValueRef.backgroundAudioFile = ''"
          />
        </n-form-item>
        <n-modal v-model:show="data.showImagePreview" preset="card" style="width: 600px">
          <img :src="data.previewImageUrl" style="width: 100%" />
        </n-modal>
      </n-space>
    </n-form>

    <video-quotes-editor v-if="videoStore.state.video.id" :video-id="videoStore.state.video.id" />

    <n-space justify="end">
      <a :href="getters.exportShotcutMltURL.value" download>
        <n-button>Export Shotcut .mlt</n-button>
      </a>
      <a :href="getters.exportFullShotcutProjectURL.value" download>
        <n-button>Export Full Shotcut Project</n-button>
      </a>
    </n-space>
  </template>
</template>

<style lang="scss" scoped>
.heading-html-switch {
  position: absolute;
  right: 0;
}
</style>
