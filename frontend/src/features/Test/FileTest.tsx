import { forwardRef, useEffect, useState, useContext } from 'react'
import { Card, Col, Form, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import { useQueryClient } from '@tanstack/react-query'
import styled from 'styled-components'

import { objDeepCopy, ApiFetchTrial, ApiFetchTrainModels, ApiCreateTestFile, QUERY_KEY } from 'helpers'
import { emptyTestFile } from 'common'
import { engine } from 'appConstants/trial'
import { TrialContext } from 'contexts'

import { FileUploader, LabelSelect2, Spinner } from 'components'
import { RadioGPU, VisionClsTrialDetail, VisionADTrialDetail, FileTestResult } from 'features'

const SpinnerContainer = styled.div`
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
`

const FileTest = forwardRef(({ toggle, engineType }: any, ref) => {
    const [t] = useTranslation('translation')

    const [trialId, setTrialId] = useState<number | undefined>()
    const [requestData, setRequestData] = useState(objDeepCopy(emptyTestFile))
    const [imageDisabled, setImageDisabled] = useState(false)
    const [disabledMsg, setDisabledMsg] = useState('')

    const { trialContextValue } = useContext(TrialContext)

    const { isLoading: isTrialLoading, trial } = ApiFetchTrial(trialId)
    const { isLoading: isModelsLoading, models } = ApiFetchTrainModels(trialId)
    const createTestFile = ApiCreateTestFile()
    const queryClient = useQueryClient()
    const testFileResult = queryClient.getQueryData<any>([QUERY_KEY.createTestFile])

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
        if (!requestData.file) {
            return
        }

        let newRequest = objDeepCopy(requestData)

        newRequest['model_list'] = [newRequest.model_list.value]

        let params = {
            'gpus': objDeepCopy(newRequest.gpus),
            'model_list': objDeepCopy(newRequest.model_list),
            'parent_id': objDeepCopy(newRequest.parent_id)
        }

        newRequest['engine_type'] = engineType 
        newRequest['params'] = JSON.stringify(params)
        newRequest.file = requestData.file

        createTestFile.mutate(newRequest)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [requestData.file])

    const [modelOptions, setModelOptions] = useState<any>([])
    useEffect(() => {
        if (!models) return

        const modelOps = models.sort((a:any, b:any) => {
            return Number(b.score.String) - Number(a.score.String)
        }).map((model: any) => {
            return {
                value: `${model.train_uuid}_${model.name}`,
                label: `Trial${model.train_id}_${model.name} [${t(`metric.${model.target_metric.String}`)}: ${(parseFloat(model.score.String.replaceAll("\"", "")) * 100).toFixed(2)}%]`
            }
        })
        setModelOptions(modelOps)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [models])

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
                },
                file: null
            }))
        } else if (key === 'model') {
            setRequestData((requestData: any) => ({
                ...requestData,
                model_list : value,
                file: null 
            }))
        } else if (key === 'file') {
            setRequestData((requestData: any) => ({
                ...requestData,
                [key] : value
            }))
        } else {
            setRequestData((requestData: any) => ({
                ...requestData,
                [key] : value,
                file: null
            }))
        }
    }

    return (
        <Row>
            <Col xs={4}>
                <Card>
                    <Card.Header>{t('ui.test.title.model')}</Card.Header>
                    {trial &&
                    <Card.Body>
                        <Form.Group>
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
                            <Row>
                                {requestData.parent_id && 
                                    (engineType === engine.vision_ad ?
                                        <VisionADTrialDetail trial={trial} />
                                        :
                                        <VisionClsTrialDetail trial={trial} />           
                                    )
                                }
                            </Row>
                        </Form.Group>

                        <Form.Group>
                            <LabelSelect2
                                title={t('ui.test.select.model')}
                                name={'multiselect-model'}
                                options={modelOptions}
                                value={requestData.model_list}
                                onChange={(options: any) => handleRequestData('model', options)}
                                isMulti={false}
                            />
                        </Form.Group>
                    </Card.Body>
                    }
                </Card>

                <Card>
                    <Card.Header>{t('ui.test.title.config')}</Card.Header>
                    <Card.Body>
                        <RadioGPU selectGPU={handleRequestData} />

                        <FileUploader
                            onFileUpload={(files) => { handleRequestData('file', files[0]) }}
                            disabled={imageDisabled}
                            disabledMsg={disabledMsg}
                        />
                    </Card.Body>
                </Card>
            </Col>
            <Col xs={8}>
                <Card style={{height: '100%'}}>
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
                                <Spinner className='spinner-border-sm' />
                            </SpinnerContainer>
                        </Card.Body>
                    ) : (
                        <Card.Body>
                            {testFileResult && <FileTestResult data={testFileResult} engineType={engineType}/>}
                        </Card.Body>
                    )}
                </Card>
            </Col>
        </Row>
    )
})

export default FileTest