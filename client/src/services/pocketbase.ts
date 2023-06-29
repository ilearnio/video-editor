import PocketBase from 'pocketbase'

// import { useUserStore } from '@/stores/user'

// import { handleApiError } from './apiErrors'

export const baseUrl = import.meta.env.VITE_API_BASE_URL

const pb = new PocketBase(baseUrl)

export const getFileUrl = (
  collectionIdOrName: string,
  recordId: string,
  filename: string,
  thumbnailSize?: string, // '100x300'
) => {
  let url = `${baseUrl}/api/files/${collectionIdOrName}/${recordId}/${filename}`
  if (thumbnailSize) url += `?thumb=${thumbnailSize}`
  return url
}

// pb.afterSend = (response, data) => {
//   const userStore = useUserStore()
//   // 401 Unauthorized; 419 Session expired
//   if ([401, 419].includes(response.status)) {
//     userStore.logoutUser()
//     return
//   }

//   if (response.status >= 400) {
//     handleApiError(data)
//     return
//   }
//   return data
// }

export default pb
