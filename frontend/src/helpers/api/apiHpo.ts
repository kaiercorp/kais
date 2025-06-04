import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query' 
import { HPOModelType } from 'common'
import { APICore } from './apiCore'
import { QUERY_KEY } from './queryKey'

const api = new APICore()

export function ApiFetchHpos() {
    const { isLoading, isError, error, data } = useQuery({
        queryKey: [`${QUERY_KEY.hpo}`],
        queryFn: async () => {
            const response = await api.get(`/hpo/list`, {})

            return response.data
        }
    })

    if (isLoading) return { isLoading }

    if (isError) return { error }

    return { hpo: data? data : []}
}

export function ApiUpdateHpos() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationKey: [`update_${QUERY_KEY.hpo}`],
        mutationFn: async (data: HPOModelType[]) => {
            return await api.post(`/hpo`, data)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey : [`${QUERY_KEY.hpo}`]})
            alert('저장되었습니다.')
        },
        onError: (error) => {
          alert('변경사항 저장에 실패했습니다.')
        }
    })
}

export function ApiInitHpos() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationKey: [`update_${QUERY_KEY.hpo}`],
        mutationFn: async () => {
            return await api.post(`/hpo/init`, null)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey : [`${QUERY_KEY.hpo}`]})
            alert('Hyper parameter 목록이 초기화되었습니다.')
        },
        onError: (error) => {
          alert('Hyper parameter 초기화 실패')
        }
    })
}