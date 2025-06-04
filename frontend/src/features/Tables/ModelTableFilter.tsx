import { ResultOrder } from "common"
import { LabelInput, LabelSelect2 } from "components"
import { objDeepCopy } from "helpers"
import { useEffect, useState } from "react"
import { Col, Row, Form } from "react-bootstrap"
import { useTranslation } from "react-i18next"
import styled from "styled-components"

const FilterContainer = styled(Row)`
    padding: 5px;
    background-color: #343a40;
`

const emptyFilter = {
    trial_list: [],
    model_name: '',
    metric: {},
}

const ModelTableFilter = ({ allModels, setFilteredModels }: any) => {
    const [t] = useTranslation('translation')

    const [trials, setTrials] = useState<any>([])
    const [metrics, setMetrics] = useState<any>([])

    const [filter, setFilter] = useState<any>(emptyFilter)

    useEffect(() => {
        if (!allModels || allModels.length === 0) return

        const trialIds = new Map<number, number>()
        allModels.forEach((model: any) => {
            trialIds.set(model.train_local_id, model.train_local_id)
        })

        const modelOps = new Array<any>()
        trialIds.forEach((tid: any) => {
            modelOps.push({
                value: tid,
                label: `Trial${tid}`,
            })
        })

        modelOps.sort((a: any, b: any) => {
            return a.value - b.value
        })

        setTrials(modelOps)

        const result_keys = Object.keys(allModels[0].all_result).sort(function (x, y) {
            const xidx = ResultOrder[x] || 999
            const yidx = ResultOrder[y] || 999
            return xidx - yidx
        }).filter((key: string) => {
            if (key.includes('b_')) return false
            if (key.includes('c_')) return false
            return true
        })

        setMetrics([[...result_keys.slice(0, Math.floor(result_keys.length / 2))], [...result_keys.slice(Math.floor(result_keys.length / 2))]])
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [allModels])

    useEffect(() => {
        if (filter['trial_list'].length === 0
            && filter['model_name'] === ''
            && Object.keys(filter['metric']).length === 0) {
            setFilteredModels(allModels)
            return
        }

        const newModels = allModels.slice().filter((model: any) => {
            let returnVal = false

            filter['trial_list'].forEach((opt: any) => {
                if (opt.value === model.train_local_id) returnVal = true
            })

            if (filter['model_name'] !== '' && model.name.includes(filter['model_name'])) {
                returnVal = true
            }

            metrics.forEach((subMetrics: any) => {subMetrics.forEach((metric:any) => {
                const mod = ['wa', 'f1', 'precision', 'recall', 'uwa', 'image_accuracy', 'image_precision', 'image_recall', 'image_f1_score', 'label_accuracy', 'label_precision', 'label_recall', 'label_f1_score'].includes(metric)?100:1
                if (filter['metric'][metric] !== undefined) {
                    const min = Number(filter['metric'][metric]?.min?filter['metric'][metric].min:0)/mod
                    const max = Number(filter['metric'][metric]?.max?filter['metric'][metric].max:100)/mod
                    if (model['all_result'][metric] >= min && model['all_result'][metric] <= max) returnVal = true
                    else returnVal = false
                }
            })})

            return returnVal
        })

        setFilteredModels(newModels)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [filter])

    const filterTrial = (value: any) => {
        let newFilter = objDeepCopy(filter)
        if (value.action === 'select-option') {
            newFilter['trial_list'].push(value.option)
        } else if (value.action === 'remove-value') {
            const index = newFilter['trial_list'].indexOf(value.removedValue)
            newFilter['trial_list'].splice(index, 1)
        } else if (value.action === 'clear') {
            newFilter['trial_list'] = []
        }
        setFilter(newFilter)
    }

    const filterModel = (value: any) => {
        let newFilter = objDeepCopy(filter)
        newFilter['model_name'] = value
        setFilter(newFilter)
    }

    const filterMetric = (metric: string, type: string, value: any) => {
        let newFilter = objDeepCopy(filter)
        if (newFilter['metric'].hasOwnProperty(metric)) {
            newFilter['metric'][metric][type] = value
        } else {
            if (type === 'min') newFilter['metric'][metric] = { min: value }
            else newFilter['metric'][metric] = { max: value }
        }

        setFilter(newFilter)
    }

    return (
        <FilterContainer>
            <Col>
                <Row>
                    <Col sm={3} style={{paddingRight: '35px', marginLeft: '10px'}}>
                        <LabelSelect2
                            title={'Trial'}
                            name={'filter_trial_list'}
                            options={trials}
                            value={filter.trial_list}
                            onChange={(values: any, options: any) => filterTrial(options)}
                        />
                    </Col>

                    <Col sm={3} style={{padding: '0 45px 0 5px'}}>
                        <LabelInput
                            title={'Model'}
                            name='filter_model_name'
                            value={filter.model_name}
                            onChange={(e: any) => filterModel(e.target.value)}
                            errors={null}
                        />
                    </Col>
                </Row>

                <Row>
                    {
                        metrics.map((subMetric: any) => {
                            return (
                                <Row>
                                    {subMetric.map((metric: string, index: number) => {
                                    return (
                                        <Col style={{ margin: '10px'}} key={`filter-metric-${index}`}>
                                            <Row>
                                                <Col>
                                                    <Form.Label>{t(`metric.${metric}`)}{['wa', 'f1', 'precision', 'recall', 'uwa', 'image_accuracy', 'image_precision', 'image_recall', 'image_f1_score', 'label_accuracy', 'label_precision', 'label_recall', 'label_f1_score'].includes(metric)?'(%)':''}</Form.Label>
                                                </Col>
                                            </Row>
                                            <Row>
                                                <Col>
                                                    <Form.Control
                                                        required
                                                        type={'number'}
                                                        size='sm'
                                                        value={filter['metric'][metric]?.min ? filter['metric'][metric]['min'] : 0}
                                                        min={0}
                                                        max={100}
                                                        onChange={(e: any) => filterMetric(metric, 'min', e.target.value)}
                                                    />

                                                </Col>

                                                <Col xs={1}>~</Col>

                                                <Col>
                                                    <Form.Control
                                                        required
                                                        type={'number'}
                                                        size='sm'
                                                        value={filter['metric'][metric]?.max ? filter['metric'][metric]['max'] : 100}
                                                        min={0}
                                                        max={100}
                                                        onChange={(e: any) => filterMetric(metric, 'max', e.target.value)}
                                                    />
                                                </Col>
                                            </Row>
                                        </Col>
                                    )})}
                                </Row>
                            )
                        })
                    }
                </Row>
            </Col>
        </FilterContainer>
    )
}

export default ModelTableFilter