import { useEffect } from "react"
import { Card, Col, Row } from "react-bootstrap"
import styled from 'styled-components'

const ProbaLabel = styled.label`
font-size: 12px;
`

interface IProbaBar {
    left: number
    right: number
}
const ProbaBar = styled.div<IProbaBar>`
background-image: linear-gradient(to right, skyblue ${(props) => props.left}%, white 0%);
width: 100%;
color: black;
font-size: 12px;
padding-left: 10px;
`

const PredictionProbabilities = ({data}: any) => {
    useEffect(() => {
       if (!data)  return
    }, [data])
    return (
        <Card>
            <Card.Header>Prediction probabilities</Card.Header>
            <Card.Body>
                {
                    data.label && (
                        data.label.map((proba:any, index:number) => {
                            const left = Number(Number(data.proba[index] || 0).toFixed(0)) * 100
                            const right = 100 - left
                            return (
                                <Row key={`proba-${index}`}>
                                    <Col sm={4}>
                                        <ProbaLabel>{proba}</ProbaLabel>
                                    </Col>  
                                    <Col sm={8}>
                                        <ProbaBar left={left} right={right}>{data.proba[index]}</ProbaBar>
                                    </Col>
                                </Row>
                            )
                        })
                    )
                }
            </Card.Body>
        </Card>
    )
}

export default PredictionProbabilities