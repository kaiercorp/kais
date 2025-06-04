import { Col, Form, Row } from 'react-bootstrap'

type Props = {
  title?: string
  name: string
  value?: string | number
  position?: string
  errors: any
  onChange: any
  type?: string
  requeired?: boolean
  onKeyPress?: (e:any) => void
}

const LabelInput = ({ title, name, value, position, errors, onChange, type='text', requeired, onKeyPress }: Props) => {
  return (
    <Form.Group className={'mb-1'}>
      <Row>
        <Form.Label column='sm' sm={4}>
          {title}
        </Form.Label>
        <Col>
          {errors && errors[name] !== '' ? (
            <Form.Control.Feedback type='invalid' className='d-block'>
              {errors[name]}
            </Form.Control.Feedback>
          ) : null}
          <Form.Control
            required={requeired?requeired:false}
            type={type}
            size='sm'
            value={value}
            isInvalid={errors && errors[name] ? true : false}
            onChange={onChange}
            onKeyDown={onKeyPress && ((e:any) => onKeyPress(e))}
          />
        </Col>
      </Row>
    </Form.Group>
  )
}

export default LabelInput
