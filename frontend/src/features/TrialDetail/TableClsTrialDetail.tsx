import { Card, Col, Form, Row } from 'react-bootstrap'
import { convertUtcTime, getDuration } from 'helpers'
import { PopoverLabel } from 'components'
import { useTranslation } from 'react-i18next'

type Props = {
    trial: any
}

const TableClsTrialDetail = ({ trial }: Props) => {
    const [t] = useTranslation('translation')
    if (!trial.params || !trial.params.train_config) return <></>
    const config = trial.params.train_config
    
    return (
        <Card>
            <Card.Header>{t('ui.train.model.title')}</Card.Header>
            <Card.Body>
                <Col>
                    {trial.created_at && (
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.train.created')}
                            </Form.Label>
                            <Form.Label column='sm' sm={8}>
                                {convertUtcTime(trial.created_at)}
                            </Form.Label>
                        </Row>
                    )}
                    {trial.updated_at && (
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.train.duration')}
                            </Form.Label>
                            <Form.Label column='sm' sm={8}>
                                {getDuration(trial.created_at, trial.updated_at)}
                            </Form.Label>
                        </Row>
                    )}
                    <Row>
                        <Form.Label column='sm' sm={4}>
                            {t('ui.train.gpus')}
                        </Form.Label>
                        {trial.gpu_name && trial.gpu_name !== '' ? (
                            <Form.Label column='sm' sm={8}>
                                {trial.gpu_name}
                            </Form.Label>
                        ) : (
                            <Form.Label column='sm' sm={8}>
                                {t('ui.train.cpu')}
                            </Form.Label>
                        )}
                    </Row>
                    {config && config.target_metric && (
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.train.targetmetric')}
                            </Form.Label>
                            <Form.Label column='sm' sm={8}>
                                {t(`metric.${trial.target_metric}`)} <PopoverLabel>{t(`metric.${trial.target_metric}.desc`)}</PopoverLabel>
                            </Form.Label>
                        </Row>
                    )}
                    <Row>
                        <Form.Label column='sm' sm={4}>
                            {t('ui.train.data_path')}
                        </Form.Label>
                        <Form.Label column='sm' sm={8}>
                            {trial.params.data_path}
                        </Form.Label>
                    </Row>

                    {config && config.index_column && (
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.train.indexCol')}
                            </Form.Label>
                            <Form.Label column='sm' sm={8}>
                                {config.index_column}
                            </Form.Label>
                        </Row>
                    )}
                    {config && config.label_column && (
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.train.labelCol')}
                            </Form.Label>
                            <Form.Label column='sm' sm={8}>
                                {config.label_column}
                            </Form.Label>
                        </Row>
                    )}
                    {config && config.input_column && (
                        <Row>
                            <Form.Label column='sm' sm={4}>
                                {t('ui.train.includeCol')}
                            </Form.Label>
                            <Col sm={8}>
                                <Form.Select
                                    name='input-columns'
                                    id='input-columns'
                                    multiple
                                >
                                    {config.input_column && config.input_column.map((input: any) => {
                                        return <option key={`input-column-${input}`}>{input}</option>
                                    })}
                                </Form.Select>
                            </Col>
                        </Row>
                    )}
                </Col>
            </Card.Body>
        </Card>
    )
}

export default TableClsTrialDetail
