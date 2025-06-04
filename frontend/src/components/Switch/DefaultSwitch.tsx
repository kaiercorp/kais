import { Form } from 'react-bootstrap'
import styled from 'styled-components'


const StyledSwitch = styled(Form.Check)`
& :checked {
    background-color: ${(props) => props.bgcolor && props.bgcolor};
    border-color: ${(props) => props.bgcolor && props.bgcolor};
}
`

const SwitchLabel = (props: any) => {
    return (
        <div {...props}>
            {props.image? <img src={props.image} style={{marginBottom:'2px'}} alt={props.label} className="me-0 me-sm-1" height="14" /> : null} {props.label}
        </div>
    )
}

const DefaultSwitch = ({label, checked, onChange, image, disabled}: any) => {

    return (
        <StyledSwitch disabled={!!disabled}
            id={`switch-${label}`}
            type='switch'
            checked={checked}
            label={<SwitchLabel htmlFor={`switch-${label}`} label={label} image={image} />}
            onChange={onChange}
        />
    )
}

export default DefaultSwitch