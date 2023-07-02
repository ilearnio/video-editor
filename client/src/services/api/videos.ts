import { ListResult, type RecordListQueryParams } from 'pocketbase'

import { objectToFormData } from '@/helpers/forms'
import type { CreateVideoDTO, UpdateVideoDTO, Video } from '@/models/videos'
import pb from '@/services/pocketbase'

const collection = pb.collection('videos')

export const getVideoList = async (
  queryParams?: RecordListQueryParams,
): Promise<ListResult<Video>> => {
  const sort: keyof Video = 'created'
  const params = { sort: `-${sort}`, ...queryParams }
  const result = await collection.getList<Video>(params.page ?? 1, params.perPage ?? 50, params)
  return result
}

export const getVideo = async (id: string): Promise<Video> => {
  const record = await collection.getOne<Video>(id)
  return record
}

export const createVideo = async (data: CreateVideoDTO): Promise<Video> => {
  const record = await collection.create<Video>(objectToFormData(data), {
    $autoCancel: false,
  })
  return record
}

export const updateVideo = async (id: string, data: UpdateVideoDTO): Promise<Video> => {
  const record = await collection.update<Video>(id, objectToFormData(data), {
    $autoCancel: false,
  })
  return record
}

export const deleteVideo = async (id: string) => {
  return await collection.delete(id)
}
