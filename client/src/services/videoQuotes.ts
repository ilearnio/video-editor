/**
 * Service functions for video quotes and video quote audios
 */
import { PLAYHT_DEFAULT_SPEED, PLAYHT_DEFAULT_VOICE } from '@/config/constants'
import { stripTags } from '@/helpers/html'
import { VideoQuoteAudio } from '@/models/videoQuoteAudios'
import { VideoQuote } from '@/models/videoQuotes'

export enum InvalidVideoQuoteAudioReason {
  ContentMismatch = 'The text content of the phrase was updated',
  VoiceMismatch = 'The voice you used here differs from other voices',
  SpeedMismatch = 'The speed you used here differs from other speeds',
}

export const extractQuotePlainText = (quote: VideoQuote): string => {
  return quote.isHtmlEnabled ? stripTags(quote.content) : quote.content
}

export const extractAudioPlainText = (audio: VideoQuoteAudio): string => {
  return audio.isHtmlEnabled ? stripTags(audio.text) : audio.text
}

export const findSelectedAudios = (
  quotes: VideoQuote[],
  audios: VideoQuoteAudio[],
): VideoQuoteAudio[] => {
  const selectedAudioIds = quotes.map((q) => q.selectedAudioId).filter(Boolean)
  return audios.filter((a) => !!a.id && selectedAudioIds.includes(a.id))
}

export const mostUsedVoice = (quotes: VideoQuote[], audios: VideoQuoteAudio[]): string => {
  const selectedAudios = findSelectedAudios(quotes, audios)
  const voices = selectedAudios.map((a) => a.voice)

  const voiceCountMap = voices.reduce<Record<string, number>>((acc, voice) => {
    acc[voice] = (acc[voice] || 0) + 1
    return acc
  }, {})

  const sortedEntries = Object.entries(voiceCountMap).sort((a, b) => b[1] - a[1])

  return sortedEntries[0]?.[0] || PLAYHT_DEFAULT_VOICE
}

export const mostUsedSpeed = (quotes: VideoQuote[], audios: VideoQuoteAudio[]): number => {
  const selectedAudios = findSelectedAudios(quotes, audios)
  const speeds = selectedAudios.map((a) => a.speed)

  const speedCountMap = speeds.reduce<Record<number, number>>((acc, speed) => {
    acc[speed] = (acc[speed] || 0) + 1
    return acc
  }, {})

  const sortedEntries = Object.entries(speedCountMap).sort((a, b) => b[1] - a[1])
  return Number(sortedEntries[0]?.[0] || PLAYHT_DEFAULT_SPEED)
}

export const getInvalidReason = (
  videoQuotes: VideoQuote[],
  audios: VideoQuoteAudio[],
  audio: VideoQuoteAudio,
) => {
  const videoQuote = videoQuotes.find((q) => q.id === audio.videoQuoteId)
  if (!videoQuote) {
    throw new Error('Video quote not found')
  }

  const videoQuoteText = videoQuote.isHtmlEnabled
    ? stripTags(videoQuote.content)
    : videoQuote.content
  const audioText = extractAudioPlainText(audio)

  if (videoQuoteText !== audioText) {
    return InvalidVideoQuoteAudioReason.ContentMismatch
  }
  if (mostUsedVoice(videoQuotes, audios) !== audio.voice) {
    return InvalidVideoQuoteAudioReason.VoiceMismatch
  }
  if (mostUsedSpeed(videoQuotes, audios) !== audio.speed) {
    return InvalidVideoQuoteAudioReason.SpeedMismatch
  }
}
