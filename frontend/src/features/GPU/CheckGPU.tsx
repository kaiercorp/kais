import { useState, useContext } from 'react'
import { Col, Form, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import { GPUContext } from 'contexts'

type Props = {
  errors?: any
  selectGPU: (key: string, value: number[]) => void
}

const CheckGPU = ({ errors, selectGPU }: Props) => {
  const [t] = useTranslation('translation')
  const { gpuContextValue } = useContext(GPUContext)

  const [gpus, setGpus] = useState<number[]>([0])
  const handleClick = (e: any) => {
    const value = Number(e.target.value)
    let newGpus = gpus
    if (newGpus.includes(value)) newGpus = newGpus.filter((gpu) => gpu !== value)
    else newGpus.push(value)

    setGpus(newGpus)
    selectGPU('gpus', newGpus)
  }

  return (
    <Form.Group className='mb-1'>
      <Row>
        <Form.Label column='sm' sm={4}>
        {t('ui.gpu.title')}
        </Form.Label>
        <Col>
          {errors && errors['gpus'] ? (
            <Form.Control.Feedback type='invalid' className='d-block'>
              {errors['gpus']}
            </Form.Control.Feedback>
          ) : null}
          {gpuContextValue &&
            gpuContextValue.gpus.map((gpu: any) => {
              const gpuname =
                '[' +
                (gpu.id + 1) +
                '] ' +
                gpu.name.replace('NVIDIA ', '').replace('GeForce ', '').replace('RTX ', '').replace('GTX ', '')
              return (
                <Form.Check
                  type='checkbox'
                  key={`gpu-id-${gpu.id}`}
                  id={`gpu-id-${gpu.id}`}
                  disabled={gpu.state !== 'idle'}
                  label={gpuname}
                  value={gpu.id}
                  onClick={handleClick}
                  required={true}
                  checked={gpus.includes(gpu.id)}
                  readOnly
                />
              )
            })}
        </Col>
      </Row>
    </Form.Group>
  )
}

export default CheckGPU
