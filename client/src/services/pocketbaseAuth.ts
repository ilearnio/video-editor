import { type User } from '@/models/users'

import pb from './pocketbase'

export const pocketbaseParseAuth = (json: string | null): { token: string; model: User } | null => {
  try {
    return json ? JSON.parse(json) : null
  } catch (error) {
    return null
  }
}

export const pocketbaseLogoutUser = () => {
  pb.authStore.clear()
}
