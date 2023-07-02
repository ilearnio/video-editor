import { acceptHMRUpdate, defineStore } from 'pinia'
import { ListResult, type RecordListQueryParams } from 'pocketbase'
import { reactive } from 'vue'

import { videoToCreateVideoDTO, videoToUpdateVideoDTO } from '@/model-functions/videos'
import { Video } from '@/models/videos'
import { createVideo, deleteVideo, getVideoList, updateVideo } from '@/services/api/videos'

let lastQueryParams: RecordListQueryParams | undefined

export interface State {
  list: ListResult<Video>
  video: Video
}

const createState = () => ({
  list: {
    page: -1,
    perPage: -1,
    totalItems: -1,
    totalPages: -1,
    items: [],
  },
  video: new Video(),
})

export const useVideoStore = defineStore('videoStore', () => {
  const state: State = reactive<State>(createState())

  const getters = {}

  const actions = {
    async fetchItems(queryParams?: RecordListQueryParams) {
      lastQueryParams = queryParams
      state.list = await getVideoList(queryParams)
    },
    async deleteItems(ids: string[]) {
      const promises = ids.map((id) => deleteVideo(id))
      await Promise.all(promises)
      await actions.fetchItems(lastQueryParams)
    },
    setVideo(video: Video) {
      state.video = video
    },
    setVideoProperty<K extends keyof Video>(property: K, value: Video[K] | null | undefined) {
      if (value === null || value === undefined) {
        value = new Video()[property]
      }
      state.video[property] = value
    },
    async deleteVideo(id: string) {
      await deleteVideo(id)
    },
    async createVideo() {
      if (state.video.id) throw new Error('Video is already created')
      const data = { ...state.video, title: 'Draft' }
      const createdVideo = await createVideo(videoToCreateVideoDTO(data))
      actions.setVideo(createdVideo)
      return createdVideo
    },
    async updateVideo() {
      if (!state.video.id) {
        throw new Error('Video is not created yet')
      }
      const updatedVideo = await updateVideo(state.video.id, videoToUpdateVideoDTO(state.video))
      actions.setVideo(updatedVideo)
      return updatedVideo
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
  import.meta.hot.accept(acceptHMRUpdate(useVideoStore, import.meta.hot))
}
