import { readFileAsArrayBuffer } from './files'

// Calculates the duration of an audio file, supporting all audio formats
export const calculateAudioDuration = async (file: File): Promise<number> => {
  const arrayBuffer = await readFileAsArrayBuffer(file)
  const audioContext = new window.AudioContext()

  return new Promise((resolve, reject) => {
    audioContext.decodeAudioData(
      arrayBuffer,
      (buffer) => {
        resolve(Math.round(buffer.duration * 1000)) // ms
      },
      (error) => {
        reject(error)
      },
    )
  })
}
