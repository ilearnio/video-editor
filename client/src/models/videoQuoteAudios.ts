export class VideoQuoteAudio {
  id?: string
  created?: string
  updated?: string

  constructor(
    public userId = '',
    public videoId = '',
    public videoQuoteId = '',
    public isHtmlEnabled = false,
    public text = '',
    public voice = '',
    public speed = 0,
    public size = 0,
    public audioFile = '' as string | File,
    public duration = 0, // ms
  ) {}
}

export type CreateVideoQuoteAudioDTO = Omit<
  VideoQuoteAudio,
  'id' | 'created' | 'collectionId' | 'collectionName'
>

export type UpdateVideoQuoteAudioDTO = CreateVideoQuoteAudioDTO
