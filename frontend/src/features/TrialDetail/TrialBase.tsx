import { convertUtcTime, getDuration, logger } from 'helpers'
import { CustomConfirm, IconButtonWithPopover, StatusLabel } from 'components'
import { Col, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import { useEffect, useState } from 'react'
import { ApiStopTrain } from 'helpers'
import { useSocket } from 'hooks'

const TrialBase = ({ trial, tconfig, testState }: any) => {
    const [t] = useTranslation('translation')

    const [trainState, setTrainState] = useState('idle')
    const [targetMetric, setTargetMetric] = useState('')
    
    const stopTrain = ApiStopTrain()

    useEffect(() => {
        if (!trial) return
        setTrainState(trial.state)
    }, [trial])

    useEffect(() => {
        if (!tconfig) return
        let key = ''
        let val = 0
        if (typeof(tconfig.train_config.target_metric) === 'string') {
            key = tconfig.train_config.target_metric
        } else if (typeof(tconfig.train_config.target_metric) === 'object') {
            Object.keys(tconfig.train_config.target_metric).forEach((k:string) => {
                if (tconfig.train_config.target_metric[k] > val) {
                    key = k
                    val = tconfig.train_config.target_metric[k]
                }
            })
        }
        setTargetMetric(key)
    }, [tconfig])

    const handleSocketMessage = (e: MessageEvent<any>) => {
        try {
            let msg = JSON.parse(e.data)
            setTrainState(msg)
        } catch (e) {
            logger.error(e)
        }       
    }
    useSocket('/trials/state/' + trial?.trial_id, 'Trial State', handleSocketMessage, {shouldCleanup: true, shouldConnect: !!trial && trial.trial_id && trial.trial_id !== 0})

    const handleStopTrain = () => {
        CustomConfirm({
            onConfirm: () => {
                stopTrain.mutate(trial?.trial_id)
            },
            onCancel: () => { },
            message: t('ui.confirm.stop.train'),
        })
    }

    return (
        <table>
            <tbody>
                <tr>
                    <th>{t('ui.train.status')}</th>
                    <td>
                        <Row>
                            <Col>
                            {
                                testState
                                ?<StatusLabel state={testState}>{testState}</StatusLabel>
                                :<StatusLabel state={trainState}>{trainState}</StatusLabel>
                            }
                            </Col>
                            {
                                trial && trial.state && ['train', 'additional_train'].includes(trial.state) && (
                                    <Col sm='3'>
                                        <IconButtonWithPopover
                                            name={`btn-stop-train`}
                                            variant='danger'
                                            onClick={() => handleStopTrain()}
                                            popTitle={t('button.stop.train')}
                                            icon='mdi-stop-circle-outline'
                                        />
                                    </Col>
                                )
                            }
                        </Row>
                    </td>
                </tr>
                <tr>
                    <th>{t('ui.train.created')}</th>
                    <td>{convertUtcTime(trial.created_at)}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.duration')}</th>
                    <td>{getDuration(trial.created_at, trial.updated_at)}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.data_path')}</th>
                    <td>{tconfig.data_path}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.gpus')}</th>
                    <td>{trial.gpu_name ? trial.gpu_name.split(",").map((gpu: any) => {
                        return (
                            <div key={gpu}>{gpu}</div>
                        )
                    }) : 'CPU'}</td>
                </tr>
                <tr>
                    <th>{t('ui.train.targetmetric')}</th>
                    <td>{t(`metric.${targetMetric}`)}</td>
                </tr>
            </tbody>
        </table>
    )
}

export default TrialBase