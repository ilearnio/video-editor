import { objectToFormData } from '@/helpers/forms'
import type {
  CreateVideoQuoteAudioDTO,
  UpdateVideoQuoteAudioDTO,
  VideoQuoteAudio,
} from '@/models/videoQuoteAudios'
import api from '@/services/api'
import pb from '@/services/pocketbase'
import { type TextToSpeechRequest } from '@/types/requests/videoQuoteAudios'

export const collection = pb.collection('videoQuoteAudios')

export const getVideoQuoteAudios = async (videoId: string): Promise<VideoQuoteAudio[]> => {
  const sort: keyof VideoQuoteAudio = 'created'
  const records = await collection.getFullList<VideoQuoteAudio>({
    sort: `-${sort}`,
    filter: `videoId = '${videoId}'`,
  })
  return records
}

export const getVideoQuoteAudio = async (id: string): Promise<VideoQuoteAudio> => {
  const record = await collection.getOne<VideoQuoteAudio>(id)
  return record
}

export const createVideoQuoteAudio = async (
  data: CreateVideoQuoteAudioDTO,
): Promise<VideoQuoteAudio> => {
  const record = await collection.create<VideoQuoteAudio>(objectToFormData(data), {
    $autoCancel: false,
  })
  return record
}

export const updateVideoQuoteAudio = async (
  id: string,
  data: UpdateVideoQuoteAudioDTO,
): Promise<VideoQuoteAudio> => {
  const record = await collection.update<VideoQuoteAudio>(id, objectToFormData(data), {
    $autoCancel: false,
  })
  return record
}

export const updateVideoQuoteAudios = async (
  videoQuoteAudios: VideoQuoteAudio[],
): Promise<VideoQuoteAudio[]> => {
  return Promise.all(
    videoQuoteAudios.map(async (videoQuoteAudio) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { id, created, ...dto } = videoQuoteAudio
      return await updateVideoQuoteAudio(videoQuoteAudio.id!, dto)
    }),
  )
}

export const textToSpeech = async (req: TextToSpeechRequest): Promise<VideoQuoteAudio> => {
  const res = await api.post('/api/videoQuoteAudio/generateFromText', req)
  return res.data
}

export const deleteVideoQuoteAudio = async (id: string): Promise<void> => {
  await collection.delete(id)
}
