import styled from "styled-components"

const DivColumn = styled.div``

const ModelColumn = ({children}:any) => {
    return (
        <DivColumn>
            {children}
        </DivColumn>
    )
}

export default ModelColumn