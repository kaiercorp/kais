import { forwardRef, useEffect, useState, useContext } from 'react'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import { Button, Card, Col, Form, Row } from 'react-bootstrap'
import styled from 'styled-components'

import { engine } from 'appConstants/trial'
import { TrialContext } from 'contexts'
import { emptyTestFile } from 'common'
import { objDeepCopy, 
    ApiFetchTrial, 
    ApiFetchTrainModels, 
    ApiCreateTestFile, 
    ApiGetRowFromFile, 
    QUERY_KEY 
} from 'helpers'

import { FileUploader, LabelSelect2, Spinner } from 'components'
import { DataTable, 
    RadioGPUWithCPU, 
    TableClsTrialDetail, 
    PredictionImportanceChart, 
    PredictionProbabilities,
    PredictionImportanceChartREG,
    PredictionProbabilitiesREG
} from 'features'

const SpinnerContainer = styled.div`
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
` 

const SingleSampleTest = forwardRef(({ toggle, engineType }: any, ref) => {
    const [t] = useTranslation('translation')

    const [trialId, setTrialId] = useState<number | undefined>()
    const [requestData, setRequestData] = useState(objDeepCopy(emptyTestFile))
    const [imageDisabled, setImageDisabled] = useState(false)
    const [disabledMsg, setDisabledMsg] = useState('')
    const [selectedRow, setSelectedRow] = useState<Number>(0)
    const [modelOptions, setModelOptions] = useState<any>([])

    const { trialContextValue } = useContext(TrialContext)
    
    const { isLoading: isTrialLoading, trial } = ApiFetchTrial(trialId)
    const { isLoading: isModelsLoading, models } = ApiFetchTrainModels(trialId)
    const createTestFile = ApiCreateTestFile()
    const getRowFromFile = ApiGetRowFromFile()
    const queryClient = useQueryClient()
    const testFileResult = queryClient.getQueryData<any>([QUERY_KEY.createTestFile])
    const rows = queryClient.getQueryData<any>([QUERY_KEY.getRowFromFile])

    useEffect(() => {
        setRequestData(objDeepCopy(emptyTestFile))
    }, [])

    useEffect(() => {
        if (requestData.parent_id < 1) {
            setImageDisabled(true)
            setDisabledMsg(t('validator.project'))
        } else if (!requestData.model_list || requestData.model_list.length < 1) {
            setImageDisabled(true)
            setDisabledMsg(t('validator.model'))
        } else if (!requestData.gpus) {
            setImageDisabled(true)
            setDisabledMsg(t('validator.gpu'))
        } else {
            setImageDisabled(false)
            setDisabledMsg('')
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [requestData])

    useEffect(() => {
        if (!models) return

        const modelOps = models.sort((a:any, b:any) => { 
            return engineType === engine.table_cls ? 
                Number(b.score.String) - Number(a.score.String)
                : Number(a.score.String) - Number(b.score.String)
        }).map((model: any) => {
            return {
                value: `${model.train_uuid}_${model.name}`,
                label: engineType === engine.table_cls ? 
                `Trial${model.train_id}_${model.name} [${t(`metric.${model.target_metric.String}`)}: ${(parseFloat(model.score.String.replaceAll("\"", "")) * 100).toFixed(2)}%]`
                : `Trial${model.train_local_id}_${model.name} [${t(`metric.${model.target_metric.String}`)}: ${(Number(model.score.String).toFixed(3))}]`
            }
        })
        setModelOptions(modelOps)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [models])

    const startTest = () => {
        if (!rows) return

        let newRequest = objDeepCopy(requestData)

        newRequest['model_list'] = [newRequest.model_list.value]
      
        let params = {
            'gpus': objDeepCopy(newRequest.gpus),
            'model_list': objDeepCopy(newRequest.model_list),
            'parent_id': objDeepCopy(newRequest.parent_id),
            'row_index': selectedRow,
            'filepath': rows.filepath
        }

        newRequest['engine_type'] = engineType
        newRequest['filepath'] = rows.filepath
        newRequest['params'] = JSON.stringify(params)

        createTestFile.mutate(newRequest)
    }

    const handleRequestData = (key: string, value: any) => {
        if (key === 'trial_id') {
            if (value === 0) {
                value = undefined
            }

            setTrialId(value)

            setRequestData((requestData: any) => ({
                ...requestData,
                parent_id : value,
                trial: {
                    ...requestData.trial,
                    trial_id: value
                } 
            }))
        } else if (key === 'model') {
            setRequestData((requestData: any) => ({
                ...requestData,
                model_list : value 
            }))
        } else {
            setRequestData((requestData: any) => ({
                ...requestData,
                [key] : value 
            }))
        }
    }

    const handleRequestFile = (file: any) => {
        let request = {file: file}
        getRowFromFile.mutate(request)
    }

    const selectRow = (index:Number) => {
        setSelectedRow(index)
    }

    return (
        <Row>
            <Col xs={5}>
                <Card>
                    <Card.Header>{t('ui.test.title.model')}</Card.Header>
                    {trial && 
                    <Card.Body>
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.test.select.train')}
                            </Form.Label>
                            <Col>
                                <Form.Select
                                    name={'select-train-id'}
                                    size='sm'
                                    value={requestData.trial.parent_id}
                                    onChange={(e) => handleRequestData('trial_id', Number(e.target.value))}
                                    className='mb-1'
                                >
                                    <option value=''>{t('ui.label.let.select')}</option>
                                    {trialContextValue.trials && 
                                        trialContextValue.trials.filter((trial: any) => {
                                            if (!trial.state) return false
                                            if (
                                                trial.state === 'finish' ||
                                                (trial.state === 'finish-fail' && trial.best_model_download_path)
                                            ) {
                                                return true
                                            }
                                            return false
                                        })
                                        .sort(function (a: any, b: any) {
                                            return a.trial_id - b.trial_id
                                        })
                                        .map((trial: any) => {
                                            return (
                                                <option key={`select-trial-option-${trial.trial_id}`} value={trial.trial_id}>
                                                    [ID: {trial.trial_id}] {trial.trial_name}
                                                </option>
                                            )
                                        })}
                                </Form.Select>
                            </Col>
                        </Row>
                        <Row>{requestData.parent_id && <TableClsTrialDetail trial={trial} />}</Row>

                        <Row>
                            <LabelSelect2
                                title={t('ui.test.select.model')}
                                name={'multiselect-model'}
                                options={modelOptions}
                                value={requestData.model_list}
                                onChange={(options: any) => handleRequestData('model', options)}
                                isMulti={false}
                            />
                        </Row>


                    </Card.Body>
                    }
                </Card>

                <Card>
                    <Card.Body>
                        <Row>
                            <Col>
                                <RadioGPUWithCPU selectGPU={handleRequestData} />
                            </Col>
                        </Row>
                        <Row>
                            <FileUploader
                                onFileUpload={(files) => { handleRequestFile(files[0]) }}
                                disabled={imageDisabled}
                                disabledMsg={disabledMsg}
                                accept='text/csv, application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
                            />
                        </Row>
                    </Card.Body>
                </Card>

            </Col>
            <Col xs={7}>
                <Card style={ getRowFromFile.isPending ? {height: '50%'} : {}}>
                    <Card.Header>{t('ui.test.title.rows')}</Card.Header>
                    {getRowFromFile.isPending && (
                        <Card.Body style={{height: '100%', position: 'relative'}}>
                            <SpinnerContainer>
                               <Spinner className='spinner-border-sm'/> 
                            </SpinnerContainer>
                        </Card.Body>
                    )}
                    {getRowFromFile.isSuccess && (
                        <>
                            <DataTable dataset={rows?.rows} selectedData={selectedRow} selectData={selectRow} />
                            <div style={{width:'200px'}}>
                                <Button variant='success' className='btn-rounded btn-sm' onClick={startTest} >
                                    Run
                                </Button>
                            </div>
                        </>
                    )}
                </Card>
                <Card style={ createTestFile.isPending ? {height: '50%'} : {}}>
                    {!isTrialLoading && !isModelsLoading ? (
                        <Card.Header>{t('ui.test.title.result')}</Card.Header>
                    ) : (
                        <Card.Header>
                            <Spinner className='spinner-border-sm' />
                            {t('ui.test.wait')}
                        </Card.Header>
                    )}
                    {createTestFile.isPending ? (
                        <Card.Body style={{height: '100%', position: 'relative'}}>
                            <SpinnerContainer>
                               <Spinner className='spinner-border-sm'/> 
                            </SpinnerContainer>
                        </Card.Body>
                    ) : (
                        <Card.Body>
                            {engineType === engine.table_cls ? (
                                    <>
                                    {testFileResult && <PredictionProbabilities data={testFileResult} />}
                                    {testFileResult && <PredictionImportanceChart data={testFileResult.visualization.feature_importances} toplabel={testFileResult.visualization.top_label} />}
                                    </>
                                ) : (
                                    <>
                                    {testFileResult && <PredictionProbabilitiesREG data={testFileResult} />}
                                    {testFileResult && <PredictionImportanceChartREG data={testFileResult.visualization.feature_importances} />}
                                    </>
                                )
                            }
                        </Card.Body>
                    )}
                </Card>
            </Col>
        </Row>
    )
})

export default SingleSampleTest