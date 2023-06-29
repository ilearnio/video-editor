import { type User } from '@/models/users'

import pb from '../pocketbase'

export const authenticateUser = async (
  email: string,
  password: string,
): Promise<{
  token: string
  user: User
}> => {
  const authData = await pb.collection('users').authWithPassword(email, password)
  if (!pb.authStore.isValid) {
    throw new Error('Invalid credentials')
  }
  const user: User = {
    id: authData.record.id,
    username: authData.record.username,
    email: authData.record.email,
    name: authData.record.name,
    avatar: authData.record.avatar,
    collectionId: authData.record.collectionId,
    collectionName: authData.record.collectionName,
    emailVisibility: authData.record.emailVisibility,
    verified: authData.record.verified,
    created: authData.record.created,
    updated: authData.record.updated,
  }

  return {
    token: authData.token,
    user,
  }
}
