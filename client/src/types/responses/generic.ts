import { ClientResponseError } from 'pocketbase'

export interface PocketbaseRequestError extends Error {
  originalError: Error
  status: number
}

export interface PocketbaseGenericErrorResponse {
  code: number
  message: string
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const isPocketbaseResponseError = (err: any): err is ClientResponseError => {
  return 'originalError' in err && 'status' in err
}
