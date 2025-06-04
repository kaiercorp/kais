import { APICore } from './apiCore'
import { QueryClient } from '@tanstack/react-query'
import { QUERY_KEY } from 'helpers'

const api = new APICore()

export const ApiDownloadFile = (queryClient: QueryClient, url: string, filename: string, params?: { [key: string]: any }) => {
  queryClient.fetchQuery({
    queryKey: [QUERY_KEY.downloadFile, url],
    queryFn: async () => {
      const response: any = await api.getFile(url, params)

      const fileObjectUrl = window.URL.createObjectURL(response)
      const link = document.createElement('a')
      link.href = fileObjectUrl
      link.style.display = 'none'
      link.download = filename

      document.body.appendChild(link)
      link.click()
      setTimeout(() => {
        window.URL.revokeObjectURL(fileObjectUrl)
      }, 1000 * 60 * 10)
      link.remove()

      return true
    }
  })
}