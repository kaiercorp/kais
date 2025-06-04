import styled from "styled-components"

const DivColumnHeader = styled.div`
font-size: 12px;
color: black;
border: 1px solid grey;
width: 120px;
height: 72px;
padding: 0 5px;
word-break: break-word;
text-align: center;
font-weight: 800;
display: flex;
align-items: center;
justify-content: center;
`

const ModelColumnHeader = ({ children }:any) => {
    return (
        <DivColumnHeader>
            {children}
        </DivColumnHeader>
    )
}

export default ModelColumnHeader