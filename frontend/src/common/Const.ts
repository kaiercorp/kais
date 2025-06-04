import { BaseModalTitleType } from './Types'

export const Engine= {
    vcls: 'vision.cls',
    vad: 'vision.ad',
    tcls: 'table.cls',
    treg: 'table.reg',
    tsad: 'ts.ad',
} as const

export const tableIndexTypes = [
  'id',
  'index',
  'idx',
  'pid',
  'kaier_id'
]

type TResultOrder = {
    [key: string]: number
}
export const ResultOrder: TResultOrder = {
    wa: 1,
    auroc: 2,
    prauc: 3,
    f1: 4,
    precision: 5,
    recall: 6,
    uwa: 7,
    mae: 8,
    mse: 9,
    rmse: 10,
    image_accuracy: 11, 
    image_f1_score: 12,
    image_precision: 13,
    image_recall: 14,
    label_accuracy: 15,
    label_f1_score: 16,
    label_precision: 17,
    label_recall: 18,
}

export const emptyBaseModalTitle: BaseModalTitleType = {
    title: '',
    icon: '',
    description: '',
    size: 'lg',
    submitText: 'Submit'
}

/* Vision Cls */
export const visionClsAutoTrain: BaseModalTitleType = {
    title:'modal.title.vision.cls.autotrain.title', 
    icon: 'modal.title.vision.cls.autotrain.icon', 
    description:'modal.title.vision.cls.autotrain.description',
    size: undefined,
    submitText: 'modal.button.ok.train'
}

export const visionClsMLAutoTrain: BaseModalTitleType = {
    title:'modal.title.vision.cls.autotrain.title', 
    icon: 'modal.title.vision.cls.autotrain.icon', 
    description:'modal.title.vision.cls.autotrain.description',
    size: 'lg',
    submitText: 'modal.button.ok.train'
}

export const visionClsExpertTrain: BaseModalTitleType = {
    title:'modal.title.vision.cls.experttrain.title', 
    icon: 'modal.title.vision.cls.experttrain.icon', 
    description:'modal.title.vision.cls.experttrain.description',
    size: 'xl',
    submitText: 'modal.button.ok.train'
}

export const visionClsMultiTest: BaseModalTitleType = {
    title:'modal.title.vision.cls.multitest.title', 
    icon: 'modal.title.vision.cls.multitest.icon', 
    description:'modal.title.vision.cls.multitest.description',
    size: undefined,
    submitText: 'modal.button.ok.test'
}

export const visionClsSingleTest: BaseModalTitleType = {
    title:'modal.title.vision.cls.singletest.title', 
    icon: 'modal.title.vision.cls.singletest.icon', 
    description:'modal.title.vision.cls.singletest.description',
    size: 'xl',
    submitText: ''
}

export const visionClsAdditionalTrain: BaseModalTitleType = {
    title:'modal.title.vision.cls.additional.title', 
    icon: 'modal.title.vision.cls.additional.icon', 
    description:'modal.title.vision.cls.additional.description',
    size: undefined,
    submitText: 'modal.button.ok.train'
}

/* Table Cls */
export const tableClsAutoTrain: BaseModalTitleType = {
    title:'modal.title.table.cls.autotrain.title', 
    icon: 'modal.title.table.cls.autotrain.icon', 
    description:'modal.title.table.cls.experttrain.description',
    size: 'lg',
    submitText: 'modal.button.ok.train'
}

export const tableClsMultiTest: BaseModalTitleType = {
    title:'modal.title.table.cls.multitest.title', 
    icon: 'modal.title.table.cls.multitest.icon', 
    description:'modal.title.table.cls.multitest.description',
    size: 'xl',
    submitText: 'modal.button.ok.test'
}

export const tableClsSingleTest: BaseModalTitleType = {
    title:'modal.title.table.cls.singletest.title', 
    icon: 'modal.title.table.cls.singletest.icon', 
    description:'modal.title.table.cls.singletest.description',
    size: 'xl',
    submitText: ''
}

export const tableClsAdditionalTrain: BaseModalTitleType = {
    title:'modal.title.table.cls.additional.title', 
    icon: 'modal.title.table.cls.additional.icon', 
    description:'modal.title.table.cls.additional.description',
    size: undefined,
    submitText: 'modal.button.ok.train'
}

/* Table Regression */
export const tableRegAutoTrain: BaseModalTitleType = {
    title:'modal.title.table.reg.autotrain.title', 
    icon: 'modal.title.table.reg.autotrain.icon', 
    description:'modal.title.table.reg.experttrain.description',
    size: 'lg',
    submitText: 'modal.button.ok.train'
}

export const tableRegMultiTest: BaseModalTitleType = {
    title:'modal.title.table.reg.multitest.title', 
    icon: 'modal.title.table.reg.multitest.icon', 
    description:'modal.title.table.reg.multitest.description',
    size: 'xl',
    submitText: 'modal.button.ok.test'
}

export const tableRegSingleTest: BaseModalTitleType = {
    title:'modal.title.table.reg.singletest.title', 
    icon: 'modal.title.table.reg.singletest.icon', 
    description:'modal.title.table.reg.singletest.description',
    size: 'xl',
    submitText: ''
}

export const tableRegAdditionalTrain: BaseModalTitleType = {
    title:'modal.title.table.reg.additional.title', 
    icon: 'modal.title.table.reg.additional.icon', 
    description:'modal.title.table.reg.additional.description',
    size: undefined,
    submitText: 'modal.button.ok.train'
}

/* Vision AD */
export const visionADAutoTrain: BaseModalTitleType = {
    title:'modal.title.vision.ad.autotrain.title', 
    icon: 'modal.title.vision.ad.autotrain.icon', 
    description:'modal.title.vision.ad.autotrain.description',
    size: undefined,
    submitText: 'modal.button.ok.train'
}

export const visionADMultiTest: BaseModalTitleType = {
    title:'modal.title.vision.ad.multitest.title', 
    icon: 'modal.title.vision.ad.multitest.icon', 
    description:'modal.title.vision.ad.multitest.description',
    size: undefined,
    submitText: 'modal.button.ok.test'
}

export const visionADSingleTest: BaseModalTitleType = {
    title:'modal.title.vision.ad.singletest.title', 
    icon: 'modal.title.vision.ad.singletest.icon', 
    description:'modal.title.vision.ad.singletest.description',
    size: 'xl',
    submitText: ''
}

export const visionADAdditionalTrain: BaseModalTitleType = {
    title:'modal.title.vision.ad.additional.title', 
    icon: 'modal.title.vision.ad.additional.icon', 
    description:'modal.title.vision.ad.additional.description',
    size: undefined,
    submitText: 'modal.button.ok.train'
}

/* TS AD */
export const tsADAutoTrain: BaseModalTitleType = {
    title:'modal.title.ts.ad.autotrain.title', 
    icon: 'modal.title.ts.ad.autotrain.icon', 
    description:'modal.title.ts.ad.autotrain.description',
    size: 'lg',
    submitText: 'modal.button.ok.train'
}

export const tsADMultiTest: BaseModalTitleType = {
    title:'modal.title.ts.ad.multitest.title', 
    icon: 'modal.title.ts.ad.multitest.icon', 
    description:'modal.title.ts.ad.multitest.description',
    size: 'xl',
    submitText: 'modal.button.ok.test'
}

export const tsADAdditionalTrain: BaseModalTitleType = {
    title:'modal.title.ts.ad.additional.title', 
    icon: 'modal.title.ts.ad.additional.icon', 
    description:'modal.title.ts.ad.additional.description',
    size: 'xl',
    submitText: 'modal.button.ok.train'
}

/** Common modals */
export const trainDetail: BaseModalTitleType = {
    title:'modal.title.train.progress.title', 
    icon: 'modal.title.train.progress.icon', 
    description:'modal.title.train.progress.description',
    size: 'lg',
    submitText: ''
}

export const testDetail: BaseModalTitleType = {
    title:'modal.title.test.detail.title', 
    icon: 'modal.title.test.detail.icon', 
    description:'modal.title.test.detail.description',
    size: 'lg',
    submitText: ''
}

export const trialFailDetail: BaseModalTitleType = {
    title:'modal.title.trial.fail.title', 
    icon: 'modal.title.trial.fail.icon', 
    description:'modal.title.trial.fail.description',
    size: 'lg',
    submitText: ''
}

export const compareTrials: BaseModalTitleType = {
    title:'modal.title.comparetrials.title', 
    icon: 'modal.title.comparetrials.icon', 
    description:'modal.title.comparetrials.description',
    size: 'lg',
    submitText: ''
}

export const filterTrials: BaseModalTitleType = {
    title:'modal.title.filtertrials.title', 
    icon: 'modal.title.filtertrials.icon', 
    description:'modal.title.filtertrials.description',
    size: 'xl',
    submitText: 'modal.button.ok.apply'
}