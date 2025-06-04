import { Col, Form, Row } from 'react-bootstrap'

type Props = {
  title?: string
  name: string
  label?: string
  value?: string | number
  position?: string
  disabled?: boolean
  checked: boolean
  errors: any
  onChange: any
  onCheck: any
}

const LabelInputWithCheck = ({ title, name, label, value, position, disabled, checked, errors, onChange, onCheck }: Props) => {
  return (
    <Form.Group className={'mb-1'}>
      <Row>
        <Form.Label column='sm' sm={4}>
          {title}
        </Form.Label>
        <Col>
          <Form.Control
            required
            type='text'
            size='sm'
            value={value}
            isInvalid={errors && errors[name] ? true : false}
            onChange={onChange}
            disabled={disabled}
          />
          {errors && errors[name] !== '' ? (
            <Form.Control.Feedback type='invalid' className='d-block'>
              {errors[name]}
            </Form.Control.Feedback>
          ) : null}
        </Col>
        <Col column='sm' sm={4}>
          <Form.Switch type='switch' checked={checked} label={<Form.Label>{label}</Form.Label>} onChange={onCheck} />
        </Col>
      </Row>
    </Form.Group>
  )
}

export default LabelInputWithCheck
