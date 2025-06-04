import { Card, Col, Row } from "react-bootstrap"
import styled from "styled-components"

interface IGPUContainer {
    state: string;
}

const GPUContainer = styled(Card) <IGPUContainer>`
    margin-right: 5px;
    background-color: #464f5b;
    padding: 5px;
    max-width: 200px;
    flex: 1 0;
    border: 1px solid ${(props) => props.state === 'train' ? 'rgb(0, 153, 0)' : props.state === 'test' ? 'rgb(255, 51, 153)' : ''};
`

interface IStyledI {
    state: string;
}
const StyledI = styled.i<IStyledI>`
    line-height: 35px;
    font-size: 25px;
    color: ${(props) => props.state === 'train' ? 'rgb(0, 153, 0)' : props.state === 'test' ? 'rgb(255, 51, 153)' : ''};
`

const StyledH6 = styled.h6`
    margin-top: 3px;
    margin-bottom: 5px;
`

interface IStyledState {
    state: string;
}
const StyledState = styled.p<IStyledState>`
    display: inline;
    position: fixed;
    margin-top: 5px;
    padding: 3px;
    font-weight: 800;
    font-size: 12px;
    text-transform: uppercase;
    color: ${(props) => props.state === 'train' ? 'rgb(0, 153, 0)' : props.state === 'test' ? 'rgb(255, 51, 153)' : ''};
`

const SyteldH6Context = styled.h6`
    margin: 0;
    padding: 3px;
`

const GPUInfo = ({ gpuInfo }: any) => {
    return (
        <GPUContainer state={gpuInfo.state}>
            <Row>
                <Col>
                    <StyledI state={gpuInfo.state} className='mdi mdi-expansion-card-variant'></StyledI>
                    <StyledState state={gpuInfo.state}>{gpuInfo.state}</StyledState>
                </Col>
            </Row>
            <Row>
                <Col>
                    <StyledH6>[{gpuInfo.id + 1}] {gpuInfo.name.replace('NVIDIA ', '')}</StyledH6>
                </Col>
            </Row>
            <Row>
                <Col>
                    <SyteldH6Context>GPU {gpuInfo.use_gpu}%</SyteldH6Context>
                </Col>
                <Col>
                    <SyteldH6Context>MEM {gpuInfo.use_mem}%</SyteldH6Context>
                </Col>
            </Row>
        </GPUContainer>
    )
}

export default GPUInfo