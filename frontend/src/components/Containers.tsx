import styled from 'styled-components'
import { Form, Row } from 'react-bootstrap'

export const CardHeaderLeft = styled.div`
    font-weight: 800;
    float: left;
`

export const CardHeaderRight = styled.div`
    float: right;

    & .btn-icon {
        padding: 0 2px;
        font-size: 15px;
    }

    & button {
        margin-left: 3px;
    }
`

export const ButtonArea = styled.div`
    float: left;
    & button {
        margin-right: 5px;
    }
`

export const ConfigLabel = styled(Form.Label)`
    width: 100%;
    margin-bottom: 0;
    padding-top: calc(0.28rem + 1px);
    padding-bottom: calc(0.28rem + 1px);
    font-size: 0.875rem;
    flex: 0 0 auto;
`

export const ConfigRow = styled(Row)`
    border: 1px dotted #999999;
    margin-bottom: 3px;
    &:hover {
        background-color: #aab8c5;
        color: #000000;
    }
`