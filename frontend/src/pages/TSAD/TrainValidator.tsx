
export type emptyErrors = {
  data_path?: string
  gpus?: string
  epochs?: string
  default_config_file?: string
  base_lr?: string
  class_list?: string
  save_top_k?: string
  train_batch_size?: string
  trial_name?: string
  parent_id?: string
  test_db?: string
  threshold_k?: string
  date_column?: any
  hasError: boolean
}

export const validateTSAD = (trainType: string | undefined, values: any, t: any) => {
  let errors: emptyErrors = {
    data_path: '',
    gpus: '',
    epochs: '',
    default_config_file: '',
    base_lr: '',
    class_list: '',
    save_top_k: '',
    train_batch_size: '',
    trial_name: '',
    parent_id: '',
    test_db: '',
    hasError: false,
    threshold_k: '',
    date_column: '',
  }

  if (!trainType || !values) {
    return { error: t('validator.invalid'), hasError: true }
  }

  // if (!values.gpus || values.gpus.length < 1) {
  //   errors.gpus = t('validator.gpu')
  //   errors.hasError = true
  // }

  if (trainType === 'additional') {
    return errors
  }

  if (!values.data_path || values.data_path === '') {
    errors.data_path = t('validator.dataset.train')
    errors.hasError = true
  }
  if (values.trial_name === '') {
    errors.trial_name = t('validator.name')
    errors.hasError = true
  }
  if (values.test_db === '') {
    errors.test_db = t('validator.name.dataset')
    errors.hasError = true
  }
  if (values.train_config.date_column === '') {
    errors.date_column = t('validator.datecolumn')
    errors.hasError = true
  }

  return errors
}