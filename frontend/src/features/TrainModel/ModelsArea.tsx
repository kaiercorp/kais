import styled from "styled-components"

interface IDivArea {
    isVertical?: boolean
}

const DivArea = styled.div<IDivArea>`
width: ${(props) => (props.isVertical?'200px':'100%')};
border: 1px solid grey;
height: ${(props) => (props.isVertical?'500px':'100%')};
overflow-y: auto;
`

const ModelsArea = ({ isVertical=true, children}: any) => {
    return (
        <DivArea isVertical={isVertical}>{children}</DivArea>
    )
}

export default ModelsArea