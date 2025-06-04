
export type emptyErrorsTrainTableCls = {
    gpus?: string
    trial_name?: string
    test_db?: string
    target_metric?: string
    label_column?: string
    hasError: boolean
}

export const validate = (trainType: any, values: any, t: any) => {
    let errors: emptyErrorsTrainTableCls = {
        trial_name: '',
        test_db: '',
        target_metric: '',
        label_column: '',
        hasError: false,
    }

    if (!values) {
        return { error: t('validator.invalid'), hasError: true }
    }

    // if (!values.gpus || values.gpus.length < 1) {
    //     errors.gpus = t('validator.gpu')
    //     errors.hasError = true
    // }

    if (trainType === 'additional') {
      return errors
    }

    if (!values.trial_name || values.trial_name === '') {
        errors.trial_name = t('validator.name')
        errors.hasError = true
    }
    if (!values.test_db || values.test_db === '') {
        errors.test_db = t('validator.name.dataset')
        errors.hasError = true
    }
    if (!values.train_config.target_metric || values.train_config.target_metric === '') {
        errors.target_metric = t('validator.targetmetric')
        errors.hasError = true
    }
    if (!values.train_config.label_column || values.train_config.label_column === '') {
        errors.label_column = t('validator.targetCol')
        errors.hasError = true
    }

    return errors
}