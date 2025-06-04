import { Col, Form, Row } from 'react-bootstrap'
import { Typeahead } from 'react-bootstrap-typeahead';
import 'react-bootstrap-typeahead/css/Typeahead.css';

export type TypeaheadOption = string | Record<string, Object>

type Props = {
    title?: any
    name: string
    className?: string
    options: TypeaheadOption[]
    value: TypeaheadOption[]
    position?: string
    errors?: any
    required?: boolean
    onChange?: any
    multiple?: boolean
}

const LabelSelectTyped = ({ title, name, options, value, position, errors, onChange, multiple=true, required = false, ...props }: Props) => {
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
                    <Typeahead
                        id="label-select-typed"
                        labelKey="label"
                        multiple={multiple}
                        onChange={onChange}
                        options={options}
                        placeholder="Choose a state..."
                        selected={value}
                    />
                </Col>
            </Row>
        </Form.Group>
    )
}

export default LabelSelectTyped
