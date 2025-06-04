import styled from 'styled-components'

import { stateColors } from 'appConstants'
import { Card, Col, Row } from 'react-bootstrap'
import { GPUStatusType } from 'common'

const StyledHeader = styled.div`
display: flex;
flex-direction: row;
justify-content: flex-start;
`

const StyledContent = styled(Col)`
  width: 100%;
  display: flex;
`

const GpuState = styled.div`
  text-align: left;
  width: 120px;
  padding: 10px 0px;
`

type GPUStatusProps = {
  gpus: GPUStatusType
  state: string
}

const GPUStatus = ({ gpus, state }: GPUStatusProps) => {
  return (
    <Card className='tilebox-one' style={{minHeight: '112px'}}>
      <Card.Body>
        <StyledHeader>
          <i className='mdi mdi-expansion-card-variant' style={{top: '-17px', left: '-10px', fontSize:'30px', height:'20px', position:'relative', color: stateColors[state]}}></i>
          <h6 className='text-uppercase mt-0' style={{color: stateColors[state]}}>{state}</h6>
        </StyledHeader>
        <Row>
          <StyledContent>
            {gpus.gpus
              .filter((gpu) => {
                return gpu.state === state
              })
              .map((gpu) => {
                const gpuname = gpu.name
                  .replace('NVIDIA ', '')
                  .replace('GeForce ', '')
                  .replace('RTX ', '')
                  .replace('GTX ', '')
                return (
                  <GpuState key={`gpu-state-${gpu.id}`}>
                    [{gpu.id + 1}] {gpuname}
                  </GpuState>
                )
              })}
          </StyledContent>
        </Row>
      </Card.Body>
    </Card>
  )
}

export default GPUStatus
