export class VideoQuote {
  id?: string
  created?: string
  updated?: string

  constructor(
    public userId = '',
    public videoId = '',
    public geometry = '', // 'x, y, width, height'
    public isHtmlEnabled = false,
    public content = '',
    public position = -1,
    public selectedAudioId = '',
  ) {}
}

export type CreateVideoQuoteDTO = Omit<
  VideoQuote,
  'id' | 'created' | 'updated' | 'collectionId' | 'collectionName'
>

export type UpdateVideoQuoteDTO = CreateVideoQuoteDTO
