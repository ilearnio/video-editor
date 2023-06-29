export const readFileAsArrayBuffer = (file: File): Promise<ArrayBuffer> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = () => {
      const arrayBuffer = reader.result as ArrayBuffer
      resolve(arrayBuffer)
    }

    reader.onerror = () => {
      reject(new Error('Error reading file.'))
    }

    reader.readAsArrayBuffer(file)
  })
}
