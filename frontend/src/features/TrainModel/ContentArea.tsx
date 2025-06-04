import styled from "styled-components"

const DivArea = styled.div`
width: 100%;
overflow-x: scroll;
display: flex;
`

const ContentArea = ({ children }:any) => {
    return (
        <DivArea>
            {children}
        </DivArea>
    )
}

export default ContentArea