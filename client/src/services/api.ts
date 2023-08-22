import axios from 'axios'

import { useUserStore } from '@/stores/user'

// import { handleApiError } from './apiErrors'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

api.interceptors.request.use((config) => {
  const userStore = useUserStore()

  if (userStore.state.token) {
    config.headers.Authorization = `Bearer ${userStore.state.token}`
  }

  return config
})

api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // const userStore = useUserStore()
    // 401 Unauthorized; 419 Session expired
    if ([401, 419].includes(error.response?.status)) {
      // userStore.logoutUser()
      return Promise.reject(error)
    }

    // handleApiError(error)
    return Promise.reject(error)
  },
)

export default api
