import styled from "styled-components"

interface IDivArea {
    isVertical?: boolean
}

const DivArea = styled.div<IDivArea>`
background-color: #ffffff;
color: #5f5f5f !important;
padding: 10px;
${(props) => (props.isVertical && 'display: flex;')}
width: ${(props) => (props.isVertical ? '200px' : '100%')};

& table {
    & th {
        color: #5f5f5f !important;
    }
}
`

const TrainModelContainer = ({ isVertical=true, children }:any) => {
    return (
        <DivArea isVertical={isVertical}>
            {children}
        </DivArea>
    )
}

export default TrainModelContainer