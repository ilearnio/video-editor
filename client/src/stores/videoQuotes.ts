import { acceptHMRUpdate, defineStore } from 'pinia'
import { computed, reactive } from 'vue'

import { videoToCreateVideoQuoteDTO } from '@/model-functions/videoQuotes'
import { VideoQuote } from '@/models/videoQuotes'
import {
  createVideoQuote,
  getVideoQuotes,
  removeVideoQuote,
  updateVideoQuote,
  updateVideoQuotes,
} from '@/services/api/videoQuotes'
import { isFulfilled, isRejected } from '@/types/general'

export interface State {
  activeVideoQuoteId: string
  videoQuotes: VideoQuote[]
}

const createState = (): State => ({
  activeVideoQuoteId: '',
  videoQuotes: [],
})

export const useVideoQuotesStore = defineStore('videoQuotesStore', () => {
  const state: State = reactive<State>(createState())

  const getters = {
    activeVideoQuoteIndex: computed(() => {
      return state.videoQuotes.findIndex((q) => q.id === state.activeVideoQuoteId)
    }),
    activeVideoQuote: computed(() => {
      return state.videoQuotes.find((q) => q.id === state.activeVideoQuoteId)
    }),
  }

  const actions = {
    setActiveVideoQuoteId(id: string) {
      state.activeVideoQuoteId = id
    },
    setVideoQuote(data: VideoQuote) {
      const index = state.videoQuotes.findIndex((x) => x.id === data.id)
      if (index === -1) throw new Error(`Video quote ${data.id} not found`)

      state.videoQuotes[index] = data
    },
    addVideoQuote(data: VideoQuote) {
      state.videoQuotes.push(data)
    },
    setVideoQuotes(quotes: VideoQuote[]) {
      state.videoQuotes = quotes
    },
    setVideoQuoteProperty<K extends keyof VideoQuote>(
      videoQuoteId: string,
      property: K,
      value: VideoQuote[K] | null | undefined,
    ) {
      const index = state.videoQuotes.findIndex((x) => x.id === videoQuoteId)
      if (index === -1) throw new Error(`Video quote ${videoQuoteId} not found`)

      if (value === null || value === undefined) {
        const newQuote = new VideoQuote()
        value = newQuote[property]
      }

      state.videoQuotes[index][property] = value
    },
    async getVideoQuotes(userId: string, videoId: string) {
      const response = await getVideoQuotes(userId, videoId)
      actions.setVideoQuotes(response.items)
    },
    async createNewVideoQuote(userId: string, videoId: string, position: number) {
      const data = { ...new VideoQuote(), userId, videoId, position }
      const createdQuote = await createVideoQuote(videoToCreateVideoQuoteDTO(data))
      state.videoQuotes.push(createdQuote)
    },
    async updateVideoQuote(id: string) {
      const quoteIndex = state.videoQuotes.findIndex((q) => q.id === id)
      const quote = state.videoQuotes[quoteIndex]
      if (!quote) throw new Error('Video quote not found')
      if (!id) throw new Error('ID is missing')

      const updatedQuote = await updateVideoQuote(id, videoToCreateVideoQuoteDTO(quote))
      state.videoQuotes[quoteIndex] = {
        ...state.videoQuotes[quoteIndex],
        ...updatedQuote,
      }
    },
    async updateActiveVideoQuote() {
      if (!getters.activeVideoQuote.value?.id) throw new Error('Video quote not found')
      return actions.updateVideoQuote(getters.activeVideoQuote.value.id)
    },
    async updateVideoQuotes(quotes = state.videoQuotes) {
      const updatedQuotes = await updateVideoQuotes(quotes)
      const videoQuotes: VideoQuote[] = JSON.parse(JSON.stringify(state.videoQuotes))
      updatedQuotes.forEach((quote) => {
        const quoteIndex = videoQuotes.findIndex((x) => x.id === quote.id)
        videoQuotes.splice(quoteIndex, 1, quote)
      })
      state.videoQuotes = videoQuotes
    },
    recalcPositions(): boolean {
      let updated = false
      state.videoQuotes.forEach((quote, index) => {
        if (quote.position === index) return
        updated = true
        quote.position = index
      })
      return updated
    },
    async removeVideoQuote(id: string) {
      if (state.activeVideoQuoteId === id) actions.setActiveVideoQuoteId('')

      const index = state.videoQuotes.findIndex((q) => q.id === id)
      if (index === -1) throw new Error(`Video quote ${id} not found`)

      state.videoQuotes.splice(index, 1)

      actions.recalcPositions()
      await removeVideoQuote(id)
      await actions.updateVideoQuotes()
    },
    async insertVideoQuotesAtIndex(quotes: VideoQuote[], index: number): Promise<VideoQuote[]> {
      const initialPosition = state.videoQuotes[index - 1]
        ? state.videoQuotes[index - 1].position + 1
        : 0
      quotes = quotes.map((q, i) => ({ ...q, position: initialPosition + i }))

      const promises = quotes.map((q) => createVideoQuote(videoToCreateVideoQuoteDTO(q)))

      const results = await Promise.allSettled(promises)
      const rejected = results.filter(isRejected)
      if (rejected.length) {
        console.error('Failed to create some of the video quotes')
      }

      const createdQuotes = results.filter(isFulfilled).map((r) => r.value)
      if (createdQuotes.length > 0) {
        state.videoQuotes.splice(index, 0, ...createdQuotes)
        actions.recalcPositions()
        await actions.updateVideoQuotes()
      }

      return createdQuotes
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
  import.meta.hot.accept(acceptHMRUpdate(useVideoQuotesStore, import.meta.hot))
}
