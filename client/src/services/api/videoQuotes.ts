import { ListResult } from 'pocketbase'

import { objectToFormData } from '@/helpers/forms'
import type { CreateVideoQuoteDTO, UpdateVideoQuoteDTO, VideoQuote } from '@/models/videoQuotes'
import pb from '@/services/pocketbase'

const collection = pb.collection('videoQuotes')

export const getVideoQuotes = async (
  userId: string,
  videoId: string,
): Promise<ListResult<VideoQuote>> => {
  const sort: keyof VideoQuote = 'position'
  const records = await collection.getList<VideoQuote>(1, 999, {
    sort,
    filter: `userId = '${userId}' && videoId = '${videoId}'`,
  })
  return records
}

export const getVideoQuote = async (id: string): Promise<VideoQuote> => {
  const record = await collection.getOne<VideoQuote>(id)
  return record
}

export const createVideoQuote = async (data: CreateVideoQuoteDTO): Promise<VideoQuote> => {
  const record = await collection.create<VideoQuote>(objectToFormData(data), {
    $autoCancel: false,
  })
  return record
}

export const updateVideoQuote = async (
  id: string,
  data: Partial<UpdateVideoQuoteDTO>,
): Promise<VideoQuote> => {
  const record = await collection.update<VideoQuote>(id, objectToFormData(data), {
    $autoCancel: false,
  })
  return record
}

export const updateVideoQuotes = async (quotes: VideoQuote[]): Promise<VideoQuote[]> => {
  return Promise.all(
    quotes.map(async (quote) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { id, created, updated, ...dto } = quote
      return await updateVideoQuote(quote.id!, dto)
    }),
  )
}

export const selectVideoQuoteAudioId = async (
  videoQuoteId: string,
  videoQuoteAudioId: string,
): Promise<void> => {
  await updateVideoQuote(videoQuoteId, { selectedAudioId: videoQuoteAudioId })
}

export const deleteVideoQuote = async (id: string): Promise<void> => {
  await collection.delete(id)
}
