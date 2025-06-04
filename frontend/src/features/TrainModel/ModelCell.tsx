import styled from "styled-components"

const DivCell = styled.div`
font-size: 12px;
color: black;
border: 1px solid grey;
width: 120px;
height: 50px;
padding: 0 5px;
word-break: break-word;
text-align: right;
font-weight: 400;
display: flex;
align-items: center;
justify-content: right;
`

const ModelCell = ({ children }:any) => {
    return (
        <DivCell>
            {children}
        </DivCell>
    )
}

export default ModelCell