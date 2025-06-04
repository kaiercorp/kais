import { PopoverLabel } from 'components'
import { Col, Form, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'

type Props = {
  width: number
  height: number
  errors: any
  selectResolution: (name: string, value: number | string) => void
}

const SelectImageResolution = ({ width, height, errors, selectResolution }: Props) => {
  const [t] = useTranslation('translation')
  const handleResolution = (e: any) => {
    selectResolution('resolution', e.target.value)
  }

  return (
    <div>
      <Form.Group>
        <Row>
          <Form.Label column='sm' sm={4}>
            {t('ui.train.image.title')}
            <PopoverLabel name='description-image-resolution'>
            {t('ui.train.image.title.desc')}
            </PopoverLabel>
          </Form.Label>
          <Col>
            <Form.Select
              name={'select-image-resolution'}
              size='sm'
              value={`${width}x${height}`}
              onChange={handleResolution}
            >
              <option key='resolution-input' value='1x1'>
              {t('ui.train.image.resolution.user')}
              </option>
              <option key='resolution-origin' value='-1x-1'>
              {t('ui.train.image.resolution.origin')}
              </option>
              <option key='resolution-512' value='512x512'>
              {t('ui.train.image.resolution.512')}
              </option>
              <option key='resolution-256' value='256x256'>
              {t('ui.train.image.resolution.256')}
              </option>
            </Form.Select>
            {
              <>
                <Row style={{marginTop:'5px'}}>
                  <Col xs={4}>
                    <Form.Control
                      type='input'
                      size='sm'
                      value={width===-1?t('ui.train.image.resolution.origin'):width}
                      onChange={(e) => selectResolution('width', e.target.value)}
                      disabled={width===-1}
                    />
                  </Col>
                    <Col xs={1}>
                      <span>x</span>
                    </Col>
                  <Col xs={4}>
                    <Form.Control
                      type='input'
                      size='sm'
                      value={height===-1?t('ui.train.image.resolution.origin'):height}
                      onChange={(e) => selectResolution('height', e.target.value)}
                      disabled={height===-1}
                    />
                  </Col>
                </Row>
                {errors && errors['width'] ? (
                  <Form.Control.Feedback type='invalid' className='d-block'>
                    {errors['width']}
                  </Form.Control.Feedback>
                ) : null}
                {errors && errors['height'] ? (
                  <Form.Control.Feedback type='invalid' className='d-block'>
                    {errors['height']}
                  </Form.Control.Feedback>
                ) : null}
                {
                  ((width < 64 && width >= 0) || (height < 64 && height >= 0) )&& (
                    <Form.Control.Feedback type='invalid' className='d-block'>
                      {t('ui.train.image.resolution.desc')}
                    </Form.Control.Feedback>
                  )
                }
              </>
            }
          </Col>
        </Row>
      </Form.Group>
    </div>
  )
}

export default SelectImageResolution
