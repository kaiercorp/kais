import { APICore } from "./apiCore"
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { QUERY_KEY } from 'helpers/api'
import { ConfigType } from 'common'

const api = new APICore()

export function ApiFetchConfigs() {
    const { isLoading, isError, error, data } = useQuery({
       queryKey: [QUERY_KEY.configurations],
       queryFn: async (): Promise<{configs: ConfigType[], version: string}> => {
            const response = await api.post(`/configuration/list`, {})

            let version = response.data.filter((config: ConfigType) => config.config_key === 'VERSION')[0]?.config_val
    
            if (!version) {
                version = 'KAI.S'
            }           
            
            return { configs: response.data , version }
       } ,
       refetchOnMount: true,
    })

    if ( isLoading ) return { isLoading, configs:[] }
    
    if ( isError ) return { error, configs: [] }
   
    return data? {configs: data.configs, version: data.version} : {configs: [], version: 'KAI.S'} 
}

export function ApiUpdateConfigs() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationKey: [`update_${QUERY_KEY.configurations}`],
        mutationFn: async (data: ConfigType[]) => {
            return await api.post(`/configuration`, data)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey : [`${QUERY_KEY.configurations}`] })
            alert('저장되었습니다.')
        },
        onError: (error) => {
          alert('변경사항 저장에 실패했습니다.')
        }
    })
}
