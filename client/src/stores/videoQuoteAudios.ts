import { acceptHMRUpdate, defineStore } from 'pinia'
import { computed, reactive } from 'vue'

import { PLAYHT_DEFAULT_SPEED } from '@/config/constants'
import { videoToCreateVideoQuoteAudioDTO } from '@/model-functions/videoQuoteAudio'
import { VideoQuoteAudio } from '@/models/videoQuoteAudios'
import {
  createVideoQuoteAudio,
  removeVideoQuoteAudio,
  textToSpeech,
} from '@/services/api/videoQuoteAudios'
import { selectVideoQuoteAudioId } from '@/services/api/videoQuotes'

import { PLAYHT_DEFAULT_VOICE } from './../config/constants'

export interface State {
  videoQuoteAudios: VideoQuoteAudio[]
  activeVideoQuoteId: string
}

const createState = (): State => ({
  videoQuoteAudios: [],
  activeVideoQuoteId: '',
})

export const useVideoQuoteAudiosStore = defineStore('videoQuoteAudiosStore', () => {
  const state: State = reactive<State>(createState())

  const getters = {
    activeVideoQuoteAudios: computed(() => {
      return state.videoQuoteAudios.filter((a) => a.videoQuoteId === state.activeVideoQuoteId)
    }),
  }

  const actions = {
    async selectAudioId(id: string) {
      await selectVideoQuoteAudioId(state.activeVideoQuoteId, id)
    },
    setActiveVideoQuoteId(id: string) {
      state.activeVideoQuoteId = id
    },
    setVideoQuoteAudio(data: VideoQuoteAudio) {
      const index = state.videoQuoteAudios.findIndex((x) => x.id === data.id)
      if (index === -1) throw new Error(`Video quote audio ${data.id} not found`)

      state.videoQuoteAudios[index] = data
    },
    addVideoQuoteAudio(audio: VideoQuoteAudio) {
      state.videoQuoteAudios.unshift(audio)
    },
    setVideoQuoteAudios(audios: VideoQuoteAudio[]) {
      state.videoQuoteAudios = audios
    },
    setVideoQuoteAudioProperty<K extends keyof VideoQuoteAudio>(
      videoQuoteAudioId: string,
      property: K,
      value: VideoQuoteAudio[K] | null | undefined,
    ) {
      const index = state.videoQuoteAudios.findIndex((x) => x.id === videoQuoteAudioId)
      if (index === -1) throw new Error(`Video quote ${videoQuoteAudioId} not found`)

      if (value === null || value === undefined) {
        value = new VideoQuoteAudio()[property]
      }
      state.videoQuoteAudios[index][property] = value
    },
    async removeVideoQuoteAudio(audioId: string) {
      state.videoQuoteAudios = state.videoQuoteAudios.filter((a) => a.id !== audioId)
      // if (state.selectedAudioId === audioId) {
      //   if (state.videoQuoteAudios.length) {
      //     actions.selectAudioId(state.videoQuoteAudios[0].id!)
      //   }
      // }
      await removeVideoQuoteAudio(audioId)
    },
    async visuallyRemoveVideoQuoteAudios(videoQuoteId: string) {
      state.videoQuoteAudios = state.videoQuoteAudios.filter((a) => a.videoQuoteId !== videoQuoteId)
      // const ids = state.videoQuoteAudios.map((a) => a.id)
      // if (ids.includes(state.selectedAudioId)) {
      //   actions.setSelectedAudioId('')
      // }
    },
    async createVideoQuoteAudio(audio: VideoQuoteAudio) {
      if (audio.id) throw new Error('Audio is already created')
      if (!audio.userId) throw new Error('User ID is required')
      if (!audio.videoId) throw new Error('Video ID is required')
      if (!audio.videoQuoteId) throw new Error('Video quote ID is required')

      const createdAudio = await createVideoQuoteAudio(videoToCreateVideoQuoteAudioDTO(audio))
      actions.addVideoQuoteAudio(createdAudio)
    },
    async generateAudio(text: string, voice: string, speed: number, seed: number) {
      const videoQuoteAudio = await textToSpeech({
        videoQuoteId: state.activeVideoQuoteId,
        voice: voice || PLAYHT_DEFAULT_VOICE,
        text,
        seed,
        speed: speed ?? PLAYHT_DEFAULT_SPEED,
      })
      actions.addVideoQuoteAudio(videoQuoteAudio)
      await actions.selectAudioId(videoQuoteAudio.id!)
      return videoQuoteAudio
    },
  }

  const reset = () => {
    const newState: State = createState()
    Object.assign(state, newState)
  }

  return {
    state,
    getters,
    actions,
    reset,
  } as const
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useVideoQuoteAudiosStore, import.meta.hot))
}
