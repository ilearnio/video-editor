import {
  type CreateVideoQuoteAudioDTO,
  type UpdateVideoQuoteAudioDTO,
  VideoQuoteAudio,
} from '@/models/videoQuoteAudios'

export const videoToCreateVideoQuoteAudioDTO = (
  video: VideoQuoteAudio,
): CreateVideoQuoteAudioDTO => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { id, created, ...dto } = video
  return dto
}

export const videoToUpdateVideoQuoteAudioDTO = (
  video: VideoQuoteAudio,
): UpdateVideoQuoteAudioDTO => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { id, created, ...dto } = video
  return dto
}
