import { Col, Form, Row } from 'react-bootstrap'
import Select from 'react-select'

type Props = {
  title?: any
  name: string
  className?: string
  options?: any[]
  value?: string | any
  position?: string
  errors?: any
  required?: boolean
  onChange?: any
  isMulti?: boolean
}

const LabelSelect2 = ({ title, name, options, value, position, errors, onChange, required=true, isMulti=true, ...props }: Props) => {
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
          <Select
            isMulti={isMulti}
            name={name}
            className="react-select"
            classNamePrefix="react-select"
            options={options}
            onChange={onChange}
            value={value}
          />
        </Col>
      </Row>
    </Form.Group>
  )
}

export default LabelSelect2
