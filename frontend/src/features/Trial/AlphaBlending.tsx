import { PopoverLabel, StatusRow } from "components"
import { Col, Form, Row } from "react-bootstrap"
import { useTranslation } from "react-i18next"


const AlphaBlending = ({ requestData, handleRequestData }: any) => {
    const [t] = useTranslation('translation')

    return (
        <Form.Group className={'mb-1'}>
            <Row>
                <Form.Label column='sm' sm={4}>
                    <span>
                        {t('ui.train.targetmetric')}
                        <PopoverLabel>
                            <StatusRow><span>{t(`metric.uwa.desc`)}</span></StatusRow>
                            <StatusRow><span>{t(`metric.wa.desc`)}</span></StatusRow>
                            <StatusRow><span>{t(`metric.precision.desc`)}</span></StatusRow>
                            <StatusRow><span>{t(`metric.recall.desc`)}</span></StatusRow>
                            <StatusRow><span>{t(`metric.f1.desc`)}</span></StatusRow>
                        </PopoverLabel>
                    </span>
                </Form.Label>
                <Col column='sm' sm={8}>
                    <Row>
                        <Col column='sm' sm={8}>
                            <span style={{ fontSize: '12px' }}>{t(`metric.wa`)}(%)</span>
                        </Col>
                        <Col column='sm' sm={4}>
                            <input
                                type='number'
                                name='metric.wa'
                                min='0' max='100'
                                value={requestData.train_config.target_metric.wa}
                                onChange={(e) => handleRequestData('atarget_metric.wa', e.target.value)}
                            />
                        </Col>
                    </Row>
                    <Row>
                        <Col column='sm' sm={8}>
                            <span style={{ fontSize: '12px' }}>{t(`metric.uwa`)}(%)</span>
                        </Col>
                        <Col column='sm' sm={4}>
                            <input
                                type='number'
                                name='metric.uwa'
                                min='0' max='100'
                                value={requestData.train_config.target_metric.uwa}
                                onChange={(e) => handleRequestData('atarget_metric.uwa', e.target.value)}
                            />
                        </Col>
                    </Row>
                    <Row>
                        <Col column='sm' sm={8}>
                            <span style={{ fontSize: '12px' }}>{t(`metric.precision`)}(%)</span>
                        </Col>
                        <Col column='sm' sm={4}>
                            <input
                                type='number'
                                name='metric.precision'
                                min='0' max='100'
                                value={requestData.train_config.target_metric.precision}
                                onChange={(e) => handleRequestData('atarget_metric.precision', e.target.value)}
                            />
                        </Col>
                    </Row>
                    <Row>
                        <Col column='sm' sm={8}>
                            <span style={{ fontSize: '12px' }}>{t(`metric.recall`)}(%)</span>
                        </Col>
                        <Col column='sm' sm={4}>
                            <input
                                type='number'
                                name='metric.recall'
                                min='0' max='100'
                                value={requestData.train_config.target_metric.recall}
                                onChange={(e) => handleRequestData('atarget_metric.recall', e.target.value)}
                            />
                        </Col>
                    </Row>
                    <Row>
                        <Col column='sm' sm={8}>
                            <span style={{ fontSize: '12px' }}>{t(`metric.f1`)}(%)</span>
                        </Col>
                        <Col column='sm' sm={4}>
                            <input
                                type='number'
                                name='metric.f1'
                                min='0' max='100'
                                value={requestData.train_config.target_metric.f1}
                                onChange={(e) => handleRequestData('atarget_metric.f1', e.target.value)}
                            />
                        </Col>
                    </Row>
                </Col>
            </Row>
        </Form.Group>
    )
}

export default AlphaBlending