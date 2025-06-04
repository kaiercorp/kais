import { Col, Form, Row } from 'react-bootstrap'

type Props = {
  title?: any
  name: string
  className?: string
  children?: React.ReactNode
  value?: string
  position?: string
  errors?: any
  required?: boolean
  onChange?: any
}

const LabelSelect = ({ title, name, children, value, position, errors, onChange, required=true, ...props }: Props) => {
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
          <Form.Select 
            name={name} 
            size='sm' 
            isInvalid={errors && errors[name] ? true : false} 
            value={value} 
            onChange={onChange}
            {...props} 
          >
            {children}
          </Form.Select>
        </Col>
      </Row>
    </Form.Group>
  )
}

export default LabelSelect
