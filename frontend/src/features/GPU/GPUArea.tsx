import { useContext, useEffect, useState } from 'react'
import { Card, Row } from 'react-bootstrap'
import { GPUContext } from 'contexts'
import { GPUInfo } from 'features/System'

const GPUArea = () => {
  const { gpuContextValue } = useContext(GPUContext)
  const [gpuInfo, setGpuInfo] = useState<any>([])
  useEffect(() => {
    if (!gpuContextValue.gpus) return
    if (gpuContextValue.gpus === undefined) return

    setGpuInfo(gpuContextValue.gpus.slice())
  }, [gpuContextValue.gpus])

  return (
    <Row>
      {
        gpuInfo && gpuInfo.map((gpu: any, index: number) => {
          return <GPUInfo key={`gpu-info-${index}`} gpuInfo={gpu} />
        })
      }
      
      {
        (!gpuInfo || gpuInfo.length < 1) && <Card><Card.Body>NO GPU</Card.Body></Card>
      }
    </Row>
  )
}

export default GPUArea