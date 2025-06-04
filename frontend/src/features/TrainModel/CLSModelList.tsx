import { useEffect, useState } from 'react'
import { Col, Row } from 'react-bootstrap'

import { ApiFetchTrainModels } from 'helpers'
import { engine } from 'appConstants/trial'

import ModelTable from 'features/Tables/ModelTable'
import { ColumnModels } from 'features/Tables/ColumnModels'
import { ModelTableFilter } from 'features/Tables'
import CFMatrix from './CFMatrix'
import MLMatrix from 'pages/Vision/Cls.ML/MLMatrix'


const CLSModelList = ({ trial, openMultiTest, downloadModel, engineType }: any) => {
    const [matrix, setMatrix] = useState<any>()
    const [allModels, setAllModels] = useState<any>([])
    const [filteredModels, setFilteredModels] = useState<any>([])

    const { models } = ApiFetchTrainModels(trial.trial_id)

    useEffect(() => {
        if (!models) return

        const refinedModels = models ? models.slice().map((model: any) => {
            model['model_id'] = "trial" + model.train_local_id + "_" + model.name

            if (model.hasOwnProperty('all_result') && model['all_result'].Valid) {
                const result_keys = Object.keys(model.all_result.String)
                result_keys.forEach((key: string) => {
                    if (key === 'model' || key === 'epoch') return
                    if (key.startsWith('c_') || key.startsWith('b_')) return
                    model[key] = Number(model.all_result[key] || 0).toFixed(2)
                })
            }

            if (model.hasOwnProperty('class_result')) {
                const class_keys = Object.keys(model.class_result)
                class_keys.forEach((key: string) => {
                    if (key === 'model') return
                    const inner_key = Object.keys(model.class_result[key])
                    if (!model.class_result[key][inner_key[0]]) return
                    const vals = model.class_result[key][inner_key[0]].split(",")
                    model[`${key}_acc`] = Number(vals[0])
                    model[`${key}_correct`] = vals[1]
                })
            }
            
            if (model.accuracy_info) {
                const accuracyInfoKeys = Object.keys(model.accuracy_info)
                accuracyInfoKeys.forEach((key: string) => {
                    const innerKeys = Object.keys(model.accuracy_info[key])
                    innerKeys.forEach((innerKey) => {
                        if (!model.accuracy_info[key][innerKey]) return
                        const vals = model.accuracy_info[key][innerKey].split(",")
                        model[`${key}_acc`] = Number(vals[0])
                        model[`correct_${key}_over_total`] = vals[1]
                    })
                })
            }

            return model
        }).sort(function (a: any, b: any) {
            let direction = -1
            if (['mse', 'mae', 'rmse'].includes(a['target_metric']['String'])) {
                direction = direction * -1
            }
            return direction * (Number(a['all_result'][a['target_metric']['String']]) - Number(b['all_result'][a['target_metric']['String']]))
        }) : []

        setAllModels(refinedModels)
        setFilteredModels(refinedModels)
    }, [models])

    const onSelect = (row: any) => {
        const matrix = engineType === engine.vision_cls_ml ? row.ml_matrix.String : row.cf_matrix.String
        setMatrix(matrix)
    }

    return (
        <Col>
            <ModelTableFilter allModels={allModels} setFilteredModels={setFilteredModels} />
            <Row>
                <Col>
                    <ModelTable
                        CustomColumn={ColumnModels}
                        filteredModels={filteredModels}
                        onSelect={onSelect}
                        openMultiTest={openMultiTest}
                        downloadModel={downloadModel}
                    />
                </Col>
            </Row>
            <Row>
                <Col>
                    {matrix && 
                        (engineType === engine.vision_cls_ml ?                                                 
                        <MLMatrix ml_matrixStr={matrix}/> : <CFMatrix cf_matrixStr={matrix} />                       
                    )}
                </Col>
            </Row>
        </Col>
    )
}

export default CLSModelList
