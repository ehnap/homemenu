import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => {
    const data = response.data
    if (data.code !== 0) {
      return Promise.reject(new Error(data.message || 'Request failed'))
    }
    return data.data
  },
  async (error) => {
    const url = error.config?.url || ''
    const isAuthRequest = url.startsWith('/auth/')

    if (error.response?.status === 401 && !isAuthRequest) {
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken && !error.config._retry) {
        error.config._retry = true
        try {
          const res = await axios.post('/api/auth/refresh', {
            refresh_token: refreshToken,
          })
          const tokens = res.data.data
          localStorage.setItem('access_token', tokens.access_token)
          localStorage.setItem('refresh_token', tokens.refresh_token)
          error.config.headers.Authorization = `Bearer ${tokens.access_token}`
          return client(error.config)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/login'
        }
      } else {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        window.location.href = '/login'
      }
    }
    const msg = error.response?.data?.message || error.message || 'Request failed'
    return Promise.reject(new Error(msg))
  }
)

export default client
