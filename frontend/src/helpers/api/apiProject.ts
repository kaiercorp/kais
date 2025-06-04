import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query' 
import { ProjectType } from 'common'
import { APICore } from './apiCore'
import { QUERY_KEY } from './queryKey'

const api = new APICore()

export function ApiFetchProjects(category: string) {
  const { isLoading, isError, error, data } = useQuery({
    queryKey: [QUERY_KEY.projects, category],
    queryFn: async () => {
      if (category.includes('config.')) return
      if (category === '') return
      const response = await api.get(`/project/list/${category}`, {})
      
      return response.data
    },
    refetchOnMount: true,
  })
  
  if (isLoading) return { isLoading }

  if (isError) return { error }
  
  return { projects: data? data : [] }
}

export function ApiCreateProject(category: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.projects, category],
    mutationFn: async (params: ProjectType) => {
      return await api.post(`/project`, params)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [QUERY_KEY.projects] })
      alert('Project 추가되었습니다.')
    },
    onError: (error) => {
      alert('Project 추가에 실패했습니다.')
    }
  })
}

export function ApiDeleteProject(category: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.projects, category],
    mutationFn: async (projectId: number) => {
      return await api.delete(`/project/${projectId}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [QUERY_KEY.projects ] })
      alert('Project 삭제되었습니다.')
    },
    onError: (error) => {
      alert('Project 삭제에 실패했습니다.')
    }
  })
}

export function ApiUpdateProject(category: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.projects, category],
    mutationFn: async (params: ProjectType) => {
      return await api.update(`/project`, params)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey : [QUERY_KEY.projects] })
      alert('저장되었습니다.')
    },
    onError: (error) => {
      alert('변경사항 저장에 실패했습니다.')
    }
  })
}