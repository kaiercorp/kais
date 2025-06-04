import { TrialType } from 'common'

export type emptyErrors = {
  data_path?: string
  gpus?: string
  width?: string
  height?: string
  class_list?: string
  trial_name?: string
  parent_id?: string
  test_db?: string
  hasError: boolean
}

export const validateCLS = (trainType: string | undefined, values: TrialType | undefined, t: any) => {
  let errors: emptyErrors = {
    data_path: '',
    gpus: '',
    width: '',
    height: '',
    class_list: '',
    trial_name: '',
    parent_id: '',
    test_db: '',
    hasError: false,
  }

  if (!trainType || !values) {
    return { error: t('validator.invalid'), hasError: true }
  }

  if (!values.gpus || values.gpus.length < 1) {
    errors.gpus = t('validator.gpu')
    errors.hasError = true
  }

  if (trainType === 'additional') {
    return errors
  }

  if (!values.data_path) {
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
  if (values.train_config.width === 0) {
    errors.width = t('validator.width')
    errors.hasError = true
  }
  if (values.train_config.height === 0) {
    errors.height = t('validator.height')
    errors.hasError = true
  }

  return errors
}