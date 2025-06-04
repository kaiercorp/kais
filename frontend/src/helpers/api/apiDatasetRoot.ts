import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query' 
import { DatasetRootType } from 'common'
import { APICore } from './apiCore'
import { QUERY_KEY } from './queryKey'

const api = new APICore()

export function ApiFetchDRs() {
  
  const { isLoading, isError, error, data } = useQuery({
    queryKey: [QUERY_KEY.dataset_root],
    queryFn: async () => {
      const response = await api.get(`/dataroot`, {})
      return response.data
    }    
  })
  
  if (isLoading) return { isLoading }

  if (isError) return { error }
  
  return { datasetroots: data? data : [] }
}

export function ApiCreateDR() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [`create_${QUERY_KEY.dataset_root}`],
    mutationFn: async (params: DatasetRootType) => {
      return await api.post(`/dataroot`, params)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [QUERY_KEY.dataset_root], exact: true})
      alert('Root Path 추가되었습니다.')
    },
    onError: (error) => {
      alert('Root Path 추가에 실패했습니다.')
    }
  })
}

export function ApiDeleteDR() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [`delete_${QUERY_KEY.dataset_root}`],
    mutationFn: async (id: number) => {
      return await api.delete(`/dataroot/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [QUERY_KEY.dataset_root], exact: true})
      alert('Root Path 삭제되었습니다.')
    },
    onError: (error) => {
      alert('Root Path 삭제에 실패했습니다.')
    }
  })
}

export function ApiUpdateDR() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [`update_${QUERY_KEY.dataset_root}`],
    mutationFn: async (params: DatasetRootType) => {
      return await api.update(`/dataroot`, params)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [QUERY_KEY.dataset_root], exact: true})
      alert('변경사항이 저장되었습니다.')
    },
    onError: (error) => {
      alert('변경사항 저장에 실패했습니다.')
    }
  })
}

export function ApiDeleteDS() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [`delete_${QUERY_KEY.dataset_root}`],
    mutationFn: async (id: number) => {
      return await api.delete(`/dataset/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [`${QUERY_KEY.dataset_root}`], exact: true})
      alert('Dataset 삭제되었습니다.')
    },
    onError: (error) => {
      alert('Dataset 삭제에 실패했습니다.')
    }
  })
}