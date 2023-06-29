export const formatDuration = (durationSec: number): string => {
  const hours = Math.floor(durationSec / 3600)
  const minutes = Math.floor((durationSec / 60) % 60)
    .toString()
    .padStart(2, '0')
  const seconds = Math.floor(durationSec % 60)
    .toString()
    .padStart(2, '0')
  const deciseconds = Math.floor((durationSec % 1) * 10)
    .toString()
    .padStart(1, '0')
  return `${hours > 0 ? `${hours}:` : ''}${minutes}:${seconds}.${deciseconds}`
}
