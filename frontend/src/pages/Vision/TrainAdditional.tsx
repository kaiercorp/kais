import { forwardRef, useImperativeHandle, useState, useContext } from "react"
import { useTranslation } from 'react-i18next'
import { Card, Col, Form, Row } from "react-bootstrap"

import { convertUtcTime, objDeepCopy, ApiCreateTrain, ApiFetchTrial } from "helpers"
import { TrialContext } from 'contexts'
import { engine } from 'appConstants/trial'

import { emptyErrors as clsEmptyErrors, validateCLS } from "pages/Vision/Cls.SL/TrainValidator"
import { emptyErrors as adEmptyErrors, validateAD} from 'pages/Vision/AD/TrainValidator'
import { RadioGPU } from "features"

interface TrainAdditionalProps {
    trialId: number
    engineType: typeof engine
}

const TrainAdditional = forwardRef(({trialId, engineType}: TrainAdditionalProps, ref) => {
    const [t] = useTranslation('translation')

    const { trialContextValue, updateTrialContextValue } = useContext(TrialContext)
    const { requestData, trainMode } = trialContextValue
    
    const createTrain = ApiCreateTrain()
    let { trial } = ApiFetchTrial(trialId)

    const handleRequestData = (key: string, value: any) => {
        if (typeof requestData === 'undefined') return 

        let newRequestData = objDeepCopy(requestData)

        if (['auto_stop'].includes(key)) {
            newRequestData.train_config[key] = value
        } else {newRequestData[key] = value}

        updateTrialContextValue({requestData: newRequestData})
    }

    const setConfig = () => {
        if (typeof requestData === 'undefined') return 
        if (typeof trainMode === 'undefined') return 

        let newRequestData = objDeepCopy(requestData)
        newRequestData.train_once = false
        newRequestData.train_config.epochs = -1
        newRequestData.train_config.parent_id = trial.trial_id

        return {
            trial_id: trial.trial_id,
            train_type: trainMode,
            params: JSON.stringify(newRequestData)
        }
    }

    const [formErrors, setFormErrors] = useState<clsEmptyErrors | adEmptyErrors>({hasError: false})
    useImperativeHandle(ref, () => ({
        handleSubmit
    }))
    const handleSubmit = () => {
        const errors = engineType === engine.vision_ad ? validateAD(trainMode, requestData, t) : validateCLS(trainMode, requestData, t)
        setFormErrors(errors)

        if (!errors.hasError) {
            createTrain.mutate(setConfig())
            return true
        }
    }

    return (
        <Form noValidate validated={formErrors.hasError}>
            <Row>
                { trial &&
                 <Card>
                     <Card.Header>{t('ui.train.title.infoorigin')}</Card.Header>
                     <Card.Body>
                         <Form.Group>
                             <Row>
                                 <Form.Label column='sm' sm='4'>
                                     {t('ui.train.name')}
                                 </Form.Label>
                                 <Col xs='8'>{trial.trial_name}</Col>
                             </Row>
                         </Form.Group>
                         <Form.Group>
                             <Row>
                                 <Form.Label column='sm' sm='4'>
                                     {t('ui.train.name.testdata')}
                                 </Form.Label>
                                 <Col xs='8'>{trial?.params?.test_db}</Col>
                             </Row>
                         </Form.Group>
                         <Form.Group>
                             <Row>
                                 <Form.Label column='sm' sm='4'>
                                     {t('ui.formatter.accuracy')}
                                 </Form.Label>
                                 <Col xs='8'>{trial.accuracy > 0 ? (trial.accuracy * 100).toFixed(2) : 0}%</Col>
                             </Row>
                         </Form.Group>
                         <Form.Group>
                             <Row>
                                 <Form.Label column='sm' sm='4'>
                                     {t('ui.formatter.inftime')}
                                 </Form.Label>
                                 <Col xs='8'>{trial.inference_time > 0 ? (trial.inference_time * 1000).toFixed(0) : 0}ms</Col>
                             </Row>
                         </Form.Group>
                         <Form.Group>
                             <Row>
                                 <Form.Label column='sm' sm='4'>
                                     {t('ui.train.created')}
                                 </Form.Label>
                                 <Col xs='8'>{convertUtcTime(trial.created_at)}</Col>
                             </Row>
                         </Form.Group>

                         <RadioGPU selectGPU={handleRequestData} errors={formErrors} />

                         <Form.Group className={'mb-1'}>
                             <Row>
                                 <Form.Label column='sm' sm={4}>{t('ui.train.autostop')}</Form.Label>
                                 <Col column='sm' sm={8} style={{ marginTop: '5px' }}>
                                     <Form.Switch
                                         type='switch'
                                         checked={requestData ? requestData.train_config.auto_stop : true}
                                         label={<Form.Label>{requestData && requestData.train_config.auto_stop ? t('ui.train.autostop.auto') : t('ui.train.autostop.user')}</Form.Label>}
                                         onChange={() => handleRequestData('auto_stop', requestData? !requestData.train_config.auto_stop : false)}
                                     />
                                 </Col>
                             </Row>
                         </Form.Group>
                     </Card.Body>
                 </Card>
                }
            </Row>
        </Form>
    )
})

export default TrainAdditional