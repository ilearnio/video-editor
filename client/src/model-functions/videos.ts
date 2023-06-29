import { type CreateVideoDTO, type UpdateVideoDTO, Video } from '@/models/videos'

export const videoToCreateVideoDTO = (video: Video): CreateVideoDTO => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { id, created, updated, ...dto } = video
  return dto
}

export const videoToUpdateVideoDTO = (video: Video): UpdateVideoDTO => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { id, created, updated, ...dto } = video
  return dto
}
