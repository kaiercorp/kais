import { useEffect, useState } from 'react'
import { Col, Row } from 'react-bootstrap'

import ModelTable from 'features/Tables/ModelTable'
import { ColumnModels } from 'features/Tables/ColumnModels'
import { logger } from 'helpers'
import TableErrorGraphImage from './TableErrorGraphImage'
import { ApiFetchTrainModels } from 'helpers'
import { useSocket } from 'hooks'
import { ModelTableFilter } from 'features/Tables'


const RegModelList = ({ trial, openMultiTest, downloadModel }: any) => {
    const [errGraph, setErrorGraph] = useState<any>()
    const [allModels, setAllModels] = useState<any>([])
    const [filteredModels, setFilteredModels] = useState<any>([])
    const [socketConnected, setSocketConnected] = useState(false)

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

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            if (e.data === 'Invalid request') {
                throw e
            }

            let msg = JSON.parse(e.data)
            if (msg.hasOwnProperty('error_graph')) {
                setErrorGraph(msg.error_graph)
            }
        } catch (e) {
            logger.error(e)
        }
    }
    const ws = useSocket('/trials/train/result/treg', 'Result Data', handleSocketMessage, { setSocketConnected, shouldCleanup: true, shouldConnect: !!trial })

    const onSelect = (row: any) => {
        if (!socketConnected || !ws || !ws.current) return
        ws.current.send(JSON.stringify({ train_id: row.train_id, model_name: row.name }))
    }

    return (
        <Col>
            <ModelTableFilter allModels={allModels} setFilteredModels={setFilteredModels} />
            <Row>
                <ModelTable
                    CustomColumn={ColumnModels}
                    filteredModels={filteredModels}
                    onSelect={onSelect}
                    openMultiTest={openMultiTest}
                    downloadModel={downloadModel}
                />
            </Row>
            <Row>
                <Col>
                    {errGraph && <TableErrorGraphImage error_graph={errGraph} />}
                </Col>
            </Row>
        </Col>
    )
}

export default RegModelList
