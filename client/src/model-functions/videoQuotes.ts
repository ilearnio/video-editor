import {
  type CreateVideoQuoteDTO,
  type UpdateVideoQuoteDTO,
  VideoQuote,
} from '@/models/videoQuotes'

export const videoToCreateVideoQuoteDTO = (record: VideoQuote): CreateVideoQuoteDTO => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { id, created, updated, ...dto } = record
  return dto
}

export const videoToUpdateVideoQuoteDTO = (record: VideoQuote): UpdateVideoQuoteDTO => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { id, created, updated, ...dto } = record
  return dto
}
