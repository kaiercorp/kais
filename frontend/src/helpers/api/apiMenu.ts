import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { APICore } from './apiCore'
import { QUERY_KEY } from 'helpers/api'
import { MenuType } from 'common'

const api = new APICore()

export const ApiFetchMenu = () => {
    const { isLoading, isError, error, data} = useQuery({
        queryKey: [QUERY_KEY.fetchMenu],
        queryFn: async () => {
            const response = await api.get('/menu', {}) 
            return response.data 
        },
        refetchOnMount: true,
    })
    
    if (isLoading) return {isLoading}

    if (isError) return {error} 

    return {menu: data ? data : []}
}

export const ApiUpdateMenu = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationKey: [QUERY_KEY.updateMenu],
        mutationFn: async (menus: MenuType[]) => {
            return await api.post('/menu', menus)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: [QUERY_KEY.fetchMenu], exact: true})
            alert('저장되었습니다.')           
        },
        onError: (error) => {
          alert('저장에 실패했습니다.')
        }
    }) 
}