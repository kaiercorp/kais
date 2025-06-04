import styled from 'styled-components'
import { Col } from 'react-bootstrap'

const ConfigCol = styled(Col)`
  font-size: 12px;
  & table {
    width: 100%;
    border: 1px solid grey;
    border-collapse: collapse;

    & th {
      font-weight: 600;
      text-align: center;
      padding: 5px;
      border: 1px solid grey;
      color: #ffffff;
    }

    & td {
      text-align: center;
      padding: 5px 10px;
      border: 1px solid grey;
    }
  }
`

export default ConfigCol