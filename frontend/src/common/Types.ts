import { Engine } from './Const'

export type ConfigType = {
    id?: number
    config_key?: string
    config_val?: string
}

export type DiskType = {
    free: number
    fstype: string
    inodesFree: number
    inodesTotal: number
    inodesUsed: number
    inodesUsedPercent: number
    path: string
    total: number
    used: number
    usedPercent: number
}

export type DirectoryType = {
    id: number
    parent_id: number
    name: string
    data_type: string
    path: string
    description: string
    is_open: boolean 
    is_deleted: boolean
    is_leaf: boolean
    is_testable: boolean
    is_trainable: boolean
    is_valid: boolean
    created_at: string
    updated_at: string
    deleted_at: string
    dirs?: DirectoryType[]
}

export type QuickSortType = 'created_at' | 'accuracy' | 'inference_time'

export type ProjectType = {
  project_id?: number
  project_name?: string
  state?: string
  description?: string
  created_at?: string
  visible?: boolean
  num_trials?: number
  num_trains?: number
  num_tests?: number
  category?: string
}

export type UserType = {
    id: number
    email: string
    username: string
    password: string
    firstName: string
    lastName: string
    role: string
    token: string
    department: string
}

export type GPUStatusType = {
  total: number
  working: number
  idle: number
  train: number
  test: number
  gpus: {
    name: string
    state: string
    id: number
  }[]
}

export type GPUType = {
  id: number
  is_running: boolean
  is_use: boolean
  name: string
  state: string
  created_at: string
  updated_at: string 
  use_gpu: string
  use_mem: string
}

export type TestFolderType = {
  [key: string]: any
  gpus: [],
  trial: any,
  data_path: string,
  model: string,
  heatmap_download: boolean,
}

export type TestMultiSampleType = {
  [key: string]: any
  gpus: [],
  trial: any,
  data_path: string,
  model: string,
}

export type FilterType = {
  state: {
    total: boolean,
    train: boolean,
    additional_train: boolean,
    cancel: boolean,
    finish: boolean,
    'finish-fail': boolean,
    fail: boolean,
    test: boolean,
    finish_test: boolean,
    test_cancel: boolean,
    test_fail: boolean,
    idle: boolean 
  } 
  accuracy: {
    min: number,
    max: number  
  }
  precision: {
    min: number,
    max: number   
  }
  recall: {
    min: number,
    max: number
  },
  f1: {
    min: number,
    max: number
  },
  inference_time: {
    min: number,
    max: number
  },
  startDate: Date
  endDate: Date
}

export type TrialType = {
  [key: string]: any
  data_path: string
  project_id?: number
  gpus?: {}[]
  trial_name?: string
  parent_id?: number
  train_once?: boolean
  test_db?: string
  train_config: {
    [key: string]: any
    width: number
    height: number
    class_list: string[]
  }
}

// 유지보수하며 정의할 예정
// export type TrainContainerTrialType = {
//   trial_id: number 
//   trial_local_id: number
//   trial_name: string
//   trial_type: string 
  
//   parent_id: number
//   parent_trial: string
  
//   project_id: number
  
//   train_type: string
//   train_count: number
//   train: {

//   }
  
//   test: {
    
//   }
//   test_db: string
//   test_result: ObjToJsonType
  
//   perf: ObjToJsonType
//   progress: number
//   state: string 
//   target_metric: ObjToJsonType
//   inference_time: number
//   is_deleted: boolean
//   is_use: boolean
  
//   params: {
//     
//   }
//   params_parent: ObjToJsonType
  
//   accuracy: number
//   best_train_uuid: ObjToJsonType
//   date_type: string
//   dataset_id: number

//   gpu: any | null 
//   gpus: any | null

//   created_at: string
//   update_at: string
//   delete_at: string | null
// }

export type TrialsType = TrialType[]

export type CreateTrainType = {
  project_id: number  
  trial_name: string
  test_db?: string
  dataset_id: number 
  train_type: string
  data_type?: string
  engine_type: string
  trial_id: number
  params: string
}
  
export type TrainConfigType = {
  base_lr: number
  class_list: []
  default_config_file: string 
  epochs: number 
  height: number 
  save_top_k: number 
  target_metric: {
    wa: number
    uwa: number 
    precision: number 
    recall: number 
    f1: number 
  },
  train_batch_size: number 
  width: number
  auto_stop: boolean
}

export type TrialConfigType = {
  test_db: string,
  train_once: true,
  trial_name: string,
}


export type BaseModalTitleType = {
    title: string
    icon: string
    description: string
    size: "xl" | "sm" | "lg" | undefined
    submitText: string
}

export type EngineType = typeof Engine[keyof typeof Engine]

export type MenuType = {
  key: string
  label: string
  icon: string
  url: string
  isUse: boolean
  isTitle: boolean
  parentKey: string
  children: MenuType[] 
}

export type HPOModelType = {
  [key: string]: any
  id?: number
  engine_type: string
  model_name: string
  class_type: string
  is_use: boolean
  params: HPOParamsType[]
}

export type HPOParamsType = {
  [key: string]: any
  id?: number
  name: string
  suggest_type: string
  data_type: string
  model_id?: number
  is_use: boolean
  dists: HPOParamsDistType[]
}

export type HPOParamsDistType = {
  [key: string]: any
  id?: number
  param_id?: number
  is_use: boolean
  dist: string
  cond: DistCondType
}

export type DistCondType = {
  [key: string]: any
  key: string
  operator: string
  value: string
}

export type DatasetRootType = {
  [key: string]: any
  id?: number
  name?: string
  path?: string
  is_use?: boolean
  datasets?: any
}

export type ChartType = {
  uuid: string
  train_id: number
  target_metric: string
  name: string
  all_result: ObjToJsonType 
  cf_matrix: ObjToJsonType 
  class_result: ObjToJsonType 
  feature_importance: ObjToJsonType 
  learning_chart_material: ObjToJsonType 
  score: string
  probs: string
  error_graph: string
  updated_at: string
} 

export type TrainType = {
  id: number
  local_id: number
  uuid: string
  trial_uuid: string
  description: string   
  inference_time: number
  acc: number
  perf: string
  progress: number
  state: string
  target_metric: string
  config_path: string
  save_path: string
  model_list: ObjToJsonType
  result: ObjToJsonType
  test_result: ObjToJsonType 
  created_at: string
  updated_at: string
}

type ObjToJsonType = {
  String: string
  value: false
}