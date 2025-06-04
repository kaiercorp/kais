import { Col, Form, Row } from 'react-bootstrap'
import { convertUtcTime, getDuration } from 'helpers'
import { useTranslation } from 'react-i18next'
import { PopoverLabel } from 'components'

type Props = {
  trial: any
}

const VisionClsTrialDetail = ({ trial }: Props) => {
  const [t] = useTranslation('translation')
  if (!trial || !trial.params || !trial.params.train_config) return <></>
  const config = trial.params.train_config
  const width = config.width === -1 ? t('ui.train.image.resolution.origin') : `${config.width} x `
  const height = config.height === -1 ? '' : config.height
  const imgSize = `${width}${height}`
  return (
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
      {trial.gpu_name && (
        <Row>
          <Form.Label column='sm' sm={4}>
          {t('ui.train.gpus')}
          </Form.Label>
          <Form.Label column='sm' sm={8}>
            {trial.gpu_name}
          </Form.Label>
        </Row>
      )}
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
      {config && config.default_config_file && (
        <Row>
          <Form.Label column='sm' sm={4}>
            {t('ui.train.config_file')}
          </Form.Label>
          <Form.Label column='sm' sm={8}>
            {config.default_config_file}
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
      {config && config.width && config.height && (
        <Row>
          <Form.Label column='sm' sm={4}>
            {t('ui.train.image.title')}
          </Form.Label>
          <Form.Label column='sm' sm={8}>
            {imgSize}
          </Form.Label>
        </Row>
      )}
      {/* {config && config.base_lr && (
        <Row>
          <Form.Label column='sm' sm={4}>
            {t('ui.train.base_lr')}
          </Form.Label>
          <Form.Label column='sm' sm={8}>
            {(config.base_lr || 0).toFixed(6)}
          </Form.Label>
        </Row>
      )}
      {config && config.epochs > 0 && (
        <Row>
          <Form.Label column='sm' sm={4}>
            {t('ui.train.epochs')}
          </Form.Label>
          <Form.Label column='sm' sm={8}>
            {config.epochs}
          </Form.Label>
        </Row>
      )} */}
      {/* <Row>
        <Form.Label column='sm' sm={4}>
          {t('ui.train.batch_size')}
        </Form.Label>
        <Form.Label column='sm' sm={8}>
          {config.train_batch_size < 1 ? <>{t('ui.train.max')}</> : config.train_batch_size}
        </Form.Label>
      </Row> */}
    </Col>
  )
}

export default VisionClsTrialDetail
