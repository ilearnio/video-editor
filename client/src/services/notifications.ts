export const showNotification = (type: 'error' | 'success', message: string) => {
  window.alert(`${type}: ${message}`)
}
