import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { APICore } from './apiCore'
import { QUERY_KEY } from 'helpers/api'

const api = new APICore()

export function ApiLogin() {
  return useMutation({
    mutationKey: [QUERY_KEY.login],
    mutationFn: async (params: { username: string; password: string }) => {
      const response = await api.post(`/auth/login`, params)
      return response.data
    },
    onSuccess: (data) => {
      api.setLoggedInUser(data)
    },
    onError: (error) => {
      api.setLoggedInUser(null)
    }
  })
}

export function ApiLogout() {
  const navigate = useNavigate()

  return useMutation({
    mutationKey: [QUERY_KEY.logout],
    mutationFn: async () => {
      const user = api.getLoggedInUser()
      const response = await api.post(`/auth/logout`, { username: user.username, token: user.token })
      return response.data
    },
    onSuccess: () => {
      api.setLoggedInUser(null)
      navigate('/auth/login')
    }
  })
}