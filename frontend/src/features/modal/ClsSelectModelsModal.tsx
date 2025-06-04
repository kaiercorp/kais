import { objDeepCopy } from 'helpers'
import { useEffect, useState } from 'react'
import { Button, Col, Form, Modal, Row } from 'react-bootstrap'
import styled from 'styled-components'
import { useTranslation } from 'react-i18next'

const Container = styled.div`
max-height: 500px;
overflow-x: hidden;
overflow-y: auto;
`

type ModalProps = {
  show: boolean
  selectOne?: boolean
  models: any
  toggle: () => void
  selectModels: (key: string, value: any) => void
}

const ClsSelectModelsModal = ({ show, selectOne = false, models, toggle, selectModels }: ModalProps) => {
  const [t] = useTranslation('translation')
  const [selectedModels, setSelectedModels] = useState<string[]>([])
  const [selectedModel, setSelectedModel] = useState<string[]>([])

  useEffect(() => {
    if (show) {
      setSelectedModels([])
    }
  }, [show, models])

  const handleSwitchAll = () => {
    if (selectedModels.length > 0) setSelectedModels([])
    else {
      let _selected = models.map((model:any) => {return `${model.uuid}_${model.name}`})
      setSelectedModels(_selected)
    }
  }

  const handleSwitch = (name: string) => {
    if (selectOne) {
      handleSwitchOne(name)
    } else {
      handleSwitchMulti(name)
    }
  }

  const handleSwitchOne = (name: string) => {
    setSelectedModel([name])
    setSelectedModels([name])
  }

  const handleSwitchMulti = (name: string) => {
    let newSelectedModels = objDeepCopy(selectedModels)
    const index = newSelectedModels.indexOf(name)
    
    if (index === -1) {
      newSelectedModels.push(name)
    } else {
      newSelectedModels.splice(index, 1)
    }

    setSelectedModels(newSelectedModels)
  }

  const handleConfirm = () => {
    if (selectOne) {
      selectModels('model', selectedModel)
    } else {
      selectModels('model', selectedModels)
    }
    toggle()
  }

  if (!models) return <></>
  return (
    <Modal show={show} keyboard={true} onHide={toggle}>
      <Modal.Header>{t('modal.title.model')}</Modal.Header>
      <Modal.Body>
        {!selectOne && (
          <Row>
            <Col sm={2} style={{paddingTop: '5px'}}>
              <Form.Check
                type='switch'
                id='switch-select-all-model'
                checked={selectedModels.length === models.length}
                onChange={() => handleSwitchAll()}
              ></Form.Check>
            </Col>
            <Form.Label column='sm'>
              {t('button.selectall')}
            </Form.Label>
          </Row>
        )}
        <Container>
        {models.map((model: any) => {
          return (
            <Row key={`models-${model.train_id}${model.name}`}>
              <Col sm={1} style={{paddingTop: '5px'}}>
                <Form.Check
                  type='switch'
                  id={`switch-select-${model.train_id}${model}`}
                  checked={selectedModels.includes(`${model.uuid}_${model.name}`)}
                  onChange={() => handleSwitch(`${model.uuid}_${model.name}`)}
                ></Form.Check>
              </Col>
              <Form.Label column='sm' sm={6} htmlFor={`switch-select-${model.train_id}${model}`}>
              {`Trial${model.train_id} - `}<b>{model.name}</b> 
              </Form.Label>
              <Form.Label column='sm' htmlFor={`switch-select-${model.train_id}${model}`}>
                {`${(parseFloat(model.score.String.replaceAll("\"", "")) * 100).toFixed(2)}%`}
              </Form.Label>
            </Row>
          )
        })}
        </Container>
      </Modal.Body>
      <Modal.Footer>
        <Button
          variant='info'
          className='btn-rounded btn-sm'
          onClick={handleConfirm}
          disabled={selectedModels.length < 1}
        >
          {t('button.apply')}
        </Button>
        <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
          {t('button.cancel')}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ClsSelectModelsModal
