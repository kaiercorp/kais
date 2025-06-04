
export const emptyTestFolder = {
  gpus: [0],
  trial: {},
  data_path: '',
  model: '',
  heatmap_download: false,
  train_type: 'test',
  model_list: [],
}

export const emptyTestFile = {
  gpus: [0],
  trial: { trial_id: 0},
  file: '',
  model: '',
  train_type: 'test',
}

export const emptyMultiSampleTest = {
  gpus: [0],
  trial: {},
  data_path: '',
  model: '',
  train_type: 'test',
  model_list: [],
}

export const initialFilter = {
  state: {
    total: true,
    train: true,
    additional_train: true,
    cancel: true,
    finish: true,
    'finish-fail': true,
    fail: true,
    test: true,
    finish_test: true,
    test_cancel: true,
    test_fail: true,
    idle: true
  },
  accuracy: {
    min: 0,
    max: 100,
  },
  precision: {
    min: 0,
    max: 100,
  },
  recall: {
    min: 0,
    max: 100,
  },
  f1: {
    min: 0,
    max: 100,
  },
  inference_time: {
    min: 0,
    max: 1000,
  },
  startDate: new Date(new Date().getTime() - 24 * 60 * 60 * 1000),
  endDate: new Date(),
}

export const emptySLClsTrial = {
  project_id: 0,
  trial_name: '',
  data_path: '',
  gpus: [0],
  trial_id:0,
  parent_id: -1,
  test_db: '',
  train_once: false,
  dataset_id: 0,
  train_config: {
    class_list: [],
    height: 512,
    target_metric: {
      wa: 100,
      uwa: 0,
      precision: 0,
      recall: 0,
      f1: 0
    },
    width: 512,
    auto_stop: true,
  },
}

export const emptyMLClsTrial = {
  project_id: 0,
  trial_name: '',
  data_path: '',
  gpus: [0],
  trial_id:0,
  parent_id: -1,
  test_db: '',
  train_once: false,
  dataset_id: 0,
  train_config: {
    class_list: [],
    height: 512,
    target_metric: {
      image_accuracy: 0,
      image_precision: 0,
      image_recall: 0,
      image_f1_score: 100,
      label_accuracy: 0,
      label_precision: 0,
      label_recall: 0,
      label_f1_score: 0,
    },
    width: 512,
    auto_stop: true,
  },
}

export const emptyVADTrial = {
  project_id: 0,
  trial_name: '',
  data_path: '',
  gpus: [0],
  trial_id:0,
  parent_id: -1,
  test_db: '',
  train_once: false,
  dataset_id: 0,
  train_config: {
    class_list: [],
    height: 512,
    target_metric: {
      wa: 0,
      uwa: 0,
      precision: 0,
      recall: 100,
      f1: 0,
      auroc: 0,
      prauc: 0
    },
    width: 512,
    auto_stop: true,
    tiff_frame_number: 0
  },
}

export const emptyTableClsRequest = {
  project_id: 0,
  trial_name: '',
  data_path: '',
  gpus: [],
  trial_id:0,
  parent_id: -1,
  test_db: '',
  train_once: false,
  dataset_id: 0,
  train_config: {
    target_metric: {
      wa: 100,
      uwa: 0,
      precision: 0,
      recall: 0,
      f1: 0
    },
    label_column: '',
    index_column: '',
    date_colulmn: '',
    input_column: [],
    except_column: [],
    auto_stop: true,
  },
}

export const emptyTableRegRequest = {
  project_id: 0,
  trial_name: '',
  data_path: '',
  gpus: [],
  trial_id:0,
  parent_id: -1,
  test_db: '',
  train_once: false,
  dataset_id: 0,
  train_config: {
    target_metric: {
      mse: 100,
      rmse: 0,
      mae: 0,
    },
    label_column: '',
    index_column: '',
    input_column: [],
    except_column: [],
    auto_stop: true,
  },
}

export const emptyTSADRequest = {
  project_id: 0,
  trial_name: '',
  data_path: '',
  gpus: [],
  trial_id:0,
  parent_id: -1,
  test_db: '',
  train_once: false,
  dataset_id: 0,
  train_config: {
    target_metric: {
      mse: 100,
      rmse: 0,
      mae: 0,
    },
    label_column: 'none',
    index_column: '',
    date_column: '',
    input_column: [],
    except_column: [],
    auto_stop: true,
    date_format: '%Y-%m-%d %H:%M:%S',
    threshold_mode: 'std',
    threshold_k: 3.0
  },
}
