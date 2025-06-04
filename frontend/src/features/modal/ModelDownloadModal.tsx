import { PopoverLabel } from 'components'
import { useState } from 'react'
import { Button, Card, Col, Form, Modal, Row } from 'react-bootstrap'
import { useTranslation } from 'react-i18next'

type ModalProps = {
  show: boolean
  toggle: () => void
  trial: any
  handleDownload: (url: string, filename: string) => void
}

const ModelDownloadModal = ({ show, toggle, trial, handleDownload }: ModalProps) => {
  const [t] = useTranslation('translation')

  const [selectedModels, setSelectedModels] = useState<string[]>([])

  const handleSwitchAll = () => {
    if (selectedModels.length > 0) setSelectedModels([])
    else setSelectedModels(trial.train.model_list)
  }

  const handleSwitch = (name: string) => {
    let newSelectedModels = selectedModels
    if (newSelectedModels.includes(name)) {
      let intermediat = newSelectedModels.join(',')
      intermediat = intermediat.replace(`${name},`, '').replace(name, '')
      if (intermediat === '') newSelectedModels = []
      else newSelectedModels = intermediat.split(',')
    } else {
      newSelectedModels.push(name)
    }
    setSelectedModels(newSelectedModels)
  }

  // const handleConfirm = () => {
  //   dispatch(
  //     trialRequestAction(TrialActionTypes.TRAIN_MODEL_LIST_LINK_REQUEST, { trial_id: trial.trial_id, model_list: selectedModels })
  //   )
  // }

  // useEffect(() => {
  //   if (!show || modelsLink === null) return

  //   let filename = modelsLink.split('/')
  //   filename = filename[filename.length - 1]
  //   handleDownload(modelsLink, filename)

  //   toggle()
  //   // eslint-disable-next-line react-hooks/exhaustive-deps
  // }, [show, modelsLink])

  if (trial.project_id < 1) return <></>

  return (
    <Modal show={show} keyboard={true} onHide={toggle}>
      <Modal.Header>
        <span>
          <i className={t(`train.savemodel.icon`)} />
          {t(`train.savemodel.title`)}
          <PopoverLabel name='train-description'>{t(`train.savemodel.description`)}</PopoverLabel>
        </span>
      </Modal.Header>
      <Modal.Body>
        <Card>
          <Card.Body>
            {trial.train !== null && trial.train.model_list && trial.train.model_list.length > 0 ? (
              <>
                <Row>
                  <Form.Label column='sm' sm={4}>
                    {t('button.selectall')}
                  </Form.Label>
                  <Col>
                    <Form.Check
                      type='switch'
                      id='switch-select-all-model'
                      checked={selectedModels.length === trial.train.model_list.length}
                      onChange={() => handleSwitchAll()}
                    ></Form.Check>
                  </Col>
                </Row>
                {trial.train.model_list.map((model: string) => (
                  <Row key={`key-models-${model}`}>
                    <Form.Label column='sm' sm={4}>
                      {model}
                    </Form.Label>
                    <Col>
                      <Form.Check
                        type='switch'
                        id={`model-check-${model}`}
                        checked={selectedModels.includes(model)}
                        onChange={() => handleSwitch(model)}
                      />
                    </Col>
                  </Row>
                ))}
              </>
            ) : (
              <Row>{t('ui.info.nomodel')}</Row>
            )}
          </Card.Body>
        </Card>
      </Modal.Body>
      <Modal.Footer>
        {/* <Button variant='info' className='btn-rounded btn-sm' onClick={handleConfirm}>
          <i className='mdi mdi-content-save' /> 선택 저장
        </Button> */}
        <Button variant='dark' className='btn-rounded btn-sm' onClick={toggle}>
          {t('button.close')}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}

export default ModelDownloadModal
