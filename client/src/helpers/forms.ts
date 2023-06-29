// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const objectToFormData = (obj: { [key: string]: any }): FormData => {
  const formData = new FormData()

  for (const key in obj) {
    // eslint-disable-next-line no-prototype-builtins
    if (obj.hasOwnProperty(key)) {
      formData.append(key, obj[key])
    }
  }

  return formData
}
