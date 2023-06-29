export class Video {
  id?: string
  created?: string
  updated?: string

  constructor(
    public userId = '',
    public type = 'quotes',
    public status = 'draft' as 'draft' | 'completed',
    public title = '',
    public heading = '',
    public headingIsHTML = false,
    public introImageFile = '' as string | File,
    public outroImageFile = '' as string | File,
    public outroOverlayImageFile = '' as string | File,
    public backgroundImageFile = '' as string | File,
    public backgroundAudioFile = '' as string | File,
    public backgroundAudioVolume = '',
    public backgroundAudioDuration = 0,
    public gapBetweenQuotes = 2000,
  ) {}
}

export type CreateVideoDTO = Omit<Video, 'id' | 'created' | 'updated'>

export type UpdateVideoDTO = CreateVideoDTO
