import { APICore } from './apiCore'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { QUERY_KEY } from 'helpers'
import { DirectoryType, CreateTrainType } from 'common'
import { logger } from 'helpers/logger/logger'

const api = new APICore()

export function ApiFetchTrial(trialId: number | undefined) {
  const { isLoading, isError, error, data } = useQuery({
    queryKey: [QUERY_KEY.fetchTrial, trialId],
    queryFn: async () => {
      if (typeof trialId === 'undefined') return {}

      const response = await api.get(`/trial/${trialId}`, {})

      let trial = response.data

      try {
        if (trial.params) trial['params'] = JSON.parse(trial.params)
        if (trial.train && trial.train.result && trial.train.result.String) trial.train['result'] = JSON.parse(trial.train.result.String)
        if (trial.train) trial.train.model_list = trial.train.model_list.split(",")
        if (trial.test && trial.test.test_result) {
          let results = JSON.parse(trial.test.test_result)
          let keys = Object.keys(results)
          let rr = Array<any>()
          keys.forEach(function (k: string) {
            rr.push({ 'model_name': k, 'model_result': results[k] })
          })
          trial.test.test_result = rr
        }
        if (trial.test && trial.test.test_summary) trial.test.test_summary = JSON.parse(trial.test.test_summary)
        if (trial.parent_trial) {
          trial.parent_trial = JSON.parse(trial.parent_trial)
          trial.parent_trial.params = JSON.parse(trial.parent_trial.params)
        }
      } catch (e) {
        logger.error(e)
      }

      return trial
    }
  })

  if (isLoading) return { isLoading }

  if (isError) return { error }

  return { trial: data ? data : {} }
}

export function ApiDeleteTrials() {
  return useMutation({
    mutationKey: [QUERY_KEY.deleteTrials],
    mutationFn: async (trialList: number[]) => {
      return await api.post(`/trial`, { trial_list: trialList })
    },
    onSuccess: (data) => {
      alert('삭제되었습니다.')
    },
    onError: (error) => {
      alert('삭제에 실패했습니다.')
    }
  })
}

export function ApiFetchTrialCompare() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.fetchTrialCompare],
    mutationFn: async (trialList: number[]) => {
      const response = await api.post(`/trial/compare`, { trial_list: trialList })

      return response.data
    },
    onSuccess: (data) => {
      queryClient.setQueryData([QUERY_KEY.fetchTrialCompare], data)
    }
  })
}

export function ApiFetchTrainModels(trialId: number | undefined) {
  const { isLoading, isError, error, data } = useQuery({
    queryKey: [QUERY_KEY.fetchTrainModels, trialId],
    queryFn: async () => {
      if (typeof trialId === 'undefined') return []

      const response = await api.get(`/trial/model/list/${trialId}`, {})

      const models = response.data

      try {
        models.forEach((m: any) => {
          if (m.all_result.Valid) {
            m.all_result = JSON.parse(m.all_result.String)
          }
  
          if (m.class_result.Valid) {
            m.class_result = JSON.parse(m.class_result.String)
          }
          
          if (m.accuracy_info.Valid) {
            m.accuracy_info = JSON.parse(m.accuracy_info.String)
          }
        })
      } catch (e) {
        logger.error(e)
      }

      return models
    }
  })

  if (isLoading) return { isLoading }

  if (isError) return { error }

  return { models: data ? data : [] }
}

export function ApiCreateTrain() {
  return useMutation({
    mutationKey: [QUERY_KEY.createTrain],
    mutationFn: async (data: Partial<CreateTrainType> | undefined) => {
      if (typeof data === 'undefined') return

      return await api.post(`/train`, data)
    },
    onSuccess: (data) => {
      alert('모델링이 시작되었습니다.')
    },
    onError: (error) => {
      alert('모델링 시작에 실패했습니다.')
    }
  })
}

export function ApiStopTrain() {
  return useMutation({
    mutationKey: [QUERY_KEY.stopTrain],
    mutationFn: async (trialId: number) => {
      return await api.delete(`/train/${trialId}`)
    }
  })
}

export function ApiFetchChartData(trialId: number | undefined, epoch: number | undefined) {
  const { isLoading, isError, error, data } = useQuery({
    queryKey: [QUERY_KEY.fetchChartData, trialId, epoch],
    queryFn: async () => {
      if (typeof trialId === 'undefined' || typeof epoch === 'undefined') return []

      const response = await api.get(`/train/chart/${trialId}/${epoch}`, {})

      return response.data
    }
  })

  if (isLoading) return { isLoading }

  if (isError) return { error }

  return { chartData: data ? data : [] }
}

export function ApiCreateTestDirectory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.createTestDirectory],
    mutationFn: async (data: any) => {
      const response = await api.post(`/test/dir`, data)

      return response.data
    },
    onSuccess: (data) => {
      queryClient.setQueryData([QUERY_KEY.createTestDirectory], data)
      alert('Multi Sample Test 시작되었습니다.')
    },
    onError: (error) => {
      alert('Multi Sample Test에 실패했습니다.')
    }
  })
}

export function ApiCreateTestFile() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.createTestFile],
    mutationFn: async (data: any) => {
      const response = await api.createWithFile(`/test/file`, data)

      return response.data
    },
    onSuccess: (data) => {
      queryClient.setQueryData([QUERY_KEY.createTestFile], data)
    }
  })
}

export function ApiStopTest() {
  return useMutation({
    mutationKey: [QUERY_KEY.stopTest],
    mutationFn: async (trialId: number) => {
      return await api.delete(`/test/${trialId}`)
    }
  })
}

export function ApiGetRowFromFile() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.getRowFromFile],
    mutationFn: async (data: any) => {
      const response = await api.createWithFile(`/test/row`, data)

      return response.data
    },
    onSuccess: (data) => {
      queryClient.setQueryData([QUERY_KEY.getRowFromFile], data)
    }
  })
}

export function ApiFetchDirs() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationKey: [QUERY_KEY.fetchDirectories],
    mutationFn: async (data: { parentId: number, dataType: string }) => {
      const response = await api.post(`/dataset`, { parent_id: data.parentId, data_type: data.dataType })
      let dirs = response.data
      if (data.parentId !== 0) {
        const directories = queryClient.getQueryData<DirectoryType[]>([QUERY_KEY.fetchDirectories])
        const target = directories?.find((elem: any) => elem.id === data.parentId)
        if (target) {
          target['dirs'] = dirs
          target['is_open'] = true
        }
        dirs = directories
      } else {
        dirs = dirs.map((dir: any) => { return { ...dir, is_open: false } })
      }

      return dirs
    },
    onSuccess: (data) => {
      queryClient.setQueryData([QUERY_KEY.fetchDirectories], data)
    }
  })
}

export function ApiFetchDatasets(dataType: string, parentId: number | undefined, engineType?: string) {
  const { isLoading, isError, error, data } = useQuery({
    queryKey: [QUERY_KEY.fetchDatasets, dataType, parentId, engineType],
    queryFn: async () => {
      if (typeof parentId === 'undefined') return dataType === 'table' ? [] : null

      let url;
      switch (dataType) {
        case 'table':
          url = `/dataset/column/${parentId}`
          break;
        case 'image':
          url = `/dataset/classes/${parentId}/${engineType}`
          break;
      }

      if (!url) return dataType === 'table' ? [] : null

      const response = await api.get(url, {})

      return response.data
    }
  })

  if (isLoading) return { isLoading }

  if (isError) return { error }

  let datasets;

  switch (dataType) {
    case 'table':
      datasets = { columns: data ? data : [] }
      break;
    case 'image':
      datasets = { classes: data ? data : null }
      break;
    default:
      datasets = { classes: null }
      break;
  }

  return datasets
}