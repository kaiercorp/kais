import { useEffect, useState, forwardRef, useImperativeHandle, useContext } from 'react'
import { Button, Card, Col, Form, InputGroup, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'

import { RadioGPUWithCPU, SelectDataModal } from 'features'
import { TableClsTrialDetail } from 'features/TrialDetail'
import { useToggle } from 'hooks'
import { emptyMultiSampleTest } from 'common'
import { objDeepCopy, ApiFetchTrial, ApiFetchTrainModels, ApiCreateTestDirectory } from 'helpers'
import { emptyTestMultiSampleErrors, validateTSADTest } from './TSADMultiSampleTestValidator'
import { engine } from 'appConstants/trial'
import { LabelSelect2 } from 'components'
import { TrialContext } from 'contexts'


const TSADMultiSampleTest = forwardRef(({ toggle, selectedTrial, selectedModels }: any, ref) => {
    const [t] = useTranslation('translation')

    const [requestData, setRequestData] = useState(objDeepCopy(emptyMultiSampleTest))
    const [formErrors, setFormErrors] = useState<emptyTestMultiSampleErrors>({ hasError: false })
    const [modelOptions, setModelOptions] = useState<any>([])
    const [directoryId, setDirectoryId] = useState<number | undefined>()
    const [trialId, setTrialId] = useState<number | undefined>()
    
    const { trialContextValue } = useContext(TrialContext)
    const { trial } = ApiFetchTrial(trialId)
    const { models } = ApiFetchTrainModels(trialId)
    const createTestDirectory = ApiCreateTestDirectory()

    const [isDataModalOpened, toggleDataModal] = useToggle()
    useImperativeHandle(ref, () => ({
        handleSubmit
    }))

    useEffect(() => {
        if (!models) return

        const modelOps = models.sort((a:any, b:any) => {
            return Number(a.score.String) - Number(b.score.String)
        }).map((model: any) => {
            return {
                value: `${model.train_uuid}_${model.name}`,
                label: `Trial${model.train_local_id}_${model.name} [${t(`metric.${model.target_metric.String}`)}: ${(Number(model.score.String).toFixed(3))}]`
            }
        })
        setModelOptions(modelOps)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [models])

    useEffect(() => {
        let newRequestData = objDeepCopy(emptyMultiSampleTest)
        if (selectedTrial) {
            setTrialId(selectedTrial.trial_id) 

            newRequestData['parent_id'] = selectedTrial.trial_id
            newRequestData['trial']['trial_id'] = selectedTrial.trial_id
            
            if (selectedModels.uuid && selectedModels.name) {
                const value = `${selectedModels.uuid}_${selectedModels.name}`
                const label = `Trial${selectedModels.train_local_id}_${selectedModels.name} [${t(`metric.${selectedModels.target_metric.String}`)}: ${(Number(selectedModels.score.String).toFixed(3))}]`
                newRequestData['model_list'] = [{ value: value, label: label }]
            }

            setRequestData(newRequestData)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const handleRequestData = (key: string, value: any) => {
        let newRequestData = objDeepCopy(requestData)
        if (key === 'trial_id' && value) {
            setTrialId(value)
            
            newRequestData['parent_id'] = value
            newRequestData['trial']['trial_id'] = value
        } else if (key === 'gpus') {
            if (value.length === 0) {
                newRequestData['gpu'] = []    
            } else {
                newRequestData['gpu'] = value
            }
        } else if (key === 'model') {
            if (value.action === 'select-option') {
                newRequestData['model_list'].push(value.option)
            } else if (value.action === 'remove-value') {
                const index = newRequestData['model_list'].indexOf(value.removedValue)
                newRequestData['model_list'].splice(index, 1)
            } else if (value.action === 'clear') {
                newRequestData['model_list'] = []
            }
        } else {
            newRequestData[key] = value
        }

        setRequestData(newRequestData)
    }

    const handleSelectData = (path: string, id: number) => {
        let newRequestData = objDeepCopy(requestData)

        newRequestData.data_path = path
        newRequestData.trial_name = 'test_' + path
        newRequestData['dataset_id'] = id
        setRequestData(newRequestData)
    }

    const handleSubmit = () => {
        let newRequestData = objDeepCopy(requestData)

        if(newRequestData['model_list']) {
            newRequestData['model_list'] = newRequestData.model_list.map((model: any) => {
              return model.value
            })
        }

        newRequestData['project_id'] = trial.project_id
        newRequestData['engine_type'] = engine.ts_ad
        newRequestData['params'] = JSON.stringify(newRequestData)

        const errors = validateTSADTest(newRequestData, t)
        setFormErrors(errors)

        if (!errors.hasError) {
            createTestDirectory.mutate(newRequestData)
            return true
        }

        return false
    }
    
    const handleDirectoryIdChange = (directoryId: number) => {
        setDirectoryId(directoryId)
    }

    return (
        <Form noValidate validated={formErrors.hasError}>
            <Row>
                <Col sm='6'>
                    <Card>
                        <Card.Header>{t('ui.test.title.model')}</Card.Header>
                        <Card.Body>
                            <Form.Group>
                                <Row>
                                    <Form.Label column='sm' sm={4}>
                                    {t('ui.test.select.train')}
                                    </Form.Label>
                                    <Col>
                                        {
                                            formErrors.hasError && formErrors.parent_id && formErrors.parent_id !== '' ? (
                                                <Form.Control.Feedback type='invalid' className='d-block'>
                                                    {formErrors.parent_id}
                                                </Form.Control.Feedback>
                                            )
                                                : null
                                        }
                                        <Form.Select
                                            name={'select-train-id'}
                                            size='sm'
                                            value={requestData.trial.trial_id}
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
                            </Form.Group>

                            <Form.Group>
                                <LabelSelect2
                                    title={t('ui.test.select.model')}
                                    name={'model_list'}
                                    options={modelOptions}
                                    value={requestData.model_list}
                                    onChange={(values: any, options: any) => handleRequestData('model', options)}
                                    errors={formErrors}
                                />
                            </Form.Group>
                        </Card.Body>
                    </Card>

                    <Card>
                        <Card.Header>{t('ui.test.title.config')}</Card.Header>
                        <Card.Body>
                            <RadioGPUWithCPU errors={formErrors} selectGPU={handleRequestData} />

                            <Form.Group>
                                <Row>
                                    <Form.Label column='sm' sm={4}>
                                    {t('ui.test.dataset')}
                                    </Form.Label>
                                    <Col>
                                        {
                                            formErrors.hasError && formErrors.dataset_id && formErrors.dataset_id !== '' ? (
                                                <Form.Control.Feedback type='invalid' className='d-block'>
                                                    {formErrors.dataset_id}
                                                </Form.Control.Feedback>
                                            )
                                                : null
                                        }
                                        <InputGroup className='mb-1'>
                                            <Form.Control value={requestData.data_path} readOnly />
                                            <Button variant='info' onClick={toggleDataModal}>
                                            {t('button.select')}
                                            </Button>
                                        </InputGroup>
                                    </Col>
                                </Row>
                            </Form.Group>
                        </Card.Body>
                    </Card>
                </Col>

                <Col sm='6'>
                    <Row>{trial && <TableClsTrialDetail trial={trial} />}</Row>
                </Col>
            </Row>

            <SelectDataModal show={isDataModalOpened} selectData={handleSelectData} toggle={toggleDataModal} isTest={true} dataType='table' directoryId={directoryId} onDirectoryIdChange={handleDirectoryIdChange}/>
        </Form>
    )
})

export default TSADMultiSampleTest