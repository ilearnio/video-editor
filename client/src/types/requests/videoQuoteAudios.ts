export interface TextToSpeechRequest {
  videoQuoteId: string
  voice: string
  text: string
  seed: number | null
  speed: number
}
