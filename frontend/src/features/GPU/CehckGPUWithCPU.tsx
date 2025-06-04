import { DefaultSwitch } from 'components'
import { useState, useContext } from 'react'
import { Col, Form, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'
import { GPUContext } from 'contexts'

type Props = {
    errors?: any
    selectGPU: (key: string, value: any[]) => void
}

const CheckGPUWithCPU = ({ errors, selectGPU }: Props) => {
    const [t] = useTranslation('translation')
    const { gpuContextValue } = useContext(GPUContext)

    const [gpus, setGpus] = useState<number[]>([])
    const handleClick = (e: any) => {
        const value = Number(e.target.value)
        let newGpus = gpus
        if (newGpus.includes(value)) newGpus = newGpus.filter((gpu) => gpu !== value)
        else newGpus.push(value)

        setGpus(newGpus)
        selectGPU('gpus', newGpus)
    }

    const [isCpu, setIsCpu] = useState(true)
    const toggleIsCpu = () => {
        setIsCpu(!isCpu)
        if (!isCpu) {
            selectGPU('gpus', [])
        }
        else {
            selectGPU('gpus', gpus)
        }
    }

    return (
        <Form.Group className='mb-1'>
            <Row>
                <Form.Label column='sm' sm={4}>
                    {t('ui.gpu.title')}
                </Form.Label>
                <Col sm={8}>
                    <Row>
                        {errors && errors['gpus'] ? (
                            <Form.Control.Feedback type='invalid' className='d-block'>
                                {errors['gpus']}
                            </Form.Control.Feedback>
                        ) : null}
                    </Row>
                    <Row>
                        <Col column='sm' sm={6}>
                            <DefaultSwitch label={t('ui.gpu.useCpu')} checked={isCpu} onChange={toggleIsCpu} />
                        </Col>
                        <Col column='sm' sm={6} disabled={isCpu}>
                            {!isCpu && gpuContextValue &&
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
                                            disabled={gpu.state !== 'idle' || isCpu}
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
                </Col>
            </Row>
        </Form.Group>
    )
}

export default CheckGPUWithCPU
