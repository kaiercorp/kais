import { useState, useContext, useEffect } from 'react'
import { Col, Form, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import { GPUContext } from 'contexts'

type Props = {
  errors?: any
  selectGPU: (key: string, value: number[]) => void
}

const RadioGPU = ({ errors, selectGPU }: Props) => {
  const [t] = useTranslation('translation')
  const { gpuContextValue } = useContext(GPUContext)

  const [gpus, setGpus] = useState<number|undefined>()
  const handleClick = (e: any) => {
    const value = Number(e.target.value)
    if (value === gpus) {
      setGpus(undefined)
      selectGPU('gpus', [])
    } else {
      setGpus(value)
      selectGPU('gpus', [value])
    }
  }

  useEffect(() => {
    if (!gpuContextValue.gpus) return

    let idle = -1
    gpuContextValue.gpus.forEach((gpu: any) => {
      if (!gpu.is_running && idle === -1) {
        idle = gpu.id
      }
    })
    handleClick({target:{value:idle}})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

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
                  type='radio'
                  key={`gpu-id-${gpu.id}`}
                  id={`gpu-id-${gpu.id}`}
                  disabled={gpu.state !== 'idle'}
                  label={gpuname}
                  value={gpu.id}
                  onClick={handleClick}
                  required={true}
                  readOnly
                  checked={gpus === gpu.id}
                  // checked={true}
                />
              )
            })}
        </Col>
      </Row>
    </Form.Group>
  )
}

export default RadioGPU
