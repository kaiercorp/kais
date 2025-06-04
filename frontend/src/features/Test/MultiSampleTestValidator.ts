import { TestMultiSampleType } from 'common'

export type emptyTestMultiSampleErrors = {
    dataset_id?: string
    gpus?: string
    heatmap_download?: string
    model_list?: string
    parent_id?: string
    hasError: boolean
}

export const validate = (values: TestMultiSampleType, t: any) => {
    let errors: emptyTestMultiSampleErrors = {
        dataset_id: '',
        gpus: '',
        heatmap_download: '',
        model_list: '',
        parent_id: '',
        hasError: false
    }

    if (!values) {
        return { error: t('validator.invalid'), hasError: true}
    }

    if (!values.dataset_id || values.dataset_id < 1) {
        errors.dataset_id = t('validator.dataset.test')
        errors.hasError = true
    }
    // if (!values.gpu || values.gpu.length < 1) {
    //     errors.gpus = t('validator.gpu')
    //     errors.hasError = true
    // }
    if (!values.model_list || values.model_list.length < 1) {
        errors.model_list = t('validator.model')
        errors.hasError = true
    }
    if (!values.parent_id || values.parent_id < 1) {
        errors.parent_id = t('validator.train')
        errors.hasError = true
    }

    return errors
}