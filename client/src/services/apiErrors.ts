import { type NotificationApiInjection } from 'naive-ui/es/notification/src/NotificationProvider'

import { isPocketbaseResponseError } from '@/types/responses/generic'

export const handleApiError = (err: unknown, notification: NotificationApiInjection) => {
  if (isPocketbaseResponseError(err)) {
    notification.error({ title: 'Error', content: err.message })
    return
  }
  console.error(err)
  notification.error({
    title: 'Error',
    content: 'Unknown server error',
  })
}
