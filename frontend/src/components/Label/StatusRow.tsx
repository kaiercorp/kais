import styled from 'styled-components'

const StatusRow = styled.div`
display: flex;
align-items: left;
text-align: left;
margin-bottom: 5px;
border-bottom: 1px dotted grey;

& div {
    margin-right: 10px;
}
& span {
    font-size: 12px;
    width: 180px;
}
`

export default StatusRow